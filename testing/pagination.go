package testing

import (
	"container/list"
	"errors"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// MockStateRangeQueryPagedIterator represents paged version of shimtest.MockStateRangeQueryIterator
type MockStateRangeQueryPagedIterator struct {
	Closed       bool
	Stub         *MockStub
	Keys         *list.List
	Current      *list.Element
	nextBookmark string
}

func (iter *MockStateRangeQueryPagedIterator) Len() int32 {
	return int32(iter.Keys.Len())
}

func (iter *MockStateRangeQueryPagedIterator) NextBookmark() string {
	return iter.nextBookmark
}

func (iter *MockStateRangeQueryPagedIterator) Close() error {
	iter.Closed = true

	return nil
}

func (iter *MockStateRangeQueryPagedIterator) HasNext() bool {
	return iter.Current != nil
}

func (iter *MockStateRangeQueryPagedIterator) Next() (*queryresult.KV, error) {
	if iter.Closed {
		err := errors.New("MockStateRangeQueryPagedIterator.Next() called after Close()")
		return nil, err
	}

	if !iter.HasNext() {
		err := errors.New("MockStateRangeQueryPagedIterator.Next() called when it does not HaveNext()")
		return nil, err
	}

	key := iter.Current.Value.(string)
	value, err := iter.Stub.GetState(key)
	iter.Current = iter.Current.Next()

	return &queryresult.KV{Key: key, Value: value}, err
}

// NewMockStatesRangeQueryPagedIterator creates iterator starting from bookmark
// and limited by pageSize
func NewMockStatesRangeQueryPagedIterator(stub *MockStub, startKey string, endKey string, pageSize int32, bookmark string) *MockStateRangeQueryPagedIterator {
	iter := new(MockStateRangeQueryPagedIterator)
	iter.Stub = stub
	iter.Keys = new(list.List)

	var elem = stub.Keys.Front()
	// rewind until bookmark if is set
	for bookmark != "" && elem != nil {
		if elem.Value.(string) == bookmark {
			break
		}
		elem = elem.Next()
	}

	// Loop through keys until pageSize exceeded and find bookmark for next page
	for elem != nil {
		comp1 := strings.Compare(elem.Value.(string), startKey)
		comp2 := strings.Compare(elem.Value.(string), endKey)
		if (comp1 >= 0 && comp2 < 0) || (startKey == "" && endKey == "") {
			if iter.Keys.Len() < int(pageSize) {
				iter.Keys.PushBack(elem.Value)
				elem = elem.Next()

				continue
			}
			iter.nextBookmark = elem.Value.(string)
			break
		}
		elem = elem.Next()
	}

	iter.Current = iter.Keys.Front()

	return iter
}

// GetStateByPartialCompositeKeyWithPagination mocked
func (stub *MockStub) GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string,
	pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {

	partialCompositeKey, err := stub.CreateCompositeKey(objectType, keys)
	if err != nil {
		return nil, nil, err
	}

	iter := NewMockStatesRangeQueryPagedIterator(
		stub, partialCompositeKey, partialCompositeKey+string(maxUnicodeRuneValue), pageSize, bookmark)

	return iter, &pb.QueryResponseMetadata{
		FetchedRecordsCount: iter.Len(),
		Bookmark:            iter.NextBookmark(),
	}, nil
}

// GetStateByRangeWithPagination mocked
func (stub *MockStub) GetStateByRangeWithPagination(startKey, endKey string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	iter := NewMockStatesRangeQueryPagedIterator(stub, startKey, endKey, pageSize, bookmark)

	return iter, &pb.QueryResponseMetadata{
		FetchedRecordsCount: iter.Len(),
		Bookmark:            iter.NextBookmark(),
	}, nil
}
