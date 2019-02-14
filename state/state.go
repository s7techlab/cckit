package state

import (
	"fmt"
	"reflect"
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
	Key []string

	//KeyerFunc func(string) ([]string, error)
	KeyFunc func() (Key, error)

	// Keyer interface for entity containing logic of its key creation
	Keyer interface {
		Key() (Key, error)
	}

	// StringsKeys interface for entity containing logic of its key creationn - backward compatibility
	StringsKeyer interface {
		Key() ([]string, error)
	}

	// KeyValue interface combines Keyer as ToByter methods - state entry representation
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

	Logger() *shim.ChaincodeLogger

	UseKeyTransformer(KeyTransformer) State
	UseStateGetTransformer(FromBytesTransformer) State
	UseStatePutTransformer(ToBytesTransformer) State
}

func (k Key) Append(key Key) Key {
	return append(k, key...)
}

type StateImpl struct {
	stub                shim.ChaincodeStubInterface
	logger              *shim.ChaincodeLogger
	StateKeyTransformer KeyTransformer
	StateGetTransformer FromBytesTransformer
	StatePutTransformer ToBytesTransformer
}

// NewState creates wrapper on shim.ChaincodeStubInterface for working with state
func NewState(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger) *StateImpl {
	return &StateImpl{
		stub:                stub,
		logger:              logger,
		StateKeyTransformer: KeyAsIs,
		StateGetTransformer: ConvertFromBytes,
		StatePutTransformer: ConvertToBytes,
	}
}

func (s *StateImpl) Logger() *shim.ChaincodeLogger {
	return s.logger
}

func (s *StateImpl) StringKey(key Key) (string, error) {
	keyParts, err := s.StateKeyTransformer(key)
	if err != nil {
		return ``, err
	}

	return StringKey(s.stub, keyParts)
}

func StringKey(stub shim.ChaincodeStubInterface, key Key) (string, error) {
	switch len(key) {
	case 0:
		return ``, ErrKeyPartsLength
	case 1:
		return key[0], nil
	default:
		return stub.CreateCompositeKey(key[0], key[1:])
	}
}

func (s *StateImpl) Key(key interface{}) (string, error) {
	// normalized key - []string
	normKey, err := NormalizeStateKey(key)
	if err != nil {
		return ``, errors.Wrap(err, `key normalizing`)
	}

	s.logger.Debugf(`state KEY: %s`, normKey)
	return s.StringKey(normKey)
}

// Get data by key from state, trying to convert to target interface
func (s *StateImpl) Get(key interface{}, config ...interface{}) (result interface{}, err error) {
	strKey, err := s.Key(key)
	if err != nil {
		return nil, err
	}

	//bytes from state
	s.logger.Debugf(`state GET %s`, strKey)
	bb, err := s.stub.GetState(strKey)
	if err != nil {
		return
	}
	if bb == nil || len(bb) == 0 {
		// config[1] default value
		if len(config) >= 2 {
			return config[1], nil
		}
		return nil, errors.Wrap(KeyError(strKey), ErrKeyNotFound.Error())
	}

	// config[0] - target type
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
	strKey, err := s.Key(entry)
	if err != nil {
		return false, err
	}
	s.logger.Debugf(`state check EXISTENCE %s`, strKey)
	bb, err := s.stub.GetState(strKey)
	if err != nil {
		return false, err
	}
	return !(bb == nil || len(bb) == 0), nil
}

// List data from state using objectType prefix in composite key, trying to conver to target interface.
// Keys -  additional components of composite key
func (s *StateImpl) List(namespace interface{}, target ...interface{}) (result []interface{}, err error) {

	key, err := NormalizeStateKey(namespace)
	if err != nil {
		return nil, errors.Wrap(err, `prepare list key parts`)
	}
	s.logger.Debugf(`state LIST namespace: %s`, key)

	key, err = s.StateKeyTransformer(key)
	if err != nil {
		return nil, err
	}
	s.logger.Debugf(`state LIST with composite key: %s`, key)

	iter, err := s.stub.GetStateByPartialCompositeKey(key[0], key[1:])
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

func NormalizeStateKey(key interface{}) (Key, error) {
	switch key.(type) {
	case Key:
		return key.(Key), nil
	case Keyer:
		return key.(Keyer).Key()
	case StringsKeyer:
		return key.(StringsKeyer).Key()
	case string:
		return Key{key.(string)}, nil
	case []string:
		return key.([]string), nil
	}
	return nil, fmt.Errorf(`%s: %s`, ErrUnableToCreateStateKey, reflect.TypeOf(key))
}

func (s *StateImpl) argKeyValue(arg interface{}, values []interface{}) (key Key, value interface{}, err error) {
	// key must be
	key, err = NormalizeStateKey(arg)
	if err != nil {
		return
	}

	switch len(values) {
	// arg is key and  value
	case 0:
		return key, arg, nil
	case 1:
		return key, values[0], nil
	default:
		return nil, nil, ErrAllowOnlyOneValue
	}
}

// Put data value in state with key, trying convert data to []byte
func (s *StateImpl) Put(entry interface{}, values ...interface{}) (err error) {
	key, value, err := s.argKeyValue(entry, values)
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

	s.logger.Debugf(`state PUT with string key: %s`, stringKey)
	return s.stub.PutState(stringKey, bb)
}

// Insert value into chaincode state, returns error if key already exists
func (s *StateImpl) Insert(entry interface{}, values ...interface{}) (err error) {
	if exists, err := s.Exists(entry); err != nil {
		return err
	} else if exists {
		strKey, _ := s.Key(entry)
		return errors.Wrap(KeyError(strKey), ErrKeyAlreadyExists.Error())
	}

	key, value, err := s.argKeyValue(entry, values)
	if err != nil {
		return err
	}
	return s.Put(key, value)
}

// Delete entry from state
func (s *StateImpl) Delete(key interface{}) (err error) {
	strKey, err := s.Key(key)
	if err != nil {
		return errors.Wrap(err, `deleting from state`)
	}

	s.logger.Debugf(`state DELETE with string key: %s`, strKey)
	return s.stub.DelState(strKey)
}

func (s *StateImpl) UseKeyTransformer(kt KeyTransformer) State {
	s.StateKeyTransformer = kt
	return s
}
func (s *StateImpl) UseStateGetTransformer(fb FromBytesTransformer) State {
	s.StateGetTransformer = fb
	return s
}

func (s *StateImpl) UseStatePutTransformer(tb ToBytesTransformer) State {
	s.StatePutTransformer = tb
	return s
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
