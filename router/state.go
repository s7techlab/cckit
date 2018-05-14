package router

import (
	"github.com/s7techlab/cckit/state"
)

// State interface for chain code CRUD operations
type State interface {
	Get(key interface{}, target ...interface{}) (result interface{}, err error)
	Exists(key interface{}) (exists bool, err error)
	Put(key interface{}, value interface{}) (err error)
	Insert(key interface{}, value interface{}) (err error)
	List(objectType interface{}, target interface{}) (result []interface{}, err error)
}

type ContextState struct {
	context Context
}

func (s ContextState) Get(key interface{}, target ...interface{}) (result interface{}, err error) {
	return state.Get(s.context.Stub(), key, target...)
}

func (s ContextState) Exists(key interface{}) (exists bool, err error) {
	return state.Exists(s.context.Stub(), key)
}

func (s ContextState) Put(key interface{}, value interface{}) (err error) {
	return state.Put(s.context.Stub(), key, value)
}

func (s ContextState) Insert(key interface{}, value interface{}) (err error) {
	return state.Insert(s.context.Stub(), key, value)
}

func (s ContextState) List(objectType interface{}, target interface{}) (result []interface{}, err error) {
	return state.List(s.context.Stub(), objectType, target)
}
