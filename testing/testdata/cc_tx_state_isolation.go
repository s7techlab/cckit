package testdata

import (
	"github.com/s7techlab/cckit/router"
)

const (
	Key1 = `abc`

	TxIsolationReadAfterWrite  = `ReadAfterWrite`
	TxIsolationReadAfterDelete = `ReadAfterDelete`
)

var (
	Value1 = []byte(`cde`)
)

func NewTxIsolationCC() *router.Chaincode {
	r := router.New(`tx_isolation`)

	r.Query(TxIsolationReadAfterWrite, ReadAfterWrite).
		Query(TxIsolationReadAfterDelete, ReadAfterDelete)

	return router.NewChaincode(r)
}

func ReadAfterWrite(c router.Context) (interface{}, error) {
	if err := c.State().Put(Key1, Value1); err != nil {
		return nil, err
	}

	// return empty, cause state changes cannot be read
	return c.State().Get(Key1)
}

func ReadAfterDelete(c router.Context) (interface{}, error) {
	if err := c.State().Delete(Key1); err != nil {
		return nil, err
	}

	// return non empty, cause state changes, include deletion, cannot be read
	return c.State().Get(Key1)
}
