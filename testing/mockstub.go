package testing

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/msp"
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
}

type CreatorTransformer func(...interface{}) (mspID string, certPEM []byte, err error)

// NewMockStub creates chaincode imitation
func NewMockStub(name string, cc shim.Chaincode) *MockStub {
	return &MockStub{
		MockStub:                *shim.NewMockStub(name, cc),
		cc:                      cc,
		ClearCreatorAfterInvoke: true, // by default tx creator data and transient map are cleared after each cc method query/invoke
		InvokablesFull:          make(map[string]*MockStub),
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
		return shim.Error(fmt.Sprintf(`%s: try to invoke chaincode "%s" in channel "%s" (%s). Available mocked chaincodes are: %s`,
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
	rand.Seed(time.Now().UnixNano())
	rand.Read(id)
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
