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

	// StringsKeys interface for entity containing logic of its key creation - backward compatibility
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
	// Get returns value from state, converted to target type
	// entry can be Key (string or []string) or type implementing Keyer interface
	Get(entry interface{}, target ...interface{}) (result interface{}, err error)

	// Get returns value from state, converted to int
	// entry can be Key (string or []string) or type implementing Keyer interface
	GetInt(entry interface{}, defaultValue int) (result int, err error)

	// GetHistory returns slice of history records for entry, with values converted to target type
	// entry can be Key (string or []string) or type implementing Keyer interface
	GetHistory(entry interface{}, target interface{}) (result HistoryEntryList, err error)

	// Exists returns entry existence in state
	// entry can be Key (string or []string) or type implementing Keyer interface
	Exists(entry interface{}) (exists bool, err error)

	// Put returns result of putting entry to state
	// entry can be Key (string or []string) or type implementing Keyer interface
	// if entry is implements Keyer interface and it's struct or type implementing
	// ToByter interface value can be omitted
	Put(entry interface{}, value ...interface{}) (err error)

	// Insert returns result of inserting entry to state
	// If same key exists in state error wil be returned
	// entry can be Key (string or []string) or type implementing Keyer interface
	// if entry is implements Keyer interface and it's struct or type implementing
	// ToByter interface value can be omitted
	Insert(entry interface{}, value ...interface{}) (err error)

	// List returns slice of target type
	// namespace can be part of key (string or []string) or entity with defined mapping
	List(namespace interface{}, target ...interface{}) (result interface{}, err error)

	// Delete returns result of deleting entry from state
	// entry can be Key (string or []string) or type implementing Keyer interface
	Delete(entry interface{}) (err error)

	Logger() *shim.ChaincodeLogger

	UseKeyTransformer(KeyTransformer) State
	UseStateGetTransformer(FromBytesTransformer) State
	UseStatePutTransformer(ToBytesTransformer) State

	// GetPrivate returns value from private state, converted to target type
	// entry can be Key (string or []string) or type implementing Keyer interface
	GetPrivate(collection string, entry interface{}, target ...interface{}) (result interface{}, err error)

	// PutPrivate returns result of putting entry to private state
	// entry can be Key (string or []string) or type implementing Keyer interface
	// if entry is implements Keyer interface and it's struct or type implementing
	// ToByter interface value can be omitted
	PutPrivate(collection string, entry interface{}, value ...interface{}) (err error)

	// InsertPrivate returns result of inserting entry to private state
	// If same key exists in state error wil be returned
	// entry can be Key (string or []string) or type implementing Keyer interface
	// if entry is implements Keyer interface and it's struct or type implementing
	// ToByter interface value can be omitted
	InsertPrivate(collection string, entry interface{}, value ...interface{}) (err error)

	// ListPrivate returns slice of target type from private state
	// namespace can be part of key (string or []string) or entity with defined mapping
	// If usePrivateDataIterator is true, used private state for iterate over objects
	// if false, used public state for iterate over keys and GetPrivateData for each key
	ListPrivate(collection string, usePrivateDataIterator bool, namespace interface{}, target ...interface{}) (result interface{}, err error)

	// DeletePrivate returns result of deleting entry from private state
	// entry can be Key (string or []string) or type implementing Keyer interface
	DeletePrivate(collection string, entry interface{}) (err error)

	// ExistsPrivate returns entry existence in private state
	// entry can be Key (string or []string) or type implementing Keyer interface
	ExistsPrivate(collection string, entry interface{}) (exists bool, err error)
}

func (k Key) Append(key Key) Key {
	return append(k, key...)
}

type Impl struct {
	stub                shim.ChaincodeStubInterface
	logger              *shim.ChaincodeLogger
	StateKeyTransformer KeyTransformer
	StateGetTransformer FromBytesTransformer
	StatePutTransformer ToBytesTransformer
}

// NewState creates wrapper on shim.ChaincodeStubInterface for working with state
func NewState(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger) *Impl {
	return &Impl{
		stub:                stub,
		logger:              logger,
		StateKeyTransformer: KeyAsIs,
		StateGetTransformer: ConvertFromBytes,
		StatePutTransformer: ConvertToBytes,
	}
}

func (s *Impl) Logger() *shim.ChaincodeLogger {
	return s.logger
}

func (s *Impl) StringKey(key Key) (string, error) {
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

func (s *Impl) Key(key interface{}) (string, error) {
	// normalized key - []string
	normKey, err := NormalizeStateKey(key)
	if err != nil {
		return ``, errors.Wrap(err, `key normalizing`)
	}

	s.logger.Debugf(`state KEY: %s`, normKey)
	return s.StringKey(normKey)
}

// Get data by key from state, trying to convert to target interface
func (s *Impl) Get(key interface{}, config ...interface{}) (result interface{}, err error) {
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
	if len(bb) == 0 {
		// config[1] default value
		if len(config) >= 2 {
			return config[1], nil
		}
		return nil, errors.Wrap(KeyError(strKey), ErrKeyNotFound.Error())
	}

	// config[0] - target type
	return s.StateGetTransformer(bb, config...)
}

func (s *Impl) GetInt(key interface{}, defaultValue int) (result int, err error) {
	val, err := s.Get(key, convert.TypeInt, defaultValue)
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}

// GetHistory by key from state, trying to convert to target interface
func (s *Impl) GetHistory(key interface{}, target interface{}) (result HistoryEntryList, err error) {
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
func (s *Impl) Exists(entry interface{}) (exists bool, err error) {
	strKey, err := s.Key(entry)
	if err != nil {
		return false, err
	}
	s.logger.Debugf(`state check EXISTENCE %s`, strKey)
	bb, err := s.stub.GetState(strKey)
	if err != nil {
		return false, err
	}
	return len(bb) != 0, nil
}

// List data from state using objectType prefix in composite key, trying to convert to target interface.
// Keys -  additional components of composite key
func (s *Impl) List(namespace interface{}, target ...interface{}) (result interface{}, err error) {

	stateList, err := NewStateList(target...)
	if err != nil {
		return nil, err
	}
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
	defer func() { _ = iter.Close() }()

	return stateList.Fill(iter, s.StateGetTransformer)
}

func NormalizeStateKey(key interface{}) (Key, error) {
	switch k := key.(type) {
	case Key:
		return k, nil
	case Keyer:
		return k.Key()
	case StringsKeyer:
		return k.Key()
	case string:
		return Key{k}, nil
	case []string:
		return k, nil
	}
	return nil, fmt.Errorf(`%s: %s`, ErrUnableToCreateStateKey, reflect.TypeOf(key))
}

func (s *Impl) argKeyValue(arg interface{}, values []interface{}) (key Key, value interface{}, err error) {
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
func (s *Impl) Put(entry interface{}, values ...interface{}) (err error) {
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
func (s *Impl) Insert(entry interface{}, values ...interface{}) (err error) {
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
func (s *Impl) Delete(key interface{}) (err error) {
	strKey, err := s.Key(key)
	if err != nil {
		return errors.Wrap(err, `deleting from state`)
	}

	s.logger.Debugf(`state DELETE with string key: %s`, strKey)
	return s.stub.DelState(strKey)
}

func (s *Impl) UseKeyTransformer(kt KeyTransformer) State {
	s.StateKeyTransformer = kt
	return s
}
func (s *Impl) UseStateGetTransformer(fb FromBytesTransformer) State {
	s.StateGetTransformer = fb
	return s
}

func (s *Impl) UseStatePutTransformer(tb ToBytesTransformer) State {
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

// Get data by key from private state, trying to convert to target interface
func (s *Impl) GetPrivate(collection string, key interface{}, config ...interface{}) (result interface{}, err error) {
	strKey, err := s.Key(key)
	if err != nil {
		return nil, err
	}

	//bytes from private state
	s.logger.Debugf(`private state GET %s`, strKey)
	bb, err := s.stub.GetPrivateData(collection, strKey)
	if err != nil {
		return
	}
	if len(bb) == 0 {
		// config[1] default value
		if len(config) >= 2 {
			return config[1], nil
		}
		return nil, errors.Wrap(KeyError(strKey), ErrKeyNotFound.Error())
	}

	// config[0] - target type
	return s.StateGetTransformer(bb, config...)
}

// PrivateExists check entry with key exists in chaincode private state
func (s *Impl) ExistsPrivate(collection string, entry interface{}) (exists bool, err error) {
	strKey, err := s.Key(entry)
	if err != nil {
		return false, err
	}
	s.logger.Debugf(`private state check EXISTENCE %s`, strKey)
	bb, err := s.stub.GetPrivateData(collection, strKey)
	if err != nil {
		return false, err
	}
	return len(bb) != 0, nil
}

// List data from private state using objectType prefix in composite key, trying to convert to target interface.
// Keys -  additional components of composite key
// If usePrivateDataIterator is true, used private state for iterate over objects
// if false, used public state for iterate over keys and GetPrivateData for each key
func (s *Impl) ListPrivate(collection string, usePrivateDataIterator bool, namespace interface{}, target ...interface{}) (result interface{}, err error) {

	stateList, err := NewStateList(target...)
	if err != nil {
		return nil, err
	}
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

	if usePrivateDataIterator {
		iter, err := s.stub.GetPrivateDataByPartialCompositeKey(collection, key[0], key[1:])
		if err != nil {
			return nil, errors.Wrap(err, `create list iterator`)
		}
		defer func() { _ = iter.Close() }()
		return stateList.Fill(iter, s.StateGetTransformer)
	}

	iter, err := s.stub.GetStateByPartialCompositeKey(key[0], key[1:])
	if err != nil {
		return nil, errors.Wrap(err, `create list iterator`)
	}
	defer func() { _ = iter.Close() }()
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return nil, errors.Wrap(err, `get key value`)
		}
		objKey, keyParts, err := s.stub.SplitCompositeKey(kv.Key)
		var curCompositeKey []string
		curCompositeKey = append(curCompositeKey, objKey)
		for _, part := range keyParts {
			curCompositeKey = append(curCompositeKey, part)
		}
		object, err := s.GetPrivate(collection, curCompositeKey, target...)
		if err != nil {
			return nil, err
		}
		stateList.AddElementToList(object)
	}

	return stateList.Get()
}

// Put data value in private state with key, trying convert data to []byte
func (s *Impl) PutPrivate(collection string, entry interface{}, values ...interface{}) (err error) {
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
	return s.stub.PutPrivateData(collection, stringKey, bb)
}

// Insert value into chaincode private state, returns error if key already exists
func (s *Impl) InsertPrivate(collection string, entry interface{}, values ...interface{}) (err error) {
	if exists, err := s.ExistsPrivate(collection, entry); err != nil {
		return err
	} else if exists {
		strKey, _ := s.Key(entry)
		return errors.Wrap(KeyError(strKey), ErrKeyAlreadyExists.Error())
	}

	key, value, err := s.argKeyValue(entry, values)
	if err != nil {
		return err
	}
	return s.PutPrivate(collection, key, value)
}

// Delete entry from private state
func (s *Impl) DeletePrivate(collection string, key interface{}) (err error) {
	strKey, err := s.Key(key)
	if err != nil {
		return errors.Wrap(err, `deleting from private state`)
	}
	s.logger.Debugf(`private state DELETE with string key: %s`, strKey)
	return s.stub.DelPrivateData(collection, strKey)
}
