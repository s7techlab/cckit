package testdata

import "github.com/hyperledger/fabric-chaincode-go/shim"

func MustCreateCompositeKey(objectType string, attributes []string) string {
	key, err := shim.CreateCompositeKey(objectType, attributes)
	if err != nil {
		panic(err)
	}

	return key
}
