package mapping

import (
	"fmt"

	"github.com/s7techlab/cckit/state/schema"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/state"
)

type (
	MappedState interface {
		state.State

		// ListWith allows to refine search criteria by adding to namespace key parts
		ListWith(schema interface{}, key state.Key) (result interface{}, err error)

		// GetByUniqKey return one entry
		GetByUniqKey(schema interface{}, idx string, idxVal []string, target ...interface{}) (result interface{}, err error)

		// GetByUniqKey return list of entries
		//GetByKey(schema interface{}, key string, keyValue []interface{}) (result interface{}, err error)
	}

	StateImpl struct {
		state    state.State
		mappings StateMappings
	}
)

func WrapState(s state.State, mappings StateMappings) *StateImpl {
	return &StateImpl{
		state:    s,
		mappings: mappings,
	}
}

func (s *StateImpl) MappingNamespace(schema interface{}) (state.Key, error) {
	m, err := s.mappings.Get(schema)
	if err != nil {
		return nil, err
	}

	return m.Namespace(), nil
}

func (s *StateImpl) Get(entry interface{}, target ...interface{}) (result interface{}, err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.Get(entry, target...) // return as is
	}

	return s.state.Get(mapped, target...)
}

func (s *StateImpl) GetInt(key interface{}, defaultValue int) (result int, err error) {
	return s.state.GetInt(key, defaultValue)
}

func (s *StateImpl) GetHistory(key interface{}, target interface{}) (result state.HistoryEntryList, err error) {
	return s.state.GetHistory(key, target)
}

func (s *StateImpl) Exists(entry interface{}) (exists bool, err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.Exists(entry) // return as is
	}

	return s.state.Exists(mapped)
}

func (s *StateImpl) Put(entry interface{}, value ...interface{}) (err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.Put(entry, value...) // return as is
	}

	keyRefs, err := mapped.Keys() // additional keys
	if err != nil {
		return
	}

	// delete previous key refs if key exists

	// put uniq key refs. if key already exists - error returned
	for _, kr := range keyRefs {
		if err = s.state.Put(kr); err != nil {
			return fmt.Errorf(`%s: %s`, ErrMappingUniqKeyExists, err)
		}
	}

	return s.state.Put(mapped)
}

func (s *StateImpl) Insert(entry interface{}, value ...interface{}) (err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.Insert(entry, value...) // return as is
	}

	keyRefs, err := mapped.Keys() // additional keys
	if err != nil {
		return
	}

	// insert uniq key refs. if key already exists - error returned
	for _, kr := range keyRefs {
		if err = s.state.Insert(kr); err != nil {
			return fmt.Errorf(`%s: %s`, ErrMappingUniqKeyExists, err)
		}
	}

	return s.state.Insert(mapped)
}

func (s *StateImpl) List(namespace interface{}, target ...interface{}) (result interface{}, err error) {
	if s.mappings.Exists(namespace) {
		m, err := s.mappings.Get(namespace)
		if err != nil {
			return nil, errors.Wrap(err, `mapping`)
		}

		namespace = m.Namespace()
		s.Logger().Debugf(`state mapped LIST with namespace: %s`, namespace)
		target = targetFromMapping(m)
	}

	return s.state.List(namespace, target...)
}

func targetFromMapping(m StateMapper) (target []interface{}) {
	target = []interface{}{m.Schema()}
	if list := m.List(); list != nil {
		target = append(target, list)
	}
	return
}

func (s *StateImpl) ListWith(entry interface{}, key state.Key) (result interface{}, err error) {
	if !s.mappings.Exists(entry) {
		return nil, ErrStateMappingNotFound
	}
	m, err := s.mappings.Get(entry)
	if err != nil {
		return nil, errors.Wrap(err, `mapping`)
	}

	namespace := m.Namespace()
	s.Logger().Debugf(`state mapped LIST with namespace: %s`, namespace, namespace.Append(key))

	return s.state.List(namespace.Append(key), targetFromMapping(m)...)
}

func (s *StateImpl) GetByUniqKey(entry interface{}, idx string, idxVal []string, target ...interface{}) (result interface{}, err error) {
	if !s.mappings.Exists(entry) {
		return nil, ErrStateMappingNotFound
	}

	keyRef, err := s.state.Get(NewKeyRefIdMapped(entry, idx, idxVal), &schema.KeyRef{})
	if err != nil {
		return nil, err
	}

	return s.state.Get(keyRef.(*schema.KeyRef).PKey, target...)
}

func (s *StateImpl) Delete(entry interface{}) (err error) {

	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.Delete(entry) // return as is
	}

	return s.state.Delete(mapped)
}

func (s *StateImpl) Logger() *shim.ChaincodeLogger {
	return s.state.Logger()
}

func (s *StateImpl) UseKeyTransformer(kt state.KeyTransformer) state.State {
	return s.state.UseKeyTransformer(kt)
}

func (s *StateImpl) UseStateGetTransformer(fb state.FromBytesTransformer) state.State {
	return s.state.UseStateGetTransformer(fb)
}

func (s *StateImpl) UseStatePutTransformer(tb state.ToBytesTransformer) state.State {
	return s.state.UseStatePutTransformer(tb)
}
