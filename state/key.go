package state

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/s7techlab/cckit/convert"
)

// StringsIdToStr helper for passing []string key
func StringsIdToStr(idSlice []string) string {
	return strings.Join(idSlice, "\000")
}

// StringsIdFromStr helper for restoring []string key
func StringsIdFromStr(idString string) []string {
	return strings.Split(idString, "\000")
}

type (
	Key []string

	// StateKey stores origin and transformed state key
	TransformedKey struct {
		Origin Key
		Parts  Key
		String string
	}

	//KeyerFunc func(string) ([]string, error)
	KeyFunc func() (Key, error)

	// KeyerFunc transforms string to key
	KeyerFunc func(string) (Key, error)

	// Keyer interface for entity containing logic of its key creation
	Keyer interface {
		Key() (Key, error)
	}

	// StringsKeys interface for entity containing logic of its key creation - backward compatibility
	StringsKeyer interface {
		Key() ([]string, error)
	}

	// KeyValue interface combines Keyer as ToByter methods - state entry representation
	KeyValue interface {
		Keyer
		convert.ToByter
	}

	stringKeyer struct {
		str   string
		keyer KeyerFunc
	}
)

func (k Key) Append(key Key) Key {
	return append(k, key...)
}

// Key human readable representation
func (k Key) String() string {
	return strings.Join(k, ` | `)
}

// Parts returns object type and attributes slice
func (k Key) Parts() (objectType string, attrs []string) {
	if len(k) > 0 {
		objectType = k[0]

		if len(k) > 1 {
			attrs = k[1:]
		}
	}
	return
}

func NormalizeKey(stub shim.ChaincodeStubInterface, key interface{}) (Key, error) {
	switch k := key.(type) {
	case Key:
		return k, nil
	case Keyer:
		return k.Key()
	case StringsKeyer:
		return k.Key()
	case string:
		return KeyFromComposite(stub, k)
	case []string:
		return k, nil
	}
	return nil, fmt.Errorf(`%s: %w`, reflect.TypeOf(key), ErrUnableToCreateStateKey)
}

func KeyFromComposite(stub shim.ChaincodeStubInterface, key string) (Key, error) {
	var (
		objectType string
		attributes []string
		err        error
	)

	// contains key delimiter
	if strings.ContainsRune(key, 0) {
		objectType, attributes, err = stub.SplitCompositeKey(key)
		if err != nil {
			return nil, fmt.Errorf(`key from composite: %w`, err)
		}
	} else {
		objectType = key
	}

	return append([]string{objectType}, attributes...), nil
}

func KeyToComposite(stub shim.ChaincodeStubInterface, key Key) (string, error) {
	compositeKey, err := stub.CreateCompositeKey(key[0], key[1:])
	if err != nil {
		return ``, fmt.Errorf(`key to composite: %w`, err)
	}

	return compositeKey, nil
}

func KeyToString(stub shim.ChaincodeStubInterface, key Key) (string, error) {
	switch len(key) {
	case 0:
		return ``, ErrKeyPartsLength
	case 1:
		return key[0], nil
	default:
		return KeyToComposite(stub, key)
	}
}

func (sk stringKeyer) Key() (Key, error) {
	return sk.keyer(sk.str)
}

// StringKeyer constructor for struct implementing Keyer interface
func StringKeyer(str string, keyer KeyerFunc) Keyer {
	return stringKeyer{str, keyer}
}
