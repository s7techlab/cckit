package state

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/convert"
)

var (
	// ErrUnableToCreateKey can occurs while creating composite key for entry
	ErrUnableToCreateKey = errors.New(`unable to create state key`)

	// ErrKeyAlreadyExists can occurs when trying to insert entry with existing key
	ErrKeyAlreadyExists = errors.New(`state key already exists`)
)

// EntryList list of entries from state, gotten by part of composite key
type EntryList []interface{}

// Get data by key from state, trying to convert to target interface
func Get(stub shim.ChaincodeStubInterface, key interface{}, target interface{}) (result interface{}, err error) {

	stringKey, err := Key(stub, key)
	if err != nil {
		return false, err
	}

	bb, err := stub.GetState(stringKey)
	if err != nil {
		return
	}
	if bb == nil || len(bb) == 0 {
		return nil, fmt.Errorf("state entry not found, key: %s", key)
	}
	return convert.FromBytes(bb, target)
}

// Exists check entry with key exists in chaincode state
func Exists(stub shim.ChaincodeStubInterface, key interface{}) (exists bool, err error) {

	stringKey, err := Key(stub, key)
	if err != nil {
		return false, err
	}

	bb, err := stub.GetState(stringKey)
	if err != nil {
		return false, err
	}
	return !(bb == nil || len(bb) == 0), nil
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
func Put(stub shim.ChaincodeStubInterface, key interface{}, value interface{}) (err error) {
	b, err := convert.ToBytes(value)
	if err != nil {
		return err
	}
	stringKey, err := Key(stub, key)
	if err != nil {
		return err
	}
	return stub.PutState(stringKey, b)
}

// Insert value into chaincode state, returns error if key already exists
func Insert(stub shim.ChaincodeStubInterface, key interface{}, value interface{}) (err error) {
	exists, err := Exists(stub, key)

	if err != nil {
		return err
	}

	if exists {
		return ErrKeyAlreadyExists
	}

	return Put(stub, key, value)
}

// Key transforms interface{} to string key
func Key(stub shim.ChaincodeStubInterface, key interface{}) (string, error) {
	switch key.(type) {
	case string:
		return key.(string), nil
	case []string:
		s := key.([]string)
		return stub.CreateCompositeKey(s[0], s[1:])
	}
	return ``, ErrUnableToCreateKey
}
