package debug

import (
	"fmt"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/state"
)

const (
	InvokeStateCleanFunc  = `StateClean`
	QueryStateKeysFunc    = `StateKeys`
	QueryStateGetFunc     = `StateGet`
	InvokeStatePutFunc    = `StatePut`
	InvokeStateDeleteFunc = `StateDelete`
)

var (
	// KeyParam parameter for get, put, delete data from state
	KeyParam = param.Strings(`key`)

	// PrefixParam parameter
	PrefixParam = param.String(`prefix`)

	// PrefixesParam parameter
	PrefixesParam = param.Strings(`prefixes`)

	// ValueParam  parameter for putting value in state
	ValueParam = param.Bytes(`value`)
)

// AddHandlers adds debug handlers to router, allows to add more middleware
// for example for access control
func AddHandlers(r *router.Group, prefix string, middleware ...router.MiddlewareFunc) {

	// clear state entries by key prefix
	r.Invoke(
		prefix+InvokeStateCleanFunc,
		InvokeStateClean,
		append([]router.MiddlewareFunc{PrefixesParam}, middleware...)...)

	// query keys by prefix
	r.Query(
		prefix+QueryStateKeysFunc,
		QueryKeysList,
		append([]router.MiddlewareFunc{PrefixParam}, middleware...)...)

	// query value by key
	r.Query(
		prefix+QueryStateGetFunc,
		QueryStateGet,
		append([]router.MiddlewareFunc{KeyParam}, middleware...)...)

	r.Invoke(
		prefix+InvokeStatePutFunc,
		InvokeStatePut,
		append([]router.MiddlewareFunc{KeyParam, ValueParam}, middleware...)...)

	r.Invoke(
		prefix+InvokeStateDeleteFunc,
		InvokeStateDelete,
		append([]router.MiddlewareFunc{KeyParam}, middleware...)...)
}

// InvokeStateClean delete entries from state, prefix []string contains key prefixes or whole key
func InvokeStateClean(c router.Context) (interface{}, error) {
	return DeleteStateByPrefixes(c.State(), c.Param(`prefixes`).([]string))
}

// InvokeStatePut router handler puts value in chaincode state with composite key,
// created with key parts ([]string)
func InvokeStatePut(c router.Context) (interface{}, error) {
	key, err := state.KeyToString(c.Stub(), c.Param(`key`).([]string))
	if err != nil {
		return nil, fmt.Errorf(`create key: %w`, err)
	}
	return nil, c.Stub().PutState(key, c.ParamBytes(`value`))
}

// QueryKeysList router handler returns string slice with keys by prefix (object type)
func QueryKeysList(c router.Context) (interface{}, error) {
	return c.State().Keys(c.ParamString(`prefix`))
}

// QueryStateGet router handler returns state entry by key ([]string)
func QueryStateGet(c router.Context) (interface{}, error) {
	key, err := state.KeyToString(c.Stub(), c.Param(`key`).([]string))
	if err != nil {
		return nil, fmt.Errorf(`create key: %w`, err)
	}
	return c.Stub().GetState(key)
}

// InvokeStateDelete router handler delete state entry by key ([]string)
func InvokeStateDelete(c router.Context) (interface{}, error) {
	key, err := state.KeyToString(c.Stub(), c.Param(`key`).([]string))
	if err != nil {
		return nil, fmt.Errorf(`create key: %w`, err)
	}
	return nil, c.Stub().DelState(key)
}
