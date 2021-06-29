package state

import (
	"sort"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/pkg/errors"
)

var (
	ErrCachedQueryIteratorNoNext = errors.New(`cached query iterator no next`)
)

type (
	TxWriteSet  map[string][]byte
	TxDeleteSet map[string]interface{}

	Cached struct {
		State
		TxWriteSet  TxWriteSet
		TxDeleteSet TxDeleteSet
	}

	CachedQueryIterator struct {
		current int
		closed  bool
		KVs     []*queryresult.KV
	}
)

// WithCached returns state with tx level state cache
func WithCache(ss State) *Cached {
	s := ss.(*Impl)
	cached := &Cached{
		State:       s,
		TxWriteSet:  make(map[string][]byte),
		TxDeleteSet: make(map[string]interface{}),
	}

	// PutState wrapper
	s.PutState = func(key string, bb []byte) error {
		cached.TxWriteSet[key] = bb
		return s.stub.PutState(key, bb)
	}

	// GetState wrapper
	s.GetState = func(key string) ([]byte, error) {
		if bb, ok := cached.TxWriteSet[key]; ok {
			return bb, nil
		}

		if _, ok := cached.TxDeleteSet[key]; ok {
			return []byte{}, nil
		}
		return s.stub.GetState(key)
	}

	s.DelState = func(key string) error {
		delete(cached.TxWriteSet, key)
		cached.TxDeleteSet[key] = nil

		return s.stub.DelState(key)
	}

	s.GetStateByPartialCompositeKey = func(objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
		iterator, err := s.stub.GetStateByPartialCompositeKey(objectType, keys)
		if err != nil {
			return nil, err
		}

		prefix, err := s.stub.CreateCompositeKey(objectType, keys)
		if err != nil {
			return nil, err
		}

		return NewCachedQueryIterator(iterator, prefix, cached.TxWriteSet, cached.TxDeleteSet)
	}

	return cached
}

func NewCachedQueryIterator(iterator shim.StateQueryIteratorInterface, prefix string, writeSet TxWriteSet, deleteSet TxDeleteSet) (*CachedQueryIterator, error) {
	queryIterator := &CachedQueryIterator{
		current: -1,
	}
	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		if _, ok := deleteSet[kv.Key]; ok {
			continue
		}

		queryIterator.KVs = append(queryIterator.KVs, kv)
	}

	for wroteKey, wroteValue := range writeSet {
		if strings.HasPrefix(wroteKey, prefix) {
			queryIterator.KVs = append(queryIterator.KVs, &queryresult.KV{
				Namespace: "",
				Key:       wroteKey,
				Value:     wroteValue,
			})
		}
	}

	sort.Slice(queryIterator.KVs, func(i, j int) bool {
		return queryIterator.KVs[i].Key > queryIterator.KVs[i].Key
	})

	return queryIterator, nil
}

func (i *CachedQueryIterator) Next() (*queryresult.KV, error) {
	if !i.HasNext() {
		return nil, errors.New(`no next items`)
	}

	i.current++
	return i.KVs[i.current], nil
}

// HasNext returns true if the range query iterator contains additional keys
// and values.
func (i *CachedQueryIterator) HasNext() bool {
	return i.current < len(i.KVs)-1
}

// Close closes the iterator. This should be called when done
// reading from the iterator to free up resources.
func (i *CachedQueryIterator) Close() error {
	i.closed = true
	return nil
}
