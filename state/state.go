package state

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	pb "github.com/hyperledger/fabric-protos-go/peer"

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

type Impl struct {
	stub   shim.ChaincodeStubInterface
	logger *zap.Logger

	// wrappers for state access methods
	PutState                                    func(string, []byte) error
	GetState                                    func(string) ([]byte, error)
	DelState                                    func(string) error
	GetStateByPartialCompositeKey               func(objectType string, keys []string) (shim.StateQueryIteratorInterface, error)
	GetStateByPartialCompositeKeyWithPagination func(objectType string, keys []string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)

	StateKeyTransformer        KeyTransformer
	StateKeyReverseTransformer KeyTransformer
	StateGetTransformer        FromBytesTransformer
	StatePutTransformer        ToBytesTransformer
}

// NewState creates wrapper on shim.ChaincodeStubInterface for working with state
func NewState(stub shim.ChaincodeStubInterface, logger *zap.Logger) *Impl {
	i := &Impl{
		stub:                       stub,
		logger:                     logger,
		StateKeyTransformer:        KeyAsIs,
		StateKeyReverseTransformer: KeyAsIs,
		StateGetTransformer:        ConvertFromBytes,
		StatePutTransformer:        ConvertToBytes,
	}

	// Get data by key from state, direct from stub
	i.GetState = func(key string) ([]byte, error) {
		return stub.GetState(key)
	}

	// PutState puts the specified `key` and `value` into the transaction's
	// writeset as a data-write proposal.
	i.PutState = func(key string, bb []byte) error {
		return stub.PutState(key, bb)
	}

	// DelState records the specified `key` to be deleted in the writeset of
	// the transaction proposal.
	i.DelState = func(key string) error {
		return stub.DelState(key)
	}

	// GetStateByPartialCompositeKey queries the state in the ledger based on
	// a given partial composite key
	i.GetStateByPartialCompositeKey = func(objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
		return stub.GetStateByPartialCompositeKey(objectType, keys)
	}

	i.GetStateByPartialCompositeKeyWithPagination = func(
		objectType string, keys []string, pageSize int32, bookmark string) (
		shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
		return stub.GetStateByPartialCompositeKeyWithPagination(objectType, keys, pageSize, bookmark)
	}

	return i
}

func (s *Impl) Clone() State {
	return &Impl{
		stub:                          s.stub,
		logger:                        s.logger,
		PutState:                      s.PutState,
		GetState:                      s.GetState,
		DelState:                      s.DelState,
		GetStateByPartialCompositeKey: s.GetStateByPartialCompositeKey,
		GetStateByPartialCompositeKeyWithPagination: s.GetStateByPartialCompositeKeyWithPagination,
		StateKeyTransformer:                         s.StateKeyTransformer,
		StateKeyReverseTransformer:                  s.StateKeyReverseTransformer,
		StateGetTransformer:                         s.StateGetTransformer,
		StatePutTransformer:                         s.StatePutTransformer,
	}
}

func (s *Impl) Logger() *zap.Logger {
	return s.logger
}

func (s *Impl) Key(key interface{}) (*TransformedKey, error) {
	var (
		trKey = &TransformedKey{}
		err   error
	)

	if trKey.Origin, err = NormalizeKey(s.stub, key); err != nil {
		return nil, errors.Wrap(err, `key normalizing`)
	}

	s.logger.Debug(`state KEY`, zap.String(`key`, trKey.Origin.String()))

	if trKey.Parts, err = s.StateKeyTransformer(trKey.Origin); err != nil {
		return nil, err
	}

	if trKey.String, err = KeyToString(s.stub, trKey.Parts); err != nil {
		return nil, err
	}

	return trKey, nil
}

// Get data by key from state, trying to convert to target interface
func (s *Impl) Get(entry interface{}, config ...interface{}) (interface{}, error) {
	key, err := s.Key(entry)
	if err != nil {
		return nil, err
	}

	//bytes from state
	s.logger.Debug(`state GET`, zap.String(`key`, key.String))
	bb, err := s.GetState(key.String)
	if err != nil {
		return nil, err
	}
	if len(bb) == 0 {
		// config[1] default value
		if len(config) >= 2 {
			return config[1], nil
		}
		return nil, errors.Errorf(`%s: %s`, ErrKeyNotFound, key.Origin)
	}

	// config[0] - target type
	return s.StateGetTransformer(bb, config...)
}

func (s *Impl) GetInt(key interface{}, defaultValue int) (int, error) {
	val, err := s.Get(key, convert.TypeInt, defaultValue)
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}

// GetHistory by key from state, trying to convert to target interface
func (s *Impl) GetHistory(entry interface{}, target interface{}) (HistoryEntryList, error) {
	key, err := s.Key(entry)
	if err != nil {
		return nil, err
	}

	iter, err := s.stub.GetHistoryForKey(key.String)
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
func (s *Impl) Exists(entry interface{}) (bool, error) {
	key, err := s.Key(entry)
	if err != nil {
		return false, err
	}

	bb, err := s.GetState(key.String)
	if err != nil {
		return false, err
	}

	exists := len(bb) != 0
	s.logger.Debug(`state check EXISTENCE`, zap.String(`key`, key.String), zap.Bool(`exists`, exists))
	return exists, nil
}

// List data from state using objectType prefix in composite key, trying to convert to target interface.
// Keys -  additional components of composite key
func (s *Impl) List(namespace interface{}, target ...interface{}) (interface{}, error) {
	stateList, err := NewStateList(target...)
	if err != nil {
		return nil, err
	}

	iter, err := s.createStateQueryIterator(namespace)
	if err != nil {
		return nil, errors.Wrap(err, `state iterator`)
	}

	defer func() { _ = iter.Close() }()

	return stateList.Fill(iter, s.StateGetTransformer)
}

func (s *Impl) createStateQueryIterator(namespace interface{}) (shim.StateQueryIteratorInterface, error) {
	n, t, err := s.normalizeAndTransformKey(namespace)
	if err != nil {
		return nil, err
	}
	s.logger.Debug(`state KEYS with composite key`,
		zap.String(`key`, n.String()), zap.String(`transformed`, t.String()))

	objectType, attrs := t.Parts()
	if objectType == `` {
		return s.stub.GetStateByRange(``, ``) // all state entries
	}

	return s.GetStateByPartialCompositeKey(objectType, attrs)
}

// normalizeAndTransformKey returns normalized and transformed key or error if occur
func (s *Impl) normalizeAndTransformKey(namespace interface{}) (Key, Key, error) {
	normal, err := NormalizeKey(s.stub, namespace)
	if err != nil {
		return nil, nil, fmt.Errorf(`list prefix: %w`, err)
	}

	transformed, err := s.StateKeyTransformer(normal)
	if err != nil {
		return nil, nil, err
	}

	return normal, transformed, nil
}

func (s *Impl) ListPaginated(
	namespace interface{}, pageSize int32, bookmark string, target ...interface{}) (
	interface{}, *pb.QueryResponseMetadata, error) {
	stateList, err := NewStateList(target...)
	if err != nil {
		return nil, nil, err
	}

	iter, md, err := s.createStateQueryPagedIterator(namespace, pageSize, bookmark)
	if err != nil {
		return nil, nil, errors.Wrap(err, `state iterator`)
	}

	defer func() { _ = iter.Close() }()
	list, err := stateList.Fill(iter, s.StateGetTransformer)

	return list, md, err
}

func (s *Impl) createStateQueryPagedIterator(namespace interface{}, pageSize int32, bookmark string) (
	shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	n, t, err := s.normalizeAndTransformKey(namespace)
	if err != nil {
		return nil, nil, err
	}

	s.logger.Debug(`state KEYS with composite key`,
		zap.String(`key`, n.String()), zap.String(`transformed`, t.String()),
		zap.Int32("pageSize", pageSize), zap.String("bookmark", bookmark))

	objectType, attrs := t.Parts()
	if objectType == `` {
		return s.stub.GetStateByRangeWithPagination(``, ``, pageSize, bookmark)
	}

	return s.GetStateByPartialCompositeKeyWithPagination(objectType, attrs, pageSize, bookmark)
}

func (s *Impl) Keys(namespace interface{}) ([]string, error) {
	iter, err := s.createStateQueryIterator(namespace)
	if err != nil {
		return nil, errors.Wrap(err, `state iterator`)
	}

	defer func() { _ = iter.Close() }()

	var keys []string
	for iter.HasNext() {
		v, err := iter.Next()
		if err != nil {
			return nil, err
		}

		key, err := KeyFromComposite(s.stub, v.Key)
		if err != nil {
			return nil, err
		}

		reverseTransformedKey, err := s.StateKeyReverseTransformer(key)
		if err != nil {
			return nil, fmt.Errorf(`reverse transform key: %w`, err)
		}

		keyStr, err := KeyToString(s.stub, reverseTransformedKey)
		if err != nil {
			return nil, err
		}

		keys = append(keys, keyStr)
	}

	return keys, nil
}

func (s *Impl) argKeyValue(arg interface{}, values []interface{}) (key Key, value interface{}, err error) {
	// key must be
	if key, err = NormalizeKey(s.stub, arg); err != nil {
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

// Put data value in state with key, trying to convert data to []byte
func (s *Impl) Put(entry interface{}, values ...interface{}) error {
	entryKey, value, err := s.argKeyValue(entry, values)
	if err != nil {
		return err
	}
	bb, err := s.StatePutTransformer(value)
	if err != nil {
		return err
	}

	key, err := s.Key(entryKey)
	if err != nil {
		return err
	}

	s.logger.Debug(`state PUT`, zap.String(`key`, key.String))
	return s.PutState(key.String, bb)
}

// Insert value into chaincode state, returns error if key already exists
func (s *Impl) Insert(entry interface{}, values ...interface{}) error {
	if exists, err := s.Exists(entry); err != nil {
		return err
	} else if exists {
		key, _ := s.Key(entry)
		return fmt.Errorf(`%w: %s`, ErrKeyAlreadyExists, key.Origin)
	}

	key, value, err := s.argKeyValue(entry, values)
	if err != nil {
		return err
	}
	return s.Put(key, value)
}

// Delete entry from state
func (s *Impl) Delete(entry interface{}) error {
	key, err := s.Key(entry)
	if err != nil {
		return errors.Wrap(err, `deleting from state`)
	}

	s.logger.Debug(`state DELETE`, zap.String(`key`, key.String))
	return s.DelState(key.String)
}

func (s *Impl) UseKeyTransformer(kt KeyTransformer) {
	s.StateKeyTransformer = kt
}

func (s *Impl) UseKeyReverseTransformer(kt KeyTransformer) {
	s.StateKeyReverseTransformer = kt
}

func (s *Impl) UseStateGetTransformer(fb FromBytesTransformer) {
	s.StateGetTransformer = fb
}

func (s *Impl) UseStatePutTransformer(tb ToBytesTransformer) {
	s.StatePutTransformer = tb
}

// GetPrivate data by key from private state, trying to convert to target interface
func (s *Impl) GetPrivate(collection string, entry interface{}, config ...interface{}) (interface{}, error) {
	key, err := s.Key(entry)
	if err != nil {
		return nil, err
	}

	//bytes from private state
	s.logger.Debug(`private state GET`, zap.String(`key`, key.String))
	bb, err := s.stub.GetPrivateData(collection, key.String)
	if err != nil {
		return nil, err
	}
	if len(bb) == 0 {
		// config[1] default value
		if len(config) >= 2 {
			return config[1], nil
		}
		return nil, errors.Errorf(`%s: %s`, ErrKeyNotFound, key.Origin.String())
	}

	// config[0] - target type
	return s.StateGetTransformer(bb, config...)
}

// ExistsPrivate check entry with key exists in chaincode private state
func (s *Impl) ExistsPrivate(collection string, entry interface{}) (bool, error) {
	key, err := s.Key(entry)
	if err != nil {
		return false, err
	}
	s.logger.Debug(`private state check EXISTENCE`, zap.String(`key`, key.String))
	bb, err := s.stub.GetPrivateData(collection, key.String)
	if err != nil {
		return false, err
	}
	return len(bb) != 0, nil
}

// ListPrivate data from private state using objectType prefix in composite key, trying to convert to target interface.
// Keys -  additional components of composite key
// If usePrivateDataIterator is true, used private state for iterate over objects
// if false, used public state for iterate over keys and GetPrivateData for each key
func (s *Impl) ListPrivate(collection string, usePrivateDataIterator bool, namespace interface{}, target ...interface{}) (interface{}, error) {
	stateList, err := NewStateList(target...)
	if err != nil {
		return nil, err
	}
	key, err := NormalizeKey(s.stub, namespace)
	if err != nil {
		return nil, errors.Wrap(err, `prepare list key parts`)
	}
	s.logger.Debug(`state LIST`, zap.String(`namespace`, key.String()))

	if key, err = s.StateKeyTransformer(key); err != nil {
		return nil, err
	}
	s.logger.Debug(`state LIST with composite key`, zap.String(`namespace`, key.String()))

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

	var (
		kv       *queryresult.KV
		objKey   string
		keyParts []string
	)
	for iter.HasNext() {
		if kv, err = iter.Next(); err != nil {
			return nil, errors.Wrap(err, `get key value`)
		}
		if objKey, keyParts, err = s.stub.SplitCompositeKey(kv.Key); err != nil {
			return nil, err
		}

		object, err := s.GetPrivate(collection, append([]string{objKey}, keyParts...), target...)
		if err != nil {
			return nil, err
		}
		stateList.AddElementToList(object)
	}

	return stateList.Get()
}

// PutPrivate data value in private state with key, trying to convert data to []byte
func (s *Impl) PutPrivate(collection string, entry interface{}, values ...interface{}) (err error) {
	entryKey, value, err := s.argKeyValue(entry, values)
	if err != nil {
		return err
	}
	bb, err := s.StatePutTransformer(value)
	if err != nil {
		return err
	}

	key, err := s.Key(entryKey)
	if err != nil {
		return err
	}

	s.logger.Debug(`state PUT`, zap.String(`key`, key.String))
	return s.stub.PutPrivateData(collection, key.String, bb)
}

// InsertPrivate value into chaincode private state, returns error if key already exists
func (s *Impl) InsertPrivate(collection string, entry interface{}, values ...interface{}) (err error) {
	if exists, err := s.ExistsPrivate(collection, entry); err != nil {
		return err
	} else if exists {
		key, _ := s.Key(entry)
		return errors.Errorf(`%s: %s`, ErrKeyAlreadyExists, key.Origin)
	}

	key, value, err := s.argKeyValue(entry, values)
	if err != nil {
		return err
	}
	return s.PutPrivate(collection, key, value)
}

// DeletePrivate entry from private state
func (s *Impl) DeletePrivate(collection string, entry interface{}) error {
	key, err := s.Key(entry)
	if err != nil {
		return errors.Wrap(err, `deleting from private state`)
	}
	s.logger.Debug(`private state DELETE`, zap.String(`key`, key.String))
	return s.stub.DelPrivateData(collection, key.String)
}
