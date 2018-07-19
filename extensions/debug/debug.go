package debug

import (
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/state"
)

const (
	InvokeStateCleanFunc  = `StateClean`
	QueryStateKeysFunc    = `StateKeysList`
	QueryStateGetFunc     = `StateGet`
	InvokeStatePutFunc    = `StatePut`
	InvokeStateDeleteFunc = `StateDelete`
)

// AddHandler adds debug handlers to router
func AddHandlers(r *router.Group, prefix string, middleware ...router.MiddlewareFunc) {
	r.Invoke(prefix+InvokeStateCleanFunc, InvokeStateClean,
		append([]router.MiddlewareFunc{param.Strings(`prefix`)}, middleware...)...)
	r.Query(prefix+QueryStateKeysFunc, QueryKeysList,
		append([]router.MiddlewareFunc{param.Strings(`prefix`)}, middleware...)...)
	r.Query(prefix+QueryStateGetFunc, QueryStateGet,
		append([]router.MiddlewareFunc{param.Strings(`key`)}, middleware...)...)
	r.Invoke(prefix+InvokeStatePutFunc, InvokeStatePut,
		append([]router.MiddlewareFunc{param.Strings(`key`), param.Bytes(`value`, 1)}, middleware...)...)
	r.Invoke(prefix+InvokeStateDeleteFunc, InvokeStateDelete,
		append([]router.MiddlewareFunc{param.Strings(`key`)}, middleware...)...)
}

// InvokeStateClean delete entries from state, prefix []string contains key prefixes or whole key
func InvokeStateClean(c router.Context) (interface{}, error) {
	prefixMatches := make(map[string]int)
	for _, prefix := range c.Arg(`prefix`).([]string) {
		iter, err := c.Stub().GetStateByPartialCompositeKey(prefix, []string{})
		prefixMatches[prefix] = 0

		if err != nil {
			return nil, err
		}
		for iter.HasNext() {
			v, err := iter.Next()
			if err != nil {
				return nil, err
			}

			if err := c.Stub().DelState(v.Key); err != nil {
				return nil, err
			}
			prefixMatches[prefix]++
		}
		iter.Close()
	}
	return prefixMatches, nil
}

// InvokeValueByKeyPut router handler puts value in chaincode state with composite key, created with key parts ([]string)
func InvokeStatePut(c router.Context) (interface{}, error) {
	key, err := state.KeyFromParts(c.Stub(), c.Arg(`key`).([]string))
	if err != nil {
		return nil, errors.Wrap(err, `unable to create key`)
	}
	return nil, c.Stub().PutState(key, c.ArgBytes(`value`))
}

// QueryKeysList router handler returns string slice with keys by prefix (object type)
func QueryKeysList(c router.Context) (interface{}, error) {
	prefixes := c.Arg(`prefix`).([]string)
	iter, err := c.Stub().GetStateByPartialCompositeKey(prefixes[0], prefixes[1:])
	if err != nil {
		return nil, err
	}

	defer func() { _ = iter.Close() }()
	var keys []string
	for iter.HasNext() {
		v, err := iter.Next()
		if err != nil {
			return nil, err
		}

		keys = append(keys, v.Key)
	}
	return keys, nil
}

// QueryStateGet router handler returns state entry by key ([]string)
func QueryStateGet(c router.Context) (interface{}, error) {
	key, err := state.KeyFromParts(c.Stub(), c.Arg(`key`).([]string))
	if err != nil {
		return nil, errors.Wrap(err, `unable to create key`)
	}
	return c.Stub().GetState(key)
}

// QueryStateGet router handler delete state entry by key ([]string)
func InvokeStateDelete(c router.Context) (interface{}, error) {
	key, err := state.KeyFromParts(c.Stub(), c.Arg(`key`).([]string))
	if err != nil {
		return nil, errors.Wrap(err, `unable to create key`)
	}
	return nil, c.Stub().DelState(key)
}
