package testing

import (
	"container/list"
	"crypto/rand"
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
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

type StateItem struct {
	Key    string
	Value  []byte
	Delete bool
}

// MockStub replacement of shim.MockStub with creator mocking facilities
type MockStub struct {
	shimtest.MockStub
	cc shim.Chaincode

	StateBuffer []*StateItem // buffer for state changes during transaction

	m sync.Mutex

	_args       [][]byte
	transient   map[string][]byte
	mockCreator []byte
	TxResult    peer.Response // last tx result

	ClearCreatorAfterInvoke bool
	creatorTransformer      CreatorTransformer // transformer for tx creator data, used in From func

	Invokables map[string]*MockStub // invokable this version of MockStub

	LastTxID                    string
	ChaincodeEvent              *peer.ChaincodeEvent        // event in last tx
	chaincodeEventSubscriptions []chan *peer.ChaincodeEvent // multiple event subscriptions

	PrivateKeys map[string]*list.List
	// flag for cc2cc invokation via InvokeChaincode to dump state on outer tx finish
	// https://github.com/s7techlab/cckit/issues/97
	cc2ccInvokation bool
}

type CreatorTransformer func(...interface{}) (mspID string, certPEM []byte, err error)

// NewMockStub creates chaincode imitation
func NewMockStub(name string, cc shim.Chaincode) *MockStub {
	return &MockStub{
		MockStub: *shimtest.NewMockStub(name, cc),
		cc:       cc,
		// by default tx creator data and transient map are cleared after each cc method query/invoke
		ClearCreatorAfterInvoke: true,
		Invokables:              make(map[string]*MockStub),
		PrivateKeys:             make(map[string]*list.List),
	}
}

// PutState wrapped functions puts state items in queue and dumps
// to state after invocation
func (stub *MockStub) PutState(key string, value []byte) error {
	if stub.TxID == "" {
		return errors.New("cannot PutState without a transactions - call stub.MockTransactionStart()?")
	}
	stub.StateBuffer = append(stub.StateBuffer, &StateItem{
		Key:   key,
		Value: value,
	})

	return nil
}

func (stub *MockStub) DelState(key string) error {
	if stub.TxID == "" {
		return errors.New("cannot PutState without a transactions - call stub.MockTransactionStart()?")
	}
	stub.StateBuffer = append(stub.StateBuffer, &StateItem{
		Key:    key,
		Delete: true,
	})

	return nil
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

	stub.ChaincodeEvent = &peer.ChaincodeEvent{
		ChaincodeId: stub.Name,
		TxId:        stub.TxID,
		EventName:   name,
		Payload:     payload,
	}
	return nil
}

// EventSubscription for new or all events
func (stub *MockStub) EventSubscription(from ...int64) (events chan *peer.ChaincodeEvent, closer func() error) {
	stub.m.Lock()
	defer stub.m.Unlock()

	events = make(chan *peer.ChaincodeEvent, EventChannelBufferSize)

	if len(from) > 0 && from[0] == 0 {
		curLen := len(stub.ChaincodeEventsChannel)

		for i := 0; i < curLen; i++ {
			e := <-stub.ChaincodeEventsChannel
			events <- e
			stub.ChaincodeEventsChannel <- e
		}
	}

	stub.chaincodeEventSubscriptions = append(stub.chaincodeEventSubscriptions, events)

	subPos := len(stub.chaincodeEventSubscriptions) - 1
	return events, func() error {
		stub.m.Lock()
		defer stub.m.Unlock()

		if stub.chaincodeEventSubscriptions[subPos] != nil {
			close(stub.chaincodeEventSubscriptions[subPos])
			stub.chaincodeEventSubscriptions[subPos] = nil
		}

		return nil
	}
}

func (stub *MockStub) EventsList() []*peer.ChaincodeEvent {
	stub.m.Lock()
	defer stub.m.Unlock()

	var eventsList []*peer.ChaincodeEvent

	curLen := len(stub.ChaincodeEventsChannel)

	for i := 0; i < curLen; i++ {
		e := <-stub.ChaincodeEventsChannel
		eventsList = append(eventsList, e)
		stub.ChaincodeEventsChannel <- e
	}

	return eventsList
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
	stub.Invokables[invokableChaincodeName] = otherStub
}

// MockedPeerChaincodes returns names of mocked chaincodes, available for invoke from current stub
func (stub *MockStub) MockedPeerChaincodes() []string {
	keys := make([]string, 0)
	for k := range stub.Invokables {
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

	otherStub, exists := stub.Invokables[chaincodeName]
	if !exists {
		return shim.Error(fmt.Sprintf(
			`%s	: try to invoke chaincode "%s" in channel "%s" (%s). Available mocked chaincodes are: %s`,
			ErrChaincodeNotExists, ccName, channel, chaincodeName, stub.MockedPeerChaincodes()))
	}

	otherStub.mockCreator = stub.mockCreator
	otherStub.cc2ccInvokation = true
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
	stub.m.Lock()
	defer stub.m.Unlock()

	stub.SetArgs(args)

	stub.MockTransactionStart(uuid)
	stub.TxResult = stub.cc.Init(stub)
	stub.MockTransactionEnd(uuid)

	return stub.TxResult
}

func (stub *MockStub) DumpStateBuffer() {
	// dump state buffer to state
	if stub.TxResult.Status == shim.OK {
		for i := range stub.StateBuffer {
			s := stub.StateBuffer[i]
			if s.Delete {
				_ = stub.MockStub.DelState(s.Key)
			} else {
				_ = stub.MockStub.PutState(s.Key, s.Value)
			}
		}
	} else {
		stub.ChaincodeEvent = nil
	}
	stub.StateBuffer = nil
}

func (stub *MockStub) dumpEvents() {
	if stub.ChaincodeEvent != nil {
		// send only last event
		for _, sub := range stub.chaincodeEventSubscriptions {
			// subscription can be closed
			if sub != nil {
				sub <- stub.ChaincodeEvent
			}
		}
		stub.ChaincodeEventsChannel <- stub.ChaincodeEvent
	}
}

// MockQuery
func (stub *MockStub) MockQuery(uuid string, args [][]byte) peer.Response {
	return stub.MockInvoke(uuid, args)
}

func (stub *MockStub) MockTransactionStart(uuid string) {
	//empty event
	stub.ChaincodeEvent = nil
	// empty state buffer
	stub.StateBuffer = nil
	stub.TxResult = peer.Response{}

	stub.MockStub.MockTransactionStart(uuid)
}

func (stub *MockStub) MockTransactionEnd(uuid string) {
	stub.LastTxID = stub.TxID
	if !stub.cc2ccInvokation { // skip for inner tx cc2cc calls
		stub.DumpStateBuffer()
		stub.dumpEvents() // events works only for outer stub in Fabric

		// dump buffer to state on outer tx finishing (https://github.com/s7techlab/cckit/issues/97)
		for _, invokableStub := range stub.Invokables {
			invokableStub.DumpStateBuffer()
			invokableStub.cc2ccInvokation = false
		}
		stub.MockStub.MockTransactionEnd(uuid)
	}

	if stub.ClearCreatorAfterInvoke {
		stub.mockCreator = nil
		stub.transient = nil
	}
}

// MockInvoke
func (stub *MockStub) MockInvoke(uuid string, args [][]byte) peer.Response {
	stub.m.Lock()
	defer stub.m.Unlock()

	// this is a hack here to set MockStub.args, because its not accessible otherwise
	stub.SetArgs(args)

	// now do the invoke with the correct stub
	stub.MockTransactionStart(uuid)
	stub.TxResult = stub.cc.Invoke(stub)
	stub.MockTransactionEnd(uuid)

	return stub.TxResult
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

// InvokeBytes mock invoke with autogenerated tx uuid
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

// HasNext returns true if the range query iterator contains additional keys
// and values.
func (iter *PrivateMockStateRangeQueryIterator) HasNext() bool {
	if iter.Closed {
		// previously called Close()
		return false
	}

	if iter.Current == nil {
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
				return true
			} else {
				return false

			}
		}
		current = current.Next()
	}

	// we've reached the end of the underlying values
	return false
}

// Next returns the next key and value in the range query iterator.
func (iter *PrivateMockStateRangeQueryIterator) Next() (*queryresult.KV, error) {
	if iter.Closed {
		err := errors.New("PrivateMockStateRangeQueryIterator.Next() called after Close()")
		return nil, err
	}

	if !iter.HasNext() {
		err := errors.New("PrivateMockStateRangeQueryIterator.Next() called when it does not HaveNext()")
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
	return nil, errors.New("PrivateMockStateRangeQueryIterator.Next() went past end of range")
}

// Close closes the range query iterator. This should be called when done
// reading from the iterator to free up resources.
func (iter *PrivateMockStateRangeQueryIterator) Close() error {
	if iter.Closed {
		return errors.New("PrivateMockStateRangeQueryIterator.Close() called after Close()")
	}

	iter.Closed = true
	return nil
}

func NewPrivateMockStateRangeQueryIterator(stub *MockStub, collection string, startKey string, endKey string) *PrivateMockStateRangeQueryIterator {

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

	return iter
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
		if comp < 0 {
			// key < elem, insert it before elem
			stub.PrivateKeys[collection].InsertBefore(key, elem)
			break
		} else if comp == 0 {
			// keys exists, no need to change
			break
		} else { // comp > 0
			// key > elem, keep looking unless this is the end of the list
			if elem.Next() == nil {
				stub.PrivateKeys[collection].PushBack(key)
				break
			}
		}
	}

	// special case for empty Keys list
	if stub.PrivateKeys[collection].Len() == 0 {
		stub.PrivateKeys[collection].PushFront(key)
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
