package state

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/convert"
)

// EntryList list of entries from state, gotten by part of composite key
type EntryList []interface{}

// Get data by key from state, trying to convert to target interface
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

// List data from state using objectType prefix in composite key, trying to conver to target interface.
// Keys -  additional components of composite key
func List(stub shim.ChaincodeStubInterface, objectType string, target interface{}, keys ...string) (result EntryList, err error) {
	iter, err := stub.GetStateByPartialCompositeKey(objectType, keys)
	if err != nil {
		return nil, err
	}

	entries := EntryList{}
	defer func() { _ = iter.Close() }()

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

// Put data value in state with key, trying convert data to []byte
func Put(stub shim.ChaincodeStubInterface, key string, value interface{}) (err error) {
	b, err := convert.ToBytes(value)
	if err != nil {
		return err
	}
	return stub.PutState(key, b)
}
