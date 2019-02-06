package state

import (
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/convert"
)

// HistoryEntry struct containing history information of a single entry
type HistoryEntry struct {
	TxId      string      `json:"txId"`
	Timestamp int64       `json:"timestamp"`
	IsDeleted bool        `json:"isDeleted"`
	Value     interface{} `json:"value"`
}

// HistoryEntryList list of history entries
type HistoryEntryList []HistoryEntry

type (
	//KeyerFunc func(string) ([]string, error)
	KeyFunc func() ([]string, error)

	// Keyer interface for entity containing logic of its key creation
	Keyer interface {
		Key() ([]string, error)
	}

	KeyValue interface {
		Keyer
		convert.ToByter
	}
)

// State interface for chain code CRUD operations
type State interface {
	Get(entry interface{}, target ...interface{}) (result interface{}, err error)
	GetInt(entry interface{}, defaultValue int) (result int, err error)
	GetHistory(entry interface{}, target interface{}) (result HistoryEntryList, err error)
	Exists(entry interface{}) (exists bool, err error)
	Put(entry interface{}, value ...interface{}) (err error)
	Insert(entry interface{}, value ...interface{}) (err error)
	List(namespace interface{}, target ...interface{}) (result []interface{}, err error)
	Delete(entry interface{}) (err error)
}

type StateImpl struct {
	stub                shim.ChaincodeStubInterface
	KeyTransformer      KeyTransformer
	StateGetTransformer FromBytesTransformer
	StatePutTransformer ToBytesTransformer
}

// New creates wrapper on shim.ChaincodeStubInterface working with state
func New(stub shim.ChaincodeStubInterface) *StateImpl {
	return &StateImpl{
		stub:                stub,
		KeyTransformer:      ConvertKey,
		StateGetTransformer: ConvertFromBytes,
		StatePutTransformer: ConvertToBytes}
}

func (s *StateImpl) Key(key interface{}) (string, error) {
	keyParts, err := s.KeyTransformer(key)
	if err != nil {
		return ``, err
	}
	return KeyFromParts(s.stub, keyParts)
}

// Get data by key from state, trying to convert to target interface
func (s *StateImpl) Get(key interface{}, config ...interface{}) (result interface{}, err error) {
	strKey, err := s.Key(key)
	if err != nil {
		return nil, err
	}
	bb, err := s.stub.GetState(strKey)
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

	return s.StateGetTransformer(bb, config...)
}

func (s *StateImpl) GetInt(key interface{}, defaultValue int) (result int, err error) {
	val, err := s.Get(key, convert.TypeInt, defaultValue)
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}

// GetHistory by key from state, trying to convert to target interface
func (s *StateImpl) GetHistory(key interface{}, target interface{}) (result HistoryEntryList, err error) {
	strKey, err := s.Key(key)
	if err != nil {
		return nil, err
	}
	iter, err := s.stub.GetHistoryForKey(strKey)
	if err != nil {
		return nil, err
	}

	defer func() { _ = iter.Close() }()

	results := HistoryEntryList{}

	for iter.HasNext() {
		state, err := iter.Next()
		if err != nil {
			return nil, err
		}
		value, err := s.StateGetTransformer(state.Value, target)
		if err != nil {
			return nil, err
		}

		entry := HistoryEntry{
			TxId:      state.GetTxId(),
			Timestamp: state.GetTimestamp().GetSeconds(),
			IsDeleted: state.GetIsDelete(),
			Value:     value,
		}
		results = append(results, entry)
	}

	return results, nil
}

// Exists check entry with key exists in chaincode state
func (s *StateImpl) Exists(entry interface{}) (exists bool, err error) {
	stringKey, err := s.Key(entry)
	if err != nil {
		return false, errors.Wrap(err, `check key existence`)
	}
	bb, err := s.stub.GetState(stringKey)
	if err != nil {
		return false, err
	}
	return !(bb == nil || len(bb) == 0), nil
}

// List data from state using objectType prefix in composite key, trying to conver to target interface.
// Keys -  additional components of composite key
func (s *StateImpl) List(namespace interface{}, target ...interface{}) (result []interface{}, err error) {
	keyParts, err := s.KeyTransformer(namespace)
	if err != nil {
		return nil, errors.Wrap(err, `prepare list key parts`)
	}

	iter, err := s.stub.GetStateByPartialCompositeKey(keyParts[0], keyParts[1:])
	if err != nil {
		return nil, errors.Wrap(err, `create list iterator`)
	}

	entries := make([]interface{}, 0)
	defer func() { _ = iter.Close() }()

	for iter.HasNext() {
		v, err := iter.Next()
		if err != nil {
			return nil, err
		}
		entry, err := s.StateGetTransformer(v.Value, target...)
		if err != nil {
			return nil, errors.Wrap(err, `transform list entry`)
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *StateImpl) ArgKeyValue(arg interface{}, values []interface{}) (key interface{}, value interface{}, err error) {
	switch len(values) {

	// key is struct implementing keyer interface or has mapping instructions
	case 0:

		switch arg.(type) {

		case KeyValue:
			key, err = arg.(KeyValue).Key()
			if err != nil {
				return nil, nil, err
			}

			value, err := arg.(KeyValue).ToBytes()
			if err != nil {
				return nil, nil, err
			}

			return key, value, nil

		case Keyer:
			key, err = arg.(Keyer).Key()
			if err != nil {
				return nil, nil, err
			}
			return key, arg, err

		default:
			return nil, nil, ErrKeyNotSupportKeyerInterface
		}

	case 1:
		return arg, values[0], nil
	default:
		return nil, nil, ErrAllowOnlyOneValue
	}
}

// Put data value in state with key, trying convert data to []byte
func (s *StateImpl) Put(entry interface{}, values ...interface{}) (err error) {
	key, value, err := s.ArgKeyValue(entry, values)
	if err != nil {
		return err
	}

	bb, err := s.StatePutTransformer(value)
	if err != nil {
		return err
	}
	stringKey, err := s.Key(key)
	if err != nil {
		return err
	}
	return s.stub.PutState(stringKey, bb)
}

// Insert value into chaincode state, returns error if key already exists
func (s *StateImpl) Insert(entry interface{}, values ...interface{}) (err error) {
	key, value, err := s.ArgKeyValue(entry, values)
	if err != nil {
		return err
	}

	exists, err := s.Exists(key)
	if err != nil {
		return err
	}

	if exists {
		strKey, _ := s.Key(entry)
		return errors.Wrap(KeyError(strKey), ErrKeyAlreadyExists.Error())
	}

	return s.Put(key, value)
}

// Delete entry from state
func (s *StateImpl) Delete(key interface{}) (err error) {
	stringKey, err := s.Key(key)
	if err != nil {
		return errors.Wrap(err, `deleting from state`)
	}
	return s.stub.DelState(stringKey)
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

// KeyError error with key
func KeyError(strKey string) error {
	return errors.New(strings.Replace(strKey, "\x00", ` | `, -1))
}

//type stringKeyer struct {
//	str   string
//	keyer KeyerFunc
//}
//
//func (sk stringKeyer) Key() ([]string, error) {
//	return sk.keyer(sk.str)
//}
//
//// StringKeyer constructor for struct implementing Keyer interface
//func StringKeyer(str string, keyer KeyerFunc) Keyer {
//	return stringKeyer{str, keyer}
//}
