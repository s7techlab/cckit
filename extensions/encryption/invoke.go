package encryption

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/state"
)

// InvokeChaincode decrypts received payload
func InvokeChaincode(
	stub shim.ChaincodeStubInterface, encKey []byte, chaincodeName string,
	args []interface{}, channel string, target interface{}) (interface{}, error) {

	// args are not encrypted because we cannot pass encryption key in transient map while invoking cc from cc
	// thus target cc cannot decrypt args
	aa, err := convert.ArgsToBytes(args...)
	if err != nil {
		return nil, fmt.Errorf(`encrypt args: %w`, err)
	}

	response := stub.InvokeChaincode(chaincodeName, aa, channel)
	if response.Status != shim.OK {
		return nil, errors.New(response.Message)
	}

	if len(response.Payload) == 0 {
		return nil, state.ErrEmptyChaincodeResponsePayload
	}

	decrypted, err := Decrypt(encKey, response.Payload)
	if err != nil {
		return nil, fmt.Errorf(`decrypt payload: %w`, err)
	}
	return convert.FromBytes(decrypted, target)
}
