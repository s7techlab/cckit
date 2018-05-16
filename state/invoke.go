package state

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/convert"
)

var (
	ErrEmptyChaincodeResponsePayload = errors.New(`empty chaincode response payload`)
)

// InvokeChaincode locally calls the specified chaincode and converts result into target data type
func InvokeChaincode(
	stub shim.ChaincodeStubInterface, chaincodeName string, args []interface{}, channel string, target interface{}) (interface{}, error) {

	convArgs, err := convert.ArgsToBytes(args...)
	if err != nil {
		return nil, err
	}

	response := stub.InvokeChaincode(chaincodeName, convArgs, channel)
	if response.Status != shim.OK {
		return nil, errors.New(response.Message)
	}

	if len(response.Payload) == 0 {
		return nil, ErrEmptyChaincodeResponsePayload
	}

	return convert.FromBytes(response.Payload, target)
}
