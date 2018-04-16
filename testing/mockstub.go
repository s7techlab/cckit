package testing

import (
	"github.com/pkg/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/convert"
)

var (
	ErrChaincodeNotExists = errors.New(`chaincode not exists`)
)

type ToBytes interface {
	ToBytes() []byte
}

type MockStub struct {
	shim.MockStub

	cc shim.Chaincode

	mockCreator             []byte
	ClearCreatorAfterInvoke bool
	_args                   [][]byte
	InvokablesFull          map[string]*MockStub
	creatorTransformer      func(...interface{}) (mspID, cert string)
}

func NewMockStub(name string, cc shim.Chaincode) *MockStub {
	s := shim.NewMockStub(name, cc)
	fs := new(MockStub)
	fs.MockStub = *s
	fs.cc = cc
	fs.InvokablesFull = make(map[string]*MockStub)
	return fs
}

func (stub *MockStub) GetArgs() [][]byte {
	return stub._args
}

func (stub *MockStub) SetArgs(args [][]byte) {
	stub._args = args
}

func (stub *MockStub) GetStringArgs() []string {
	args := stub.GetArgs()
	strargs := make([]string, 0, len(args))
	for _, barg := range args {
		strargs = append(strargs, string(barg))
	}
	return strargs
}

func (stub *MockStub) MockPeerChaincode(invokableChaincodeName string, otherStub *MockStub) {
	stub.InvokablesFull[invokableChaincodeName] = otherStub
}

func (stub *MockStub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) peer.Response {

	// TODO "args" here should possibly be a serialized pb.ChaincodeInput
	// Internally we use chaincode name as a composite name
	if channel != "" {
		chaincodeName = chaincodeName + "/" + channel
	}

	otherStub, exists := stub.InvokablesFull[chaincodeName]
	if !exists {
		return shim.Error(ErrChaincodeNotExists.Error())
	}

	res := otherStub.MockInvoke(stub.TxID, args)
	return res
}

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

func (stub *MockStub) RegisterCreatorTransformer(transformer func(...interface{}) (mspID, cert string)) *MockStub {
	stub.creatorTransformer = transformer
	return stub
}

func (stub *MockStub) MockCreator(mspID string, cert string) {
	stub.mockCreator, _ = msp.NewSerializedIdentity(mspID, []byte(cert))
}

func (stub *MockStub) generateTxUid() string {
	return "xxx"
}

func (stub *MockStub) Init(iargs ...interface{}) peer.Response {
	args, err := convert.ArgsToBytes(iargs...)
	if err != nil {
		return shim.Error(err.Error())
	}

	return stub.MockInit(stub.generateTxUid(), args)
}

func (stub *MockStub) MockInit(uuid string, args [][]byte) peer.Response {

	//default method name
	//if len(args) == 0 || string(args[0]) != "Init" {
	//	args = append([][]byte{[]byte("Init")}, args...)
	//}

	stub.SetArgs(args)
	stub.MockTransactionStart(uuid)
	res := stub.cc.Init(stub)
	stub.MockTransactionEnd(uuid)

	if stub.ClearCreatorAfterInvoke {
		stub.mockCreator = nil
	}

	return res
}

func (stub *MockStub) MockInvoke(uuid string, args [][]byte) peer.Response {

	// this is a hack here to set MockStub.args, because its not accessible otherwise
	stub.SetArgs(args)

	// now do the invoke with the correct stub
	stub.MockTransactionStart(uuid)
	res := stub.cc.Invoke(stub)
	stub.MockTransactionEnd(uuid)

	if stub.ClearCreatorAfterInvoke {
		stub.mockCreator = nil
	}

	return res
}


func (stub *MockStub) Invoke(funcName string, iargs ...interface{}) peer.Response {
	return stub.MockInvokeFunc(funcName, iargs...)
}

func (stub *MockStub) MockInvokeFunc(funcName string, iargs ...interface{}) peer.Response {

	fargs, err := convert.ArgsToBytes(iargs...)
	if err != nil {
		return shim.Error(err.Error())
	}
	args := append([][]byte{[]byte(funcName)}, fargs...)
	return stub.MockInvoke(stub.generateTxUid(), args)
}

func (stub *MockStub) GetCreator() ([]byte, error) {
	return stub.mockCreator, nil
}

func (stub *MockStub) From(mspParams ...interface{}) *MockStub {
	var mspID, cert string

	if stub.creatorTransformer != nil {
		mspID, cert = stub.creatorTransformer(mspParams...)
	} else if len(mspParams) == 1 {

		switch mspParams[0].(type) {

		// array with 2 elements  - mspId and ca cert
		case [2]string:
			mspID = mspParams[0].([2]string)[0]
			cert = mspParams[0].([2]string)[1]
			//stub.MockCreator(mspParams[0].([2]string)[0], mspParams[0].([2]string)[1])
		default:
			panic(`Unknow params type to Fullmockstub.From func`)
		}
	} else if len(mspParams) == 2 {
		mspID = mspParams[0].(string)
		cert = mspParams[1].(string)
	}

	stub.MockCreator(mspID, cert)
	return stub
}
