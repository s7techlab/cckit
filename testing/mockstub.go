package testing

import (
	"container/list"
	"crypto/rand"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/convert"
)

const EventChannelBufferSize = 100

var (
	// ErrChaincodeNotExists occurs when attempting to invoke a nonexostent external chaincode
	ErrChaincodeNotExists = errors.New(`chaincode not exists`)
	// ErrUnknownFromArgsType occurs when attempting to set unknown args in From func
	ErrUnknownFromArgsType = errors.New(`unknown args type to cckit.MockStub.From func`)
	// ErrKeyAlreadyExistsInTransientMap occurs when attempting to set existing key in transient map
	ErrKeyAlreadyExistsInTransientMap = errors.New(`key already exists in transient map`)
)

// MockStub replacement of shim.MockStub with creator mocking facilities
type MockStub struct {
	shim.MockStub
	cc                          shim.Chaincode
	mockCreator                 []byte
	transient                   map[string][]byte
	ClearCreatorAfterInvoke     bool
	_args                       [][]byte
	InvokablesFull              map[string]*MockStub        // invokable this version of MockStub
	creatorTransformer          CreatorTransformer          // transformer for tx creator data, used in From func
	ChaincodeEvent              *peer.ChaincodeEvent        // event in last tx
	chaincodeEventSubscriptions []chan *peer.ChaincodeEvent // multiple event subscriptions
	PrivateKeys                 map[string]*list.List
}

type CreatorTransformer func(...interface{}) (mspID string, certPEM []byte, err error)

// NewMockStub creates chaincode imitation
func NewMockStub(name string, cc shim.Chaincode) *MockStub {
	return &MockStub{
		MockStub: *shim.NewMockStub(name, cc),
		cc:       cc,
		// by default tx creator data and transient map are cleared after each cc method query/invoke
		ClearCreatorAfterInvoke: true,
		InvokablesFull:          make(map[string]*MockStub),
		PrivateKeys:             make(map[string]*list.List),
	}
}

// GetArgs mocked args
func (stub *MockStub) GetArgs() [][]byte {
	return stub._args
}

// SetArgs set mocked args
func (stub *MockStub) SetArgs(args [][]byte) {
	stub._args = args
}

// SetEvent sets chaincode event
func (stub *MockStub) SetEvent(name string, payload []byte) error {
	if name == "" {
		return errors.New("event name can not be nil string")
	}

	stub.ChaincodeEvent = &peer.ChaincodeEvent{EventName: name, Payload: payload}
	for _, sub := range stub.chaincodeEventSubscriptions {
		sub <- stub.ChaincodeEvent
	}

	return stub.MockStub.SetEvent(name, payload)
}

func (stub *MockStub) EventSubscription() chan *peer.ChaincodeEvent {
	subscription := make(chan *peer.ChaincodeEvent, EventChannelBufferSize)
	stub.chaincodeEventSubscriptions = append(stub.chaincodeEventSubscriptions, subscription)
	return subscription
}

// ClearEvents clears chaincode events channel
func (stub *MockStub) ClearEvents() {
	for len(stub.ChaincodeEventsChannel) > 0 {
		<-stub.ChaincodeEventsChannel
	}
}

// GetStringArgs get mocked args as strings
func (stub *MockStub) GetStringArgs() []string {
	args := stub.GetArgs()
	strargs := make([]string, 0, len(args))
	for _, barg := range args {
		strargs = append(strargs, string(barg))
	}
	return strargs
}

// MockPeerChaincode link to another MockStub
func (stub *MockStub) MockPeerChaincode(invokableChaincodeName string, otherStub *MockStub) {
	stub.InvokablesFull[invokableChaincodeName] = otherStub
}

// MockedPeerChaincodes returns names of mocked chaincodes, available for invoke from current stub
func (stub *MockStub) MockedPeerChaincodes() []string {
	keys := make([]string, 0)
	for k := range stub.InvokablesFull {
		keys = append(keys, k)
	}
	return keys
}

// InvokeChaincode using another MockStub
func (stub *MockStub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) peer.Response {
	// Internally we use chaincode name as a composite name
	ccName := chaincodeName
	if channel != "" {
		chaincodeName = chaincodeName + "/" + channel
	}

	otherStub, exists := stub.InvokablesFull[chaincodeName]
	if !exists {
		return shim.Error(fmt.Sprintf(
			`%s: try to invoke chaincode "%s" in channel "%s" (%s). Available mocked chaincodes are: %s`,
			ErrChaincodeNotExists, ccName, channel, chaincodeName, stub.MockedPeerChaincodes()))
	}

	res := otherStub.MockInvoke(stub.TxID, args)
	return res
}

// GetFunctionAndParameters mocked
func (stub *MockStub) GetFunctionAndParameters() (function string, params []string) {
	allargs := stub.GetStringArgs()
	function = ""
	params = []string{}
	if len(allargs) >= 1 {
		function = allargs[0]
		params = allargs[1:]
	}
	return
}

// RegisterCreatorTransformer  that transforms creator data to MSP_ID and X.509 certificate
func (stub *MockStub) RegisterCreatorTransformer(creatorTransformer CreatorTransformer) *MockStub {
	stub.creatorTransformer = creatorTransformer
	return stub
}

// MockCreator of tx
func (stub *MockStub) MockCreator(mspID string, certPEM []byte) {
	stub.mockCreator, _ = msp.NewSerializedIdentity(mspID, certPEM)
}

func (stub *MockStub) generateTxUID() string {
	id := make([]byte, 32)
	if _, err := rand.Read(id); err != nil {
		panic(err)
	}
	return fmt.Sprintf("0x%x", id)
}

// Init func of chaincode - sugared version with autogenerated tx uuid
func (stub *MockStub) Init(iargs ...interface{}) peer.Response {
	args, err := convert.ArgsToBytes(iargs...)
	if err != nil {
		return shim.Error(err.Error())
	}

	return stub.MockInit(stub.generateTxUID(), args)
}

// InitBytes init func with ...[]byte args
func (stub *MockStub) InitBytes(args ...[]byte) peer.Response {
	return stub.MockInit(stub.generateTxUID(), args)
}

// MockInit mocked init function
func (stub *MockStub) MockInit(uuid string, args [][]byte) peer.Response {
	stub.SetArgs(args)
	stub.MockTransactionStart(uuid)
	res := stub.cc.Init(stub)
	stub.MockTransactionEnd(uuid)

	if stub.ClearCreatorAfterInvoke {
		stub.mockCreator = nil
		stub.transient = nil
	}

	return res
}

// MockQuery
func (stub *MockStub) MockQuery(uuid string, args [][]byte) peer.Response {
	return stub.MockInvoke(uuid, args)
}

// MockInvoke
func (stub *MockStub) MockInvoke(uuid string, args [][]byte) peer.Response {
	// this is a hack here to set MockStub.args, because its not accessible otherwise
	stub.SetArgs(args)

	//empty event
	stub.ChaincodeEvent = nil

	// now do the invoke with the correct stub
	stub.MockTransactionStart(uuid)
	res := stub.cc.Invoke(stub)
	stub.MockTransactionEnd(uuid)

	if stub.ClearCreatorAfterInvoke {
		stub.mockCreator = nil
		stub.transient = nil
	}

	return res
}

// Invoke sugared invoke function with autogenerated tx uuid
func (stub *MockStub) Invoke(funcName string, iargs ...interface{}) peer.Response {
	fargs, err := convert.ArgsToBytes(iargs...)
	if err != nil {
		return shim.Error(err.Error())
	}
	args := append([][]byte{[]byte(funcName)}, fargs...)
	return stub.InvokeBytes(args...)
}

// InvokeByte mock invoke with autogenerated tx uuid
func (stub *MockStub) InvokeBytes(args ...[]byte) peer.Response {
	return stub.MockInvoke(stub.generateTxUID(), args)
}

// QueryBytes mock query with autogenerated tx uuid
func (stub *MockStub) QueryBytes(args ...[]byte) peer.Response {
	return stub.MockQuery(stub.generateTxUID(), args)
}

func (stub *MockStub) Query(funcName string, iargs ...interface{}) peer.Response {
	return stub.Invoke(funcName, iargs...)
}

// GetCreator mocked
func (stub *MockStub) GetCreator() ([]byte, error) {
	return stub.mockCreator, nil
}

// From mock tx creator
func (stub *MockStub) From(txCreator ...interface{}) *MockStub {

	var mspID string
	var certPEM []byte
	var err error

	if stub.creatorTransformer != nil {
		mspID, certPEM, err = stub.creatorTransformer(txCreator...)
	} else {
		mspID, certPEM, err = TransformCreator(txCreator...)
	}

	if err != nil {
		panic(err)
	}
	stub.MockCreator(mspID, certPEM)
	return stub
}

func (stub *MockStub) GetTransient() (map[string][]byte, error) {
	return stub.transient, nil
}

// WithTransient sets transient map
func (stub *MockStub) WithTransient(transient map[string][]byte) *MockStub {
	stub.transient = transient
	return stub
}

// AddTransient adds key-value pairs to transient map
func (stub *MockStub) AddTransient(transient map[string][]byte) *MockStub {
	if stub.transient == nil {
		stub.transient = make(map[string][]byte)
	}
	for k, v := range transient {
		if _, ok := stub.transient[k]; ok {
			panic(ErrKeyAlreadyExistsInTransientMap)
		}
		stub.transient[k] = v
	}
	return stub
}

// At mock tx timestamp
//func (stub *MockStub) At(txTimestamp *timestamp.Timestamp) *MockStub {
//	stub.TxTimestamp = txTimestamp
//	return stub
//}

// DelPrivateData mocked
func (stub *MockStub) DelPrivateData(collection string, key string) error {
	m, in := stub.PvtState[collection]
	if !in {
		return errors.Errorf("Collection %s not found.", collection)
	}

	if _, ok := m[key]; !ok {
		return errors.Errorf("Key %s not found.", key)
	}
	delete(m, key)

	for elem := stub.PrivateKeys[collection].Front(); elem != nil; elem = elem.Next() {
		if strings.Compare(key, elem.Value.(string)) == 0 {
			stub.PrivateKeys[collection].Remove(elem)
		}
	}
	return nil
}

type PrivateMockStateRangeQueryIterator struct {
	Closed     bool
	Stub       *MockStub
	StartKey   string
	EndKey     string
	Current    *list.Element
	Collection string
}

// Logger for the shim package.
var mockLogger = logging.MustGetLogger("mock")

// HasNext returns true if the range query iterator contains additional keys
// and values.
func (iter *PrivateMockStateRangeQueryIterator) HasNext() bool {
	if iter.Closed {
		// previously called Close()
		mockLogger.Debug("HasNext() but already closed")
		return false
	}

	if iter.Current == nil {
		mockLogger.Error("HasNext() couldn't get Current")
		return false
	}

	current := iter.Current
	for current != nil {
		// if this is an open-ended query for all keys, return true
		if iter.StartKey == "" && iter.EndKey == "" {
			return true
		}
		comp1 := strings.Compare(current.Value.(string), iter.StartKey)
		comp2 := strings.Compare(current.Value.(string), iter.EndKey)
		if comp1 >= 0 {
			if comp2 < 0 {
				mockLogger.Debug("HasNext() got next")
				return true
			} else {
				mockLogger.Debug("HasNext() but no next")
				return false

			}
		}
		current = current.Next()
	}

	// we've reached the end of the underlying values
	mockLogger.Debug("HasNext() but no next")
	return false
}

// Next returns the next key and value in the range query iterator.
func (iter *PrivateMockStateRangeQueryIterator) Next() (*queryresult.KV, error) {
	if iter.Closed == true {
		err := errors.New("PrivateMockStateRangeQueryIterator.Next() called after Close()")
		mockLogger.Errorf("%+v", err)
		return nil, err
	}

	if iter.HasNext() == false {
		err := errors.New("PrivateMockStateRangeQueryIterator.Next() called when it does not HaveNext()")
		mockLogger.Errorf("%+v", err)
		return nil, err
	}

	for iter.Current != nil {
		comp1 := strings.Compare(iter.Current.Value.(string), iter.StartKey)
		comp2 := strings.Compare(iter.Current.Value.(string), iter.EndKey)
		// compare to start and end keys. or, if this is an open-ended query for
		// all keys, it should always return the key and value
		if (comp1 >= 0 && comp2 < 0) || (iter.StartKey == "" && iter.EndKey == "") {
			key := iter.Current.Value.(string)
			value, err := iter.Stub.GetPrivateData(iter.Collection, key)
			iter.Current = iter.Current.Next()
			return &queryresult.KV{Key: key, Value: value}, err
		}
		iter.Current = iter.Current.Next()
	}
	err := errors.New("PrivateMockStateRangeQueryIterator.Next() went past end of range")
	mockLogger.Errorf("%+v", err)
	return nil, err
}

// Close closes the range query iterator. This should be called when done
// reading from the iterator to free up resources.
func (iter *PrivateMockStateRangeQueryIterator) Close() error {
	if iter.Closed == true {
		err := errors.New("PrivateMockStateRangeQueryIterator.Close() called after Close()")
		mockLogger.Errorf("%+v", err)
		return err
	}

	iter.Closed = true
	return nil
}

func NewPrivateMockStateRangeQueryIterator(stub *MockStub, collection string, startKey string, endKey string) *PrivateMockStateRangeQueryIterator {
	mockLogger.Debug("NewPrivateMockStateRangeQueryIterator(", stub, startKey, endKey, ")")
	if _, ok := stub.PrivateKeys[collection]; !ok {
		stub.PrivateKeys[collection] = list.New()
	}
	iter := new(PrivateMockStateRangeQueryIterator)
	iter.Closed = false
	iter.Stub = stub
	iter.StartKey = startKey
	iter.EndKey = endKey
	iter.Current = stub.PrivateKeys[collection].Front()
	iter.Collection = collection

	iter.Print()

	return iter
}

func (iter *PrivateMockStateRangeQueryIterator) Print() {
	mockLogger.Debug("PrivateMockStateRangeQueryIterator {")
	mockLogger.Debug("Closed?", iter.Closed)
	mockLogger.Debug("Stub", iter.Stub)
	mockLogger.Debug("StartKey", iter.StartKey)
	mockLogger.Debug("EndKey", iter.EndKey)
	mockLogger.Debug("Current", iter.Current)
	mockLogger.Debug("HasNext?", iter.HasNext())
	mockLogger.Debug("Collection", iter.Collection)
	mockLogger.Debug("}")
}

// PutPrivateData mocked
func (stub *MockStub) PutPrivateData(collection string, key string, value []byte) error {
	if _, in := stub.PvtState[collection]; !in {
		stub.PvtState[collection] = make(map[string][]byte)
	}
	stub.PvtState[collection][key] = value

	if _, ok := stub.PrivateKeys[collection]; !ok {
		stub.PrivateKeys[collection] = list.New()
	}

	for elem := stub.PrivateKeys[collection].Front(); elem != nil; elem = elem.Next() {
		elemValue := elem.Value.(string)
		comp := strings.Compare(key, elemValue)
		mockLogger.Debug("MockStub", stub.Name, "Compared", key, elemValue, " and got ", comp)
		if comp < 0 {
			// key < elem, insert it before elem
			stub.PrivateKeys[collection].InsertBefore(key, elem)
			mockLogger.Debug("MockStub", stub.Name, "Key", key, " inserted before", elem.Value)
			break
		} else if comp == 0 {
			// keys exists, no need to change
			mockLogger.Debug("MockStub", stub.Name, "Key", key, "already in State")
			break
		} else { // comp > 0
			// key > elem, keep looking unless this is the end of the list
			if elem.Next() == nil {
				stub.PrivateKeys[collection].PushBack(key)
				mockLogger.Debug("MockStub", stub.Name, "Key", key, "appended")
				break
			}
		}
	}

	// special case for empty Keys list
	if stub.PrivateKeys[collection].Len() == 0 {
		stub.PrivateKeys[collection].PushFront(key)
		mockLogger.Debug("MockStub", stub.Name, "Key", key, "is first element in list")
	}

	return nil
}

const maxUnicodeRuneValue = utf8.MaxRune

// GetPrivateDataByPartialCompositeKey mocked
func (stub *MockStub) GetPrivateDataByPartialCompositeKey(collection, objectType string, attributes []string) (shim.StateQueryIteratorInterface, error) {
	partialCompositeKey, err := stub.CreateCompositeKey(objectType, attributes)
	if err != nil {
		return nil, err
	}
	return NewPrivateMockStateRangeQueryIterator(stub, collection, partialCompositeKey, partialCompositeKey+string(maxUnicodeRuneValue)), nil
}
