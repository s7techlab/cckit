package router

import (
	"github.com/s7techlab/cckit/state"
)

// State interface for chain code CRUD operations
type State interface {
	Get(key interface{}, target interface{}) (result interface{}, err error)
	Exists(key interface{}) (exists bool, err error)
	Put(key interface{}, target interface{}) (err error)
	List(objectType string, target interface{}) (result []interface{}, err error)
}

type stateOp struct {
	context Context
}

func (s *stateOp) Get(key interface{}, target interface{}) (result interface{}, err error) {
	return state.Get(s.context.Stub(), key, target)
}

func (s *stateOp) Exists(key interface{}) (exists bool, err error) {
	return state.Exists(s.context.Stub(), key)
}

func (s *stateOp) Put(key interface{}, value interface{}) (err error) {
	return state.Put(s.context.Stub(), key, value)
}

func (s *stateOp) List(objectType string, target interface{}) (result []interface{}, err error) {
	return state.List(s.context.Stub(), objectType, target)
}
