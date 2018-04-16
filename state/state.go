package state

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/vitiko/cckit/convert"
)

type  entryList []interface{}

func Get(stub shim.ChaincodeStubInterface, key string, target interface{}) (result interface{}, err error) {
	bb, err := stub.GetState(key)
	if err != nil {
		return
	}

	if bb == nil {
		return nil, fmt.Errorf("state entry not found, key: %s", key)
	}

	return convert.FromBytes(bb, target)
}

func List(stub shim.ChaincodeStubInterface, objectType string, target interface{}, keys ...string) (result entryList, err error) {

	iter, err := stub.GetStateByPartialCompositeKey(objectType, keys)
	if err != nil {
		return nil, err
	}

	entries := entryList{}
	defer iter.Close()
	for iter.HasNext() {
		v, err := iter.Next()
		if err != nil {
			return nil, err
		}
		entry, err := convert.FromBytes(v.Value, target)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func Put(stub shim.ChaincodeStubInterface, key string, value interface{}) (err error) {
	if b, err := convert.ToBytes(value); err != nil {
		return err
	} else {
		return stub.PutState(key, b)
	}
}
