package cckit

import (
	"github.com/s7techlab/cckit/state"
)

type State interface {
	Get(key string, target interface{}) (result interface{}, err error)
	Put(key string, target interface{}) (err error)
	List(objectType string, target interface{}) (result []interface{}, err error)
}

type stateOp struct {
	context Context
}

func (s *stateOp) Get(key string, target interface{}) (result interface{}, err error) {
	return state.Get(s.context.Stub(), key, target)
}

func (s *stateOp) Put(key string, value interface{}) (err error) {
	return state.Put(s.context.Stub(), key, value)
}

func (s *stateOp) List(objectType string, target interface{}) (result []interface{}, err error) {
	return state.List(s.context.Stub(), objectType, target)
}
