package testdata

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
)

const (
	TxStateCachedReadAfterWrite  = `ReadAfterWrite`
	TxStateCachedReadAfterDelete = `ReadAfterDelete`
	TxStateCachedListAfterWrite  = `ListAfterWrite`
	TxStateCachedListAfterDelete = `ListAfterDelete`

	BasePrefix = `prefix`
)

var (
	Keys = []string{`a`, `b`, `c`}
)

type Value struct {
	Value string
}

func Key(key string) []string {
	return []string{BasePrefix, key}
}

func KeyValue(key string) Value {
	return Value{Value: key + `_value`}
}
func NewStateCachedCC() *router.Chaincode {
	r := router.New(`state_cached`)

	r.Query(TxStateCachedReadAfterWrite, ReadAfterWrite).
		Query(TxStateCachedReadAfterDelete, ReadAfterDelete).
		Query(TxStateCachedListAfterWrite, ListAfterWrite).
		Query(TxStateCachedListAfterDelete, ListAfterDelete)

	return router.NewChaincode(r)
}

func ReadAfterWrite(ctx router.Context) (interface{}, error) {
	stateWithCache := state.WithCache(ctx.State())
	for _, k := range Keys {
		if err := stateWithCache.Put(Key(k), KeyValue(k)); err != nil {
			return nil, err
		}
	}

	// return non-empty, cause state changes cached
	return stateWithCache.Get(Key(Keys[0]), &Value{})
}

func ReadAfterDelete(ctx router.Context) (interface{}, error) {
	ctxWithStateCache := router.ContextWithStateCache(ctx)
	// delete all keys
	for _, k := range Keys {
		if err := ctxWithStateCache.State().Delete(Key(k)); err != nil {
			return nil, err
		}
	}

	// return empty, cause state changes cached
	val, _ := ctxWithStateCache.State().Get(Key(Keys[0]))
	// if we return error - state changes will not apply
	return val, nil
}

func ListAfterWrite(ctx router.Context) (interface{}, error) {
	stateWithCache := state.WithCache(ctx.State())
	for _, k := range Keys {
		if err := stateWithCache.Put(Key(k), KeyValue(k)); err != nil {
			return nil, err
		}
	}

	// return list, cause state changes cached
	return stateWithCache.List(BasePrefix, &Value{})
}

func ListAfterDelete(ctx router.Context) (interface{}, error) {
	ctxWithStateCache := router.ContextWithStateCache(ctx)
	// delete only one key, two keys remained
	if err := ctxWithStateCache.State().Delete(Key(Keys[0])); err != nil {
		return nil, err
	}

	// return list with 2 items, cause first item is deleted and state is cached
	return ctxWithStateCache.State().List(BasePrefix, &Value{})
}
