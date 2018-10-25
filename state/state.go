package state

import (
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/convert"
)

var (
	// ErrUnableToCreateKey can occurs while creating composite key for entry
	ErrUnableToCreateKey = errors.New(`unable to create state key`)

	// ErrKeyAlreadyExists can occurs when trying to insert entry with existing key
	ErrKeyAlreadyExists = errors.New(`state key already exists`)

	// ErrrKeyNotFound key not found in chaincode state
	ErrKeyNotFound = errors.New(`state entry not found`)

	// ErrAllowOnlyOneValue can occurs when trying to call Insert or Put with more than 2 arguments
	ErrAllowOnlyOneValue = errors.New(`allow only one value`)

	// ErrKeyNotSupportKeyerInterface can occurs when trying to Insert or Put struct without providing key and struct not support Keyer interface
	ErrKeyNotSupportKeyerInterface = errors.New(`key not support keyer interface`)

	// ErrKeyPartsLength can occurs when trying to create key consisting of zero parts
	ErrKeyPartsLength = errors.New(`key parts length must be greater than zero`)
)

// EntryList list of entries from state, gotten by part of composite key
type EntryList []interface{}

// Keyer interface for entity containing logic of its key creation
type Keyer interface {
	Key() ([]string, error)
}

type KeyerFunc func(string) ([]string, error)

// Get data by key from state, trying to convert to target interface
func Get(stub shim.ChaincodeStubInterface, key interface{}, config ...interface{}) (result interface{}, err error) {
	strKey, err := Key(stub, key)
	if err != nil {
		return false, err
	}
	bb, err := stub.GetState(strKey)
	if err != nil {
		return
	}
	if bb == nil || len(bb) == 0 {
		// default value
		if len(config) >= 2 {
			return config[1], nil
		}
		return nil, errors.Wrap(KeyError(strKey), ErrKeyNotFound.Error())
	}
	// converting to target type
	if len(config) >= 1 {
		return convert.FromBytes(bb, config[0])
	}

	// or return raw
	return bb, nil
}

// Exists check entry with key exists in chaincode state
func Exists(stub shim.ChaincodeStubInterface, key interface{}) (exists bool, err error) {
	stringKey, err := Key(stub, key)
	if err != nil {
		return false, errors.Wrap(err, `check key existence`)
	}
	bb, err := stub.GetState(stringKey)
	if err != nil {
		return false, err
	}
	return !(bb == nil || len(bb) == 0), nil
}

// List data from state using objectType prefix in composite key, trying to conver to target interface.
// Keys -  additional components of composite key
func List(stub shim.ChaincodeStubInterface, objectType interface{}, target interface{}) (result EntryList, err error) {
	keyParts, err := KeyParts(objectType)
	if err != nil {
		return nil, errors.Wrap(err, `unable to get key parts`)
	}
	iter, err := stub.GetStateByPartialCompositeKey(keyParts[0], keyParts[1:])
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

func getValue(key interface{}, values []interface{}) (interface{}, error) {
	switch len(values) {
	case 0:
		if _, ok := key.(Keyer); !ok {
			return nil, ErrKeyNotSupportKeyerInterface
		}
		return key, nil
	case 1:
		return values[0], nil
	default:
		return nil, ErrAllowOnlyOneValue
	}
}

// Put data value in state with key, trying convert data to []byte
func Put(stub shim.ChaincodeStubInterface, key interface{}, values ...interface{}) (err error) {
	value, err := getValue(key, values)
	if err != nil {
		return err
	}

	bb, err := convert.ToBytes(value)
	if err != nil {
		return err
	}
	stringKey, err := Key(stub, key)
	if err != nil {
		return err
	}
	return stub.PutState(stringKey, bb)
}

// Insert value into chaincode state, returns error if key already exists
func Insert(stub shim.ChaincodeStubInterface, key interface{}, values ...interface{}) (err error) {
	exists, err := Exists(stub, key)
	if err != nil {
		return err
	}

	if exists {
		strKey, _ := Key(stub, key)
		return errors.Wrap(KeyError(strKey), ErrKeyAlreadyExists.Error())
	}

	value, err := getValue(key, values)
	if err != nil {
		return err
	}

	return Put(stub, key, value)
}

// Delete entry from state
func Delete(stub shim.ChaincodeStubInterface, key interface{}) (err error) {
	stringKey, err := Key(stub, key)
	if err != nil {
		return errors.Wrap(err, `deleting from state`)
	}
	return stub.DelState(stringKey)
}

// Key transforms interface{} to string key
func Key(stub shim.ChaincodeStubInterface, key interface{}) (string, error) {
	keyParts, err := KeyParts(key)
	if err != nil {
		return ``, err
	}

	return KeyFromParts(stub, keyParts)
}

// KeyFromParts creates composite key by string slice
func KeyFromParts(stub shim.ChaincodeStubInterface, keyParts []string) (string, error) {
	switch len(keyParts) {
	case 0:
		return ``, ErrKeyPartsLength
	case 1:
		return keyParts[0], nil
	default:
		return stub.CreateCompositeKey(keyParts[0], keyParts[1:])
	}
}

// KeyParts returns string parts of composite key
func KeyParts(key interface{}) ([]string, error) {
	switch key.(type) {
	case Keyer:
		return key.(Keyer).Key()
	case string:
		return []string{key.(string)}, nil
	case []string:
		return key.([]string), nil
	}

	return nil, ErrUnableToCreateKey
}

// KeyError error with key
func KeyError(strKey string) error {
	return errors.New(strings.Replace(strKey, "\x00", ` | `, -1))
}

type stringKeyer struct {
	str   string
	keyer KeyerFunc
}

func (sk stringKeyer) Key() ([]string, error) {
	return sk.keyer(sk.str)
}

// StringKeyer constructor for struct implementing Keyer interface
func StringKeyer(str string, keyer KeyerFunc) Keyer {
	return stringKeyer{str, keyer}
}
