package encryption

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/state"
)

// InvokeChaincode encrypts args before invoking external chaincode and decrypts received payload
func InvokeChaincode(
	stub shim.ChaincodeStubInterface, encKey []byte, chaincodeName string,
	args []interface{}, channel string, target interface{}) (interface{}, error) {

	encryptedArgs, err := EncryptArgs(encKey, args...)
	if err != nil {
		return nil, errors.Wrap(err, `encrypt args`)
	}

	response := stub.InvokeChaincode(chaincodeName, encryptedArgs, channel)
	if response.Status != shim.OK {
		return nil, errors.New(response.Message)
	}

	if len(response.Payload) == 0 {
		return nil, state.ErrEmptyChaincodeResponsePayload
	}

	decrypted, err := Decrypt(encKey, response.Payload)
	if err != nil {
		return nil, errors.Wrap(err, `decrypt payload`)
	}
	return convert.FromBytes(decrypted, target)
}
