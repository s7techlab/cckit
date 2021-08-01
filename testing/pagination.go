package testing

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// MockStateRangeQueryPagedIterator represents paged version of shimtest.MockStateRangeQueryIterator
type MockStateRangeQueryPagedIterator struct {
	*shimtest.MockStateRangeQueryIterator

	position, size int
}

func (page *MockStateRangeQueryPagedIterator) HasNext() bool {
	if page.position >= page.size {
		return false
	}

	return page.MockStateRangeQueryIterator.HasNext()
}

func (page *MockStateRangeQueryPagedIterator) Next() (*queryresult.KV, error) {
	r, err := page.MockStateRangeQueryIterator.Next()
	if err == nil {
		page.position++
	}

	return r, err
}

// NewMockStatesRangeQueryPagedIterator creates MockStateRangeQueryIterator starting from bookmark
// and limited by pageSize
func NewMockStatesRangeQueryPagedIterator(stub *MockStub, startKey string, endKey string, pageSize int32, bookmark string) *MockStateRangeQueryPagedIterator {
	iter := new(MockStateRangeQueryPagedIterator)
	iter.MockStateRangeQueryIterator = shimtest.NewMockStateRangeQueryIterator(&stub.MockStub, startKey, endKey)
	iter.size = int(pageSize)

	// Forward iterator key to non empty bookmark
	for iter.Current != nil && bookmark != "" {
		iter.Current = iter.Current.Next()
		if iter.Current.Value.(string) == bookmark {
			break
		}
	}

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
		FetchedRecordsCount: pageSize,
		Bookmark:            bookmark,
	}, nil
}
