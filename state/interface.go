package state

import (
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"go.uber.org/zap"
)

// State interface for chain code CRUD operations
type State interface {
	Gettable
	Settable
	Listable
	ListablePaginated
	Deletable
	Historyable
	Transformable
	Privateable

	// Deprecated: GetInt returns value from state, converted to int
	// entry can be Key (string or []string) or type implementing Keyer interface
	GetInt(entry interface{}, defaultValue int) (int, error)

	// Keys returns slice of keys
	// namespace can be part of key (string or []string) or entity with defined mapping
	Keys(namespace interface{}) ([]string, error)

	Logger() *zap.Logger

	// Clone state for next changing transformers, state access methods etc
	Clone() State
}

type GetSettable interface {
	Gettable
	Settable
}

type (
	Gettable interface {
		// Get returns value from state, converted to target type
		// entry can be Key (string or []string) or type implementing Keyer interface
		Get(entry interface{}, target ...interface{}) (interface{}, error)
		// Exists returns entry existence in state
		// entry can be Key (string or []string) or type implementing Keyer interface
		Exists(entry interface{}) (bool, error)
	}

	Settable interface {
		// Put returns result of putting entry to state
		// entry can be Key (string or []string) or type implementing Keyer interface
		// if entry is implements Keyer interface, and it's struct or type implementing
		// ToByter interface value can be omitted
		Put(entry interface{}, value ...interface{}) error
		// Insert returns result of inserting entry to state
		// If same key exists in state error wil be returned
		// entry can be Key (string or []string) or type implementing Keyer interface
		// if entry is implements Keyer interface, and it's struct or type implementing
		// ToByter interface value can be omitted
		Insert(entry interface{}, value ...interface{}) error
	}

	Listable interface {
		// List returns slice of target type
		// namespace can be part of key (string or []string) or entity with defined mapping
		List(namespace interface{}, target ...interface{}) (interface{}, error)
	}

	ListablePaginated interface {
		// ListPaginated returns slice of target type with pagination
		// namespace can be part of key (string or []string) or entity with defined mapping
		ListPaginated(namespace interface{}, pageSize int32, bookmark string, target ...interface{}) (
			interface{}, *pb.QueryResponseMetadata, error)
	}

	Deletable interface {
		// Delete returns result of deleting entry from state
		// entry can be Key (string or []string) or type implementing Keyer interface
		Delete(entry interface{}) (err error)
	}

	Historyable interface {
		// GetHistory returns slice of history records for entry, with values converted to target type
		// entry can be Key (string or []string) or type implementing Keyer interface
		GetHistory(entry interface{}, target interface{}) (HistoryEntryList, error)
	}

	Privateable interface {
		// GetPrivate returns value from private state, converted to target type
		// entry can be Key (string or []string) or type implementing Keyer interface
		GetPrivate(collection string, entry interface{}, target ...interface{}) (interface{}, error)

		// PutPrivate returns result of putting entry to private state
		// entry can be Key (string or []string) or type implementing Keyer interface
		// if entry is implements Keyer interface, and it's struct or type implementing
		// ToByter interface value can be omitted
		PutPrivate(collection string, entry interface{}, value ...interface{}) error

		// InsertPrivate returns result of inserting entry to private state
		// If same key exists in state error wil be returned
		// entry can be Key (string or []string) or type implementing Keyer interface
		// if entry is implements Keyer interface, and it's struct or type implementing
		// ToByter interface value can be omitted
		InsertPrivate(collection string, entry interface{}, value ...interface{}) error

		// ListPrivate returns slice of target type from private state
		// namespace can be part of key (string or []string) or entity with defined mapping
		// If usePrivateDataIterator is true, used private state for iterate over objects
		// if false, used public state for iterate over keys and GetPrivateData for each key
		ListPrivate(collection string, usePrivateDataIterator bool, namespace interface{}, target ...interface{}) (interface{}, error)

		// DeletePrivate returns result of deleting entry from private state
		// entry can be Key (string or []string) or type implementing Keyer interface
		DeletePrivate(collection string, entry interface{}) error

		// ExistsPrivate returns entry existence in private state
		// entry can be Key (string or []string) or type implementing Keyer interface
		ExistsPrivate(collection string, entry interface{}) (bool, error)
	}

	Transformable interface {
		UseKeyTransformer(KeyTransformer)
		UseKeyReverseTransformer(KeyTransformer)
		UseStateGetTransformer(FromBytesTransformer)
		UseStatePutTransformer(ToBytesTransformer)
	}
)
