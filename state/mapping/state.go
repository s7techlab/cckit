package mapping

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/state"
	"github.com/s7techlab/cckit/state/schema"
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

	Impl struct {
		state    state.State
		mappings StateMappings
	}
)

func WrapState(s state.State, mappings StateMappings) *Impl {
	return &Impl{
		state:    s,
		mappings: mappings,
	}
}

func (s *Impl) MappingNamespace(schema interface{}) (state.Key, error) {
	m, err := s.mappings.Get(schema)
	if err != nil {
		return nil, err
	}

	return m.Namespace(), nil
}

func (s *Impl) Get(entry interface{}, target ...interface{}) (interface{}, error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.Get(entry, target...) // return as is
	}

	// target was not set, but we can knew about target from mapping
	if len(target) == 0 {
		var targetFromMapping interface{}
		if mapped.Mapper().KeyerFor() != nil {
			targetFromMapping = mapped.Mapper().KeyerFor()
		} else {
			targetFromMapping = mapped.Mapper().Schema()
		}
		target = append(target, targetFromMapping)
	}

	return s.state.Get(mapped, target...)
}

func (s *Impl) GetInt(entry interface{}, defaultValue int) (int, error) {
	return s.state.GetInt(entry, defaultValue)
}

func (s *Impl) GetHistory(entry interface{}, target interface{}) (state.HistoryEntryList, error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.GetHistory(entry, target) // return as is
	}

	return s.state.GetHistory(mapped, target)
}

func (s *Impl) Exists(entry interface{}) (bool, error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.Exists(entry) // return as is
	}

	return s.state.Exists(mapped)
}

func (s *Impl) Put(entry interface{}, value ...interface{}) error {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.Put(entry, value...) // return as is
	}

	keyRefs, err := mapped.Keys() // additional keys
	if err != nil {
		return err
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

func (s *Impl) Insert(entry interface{}, value ...interface{}) error {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.Insert(entry, value...) // return as is
	}

	keyRefs, err := mapped.Keys() // additional keys
	if err != nil {
		return err
	}

	// insert uniq key refs. if key already exists - error returned
	for _, kr := range keyRefs {
		if err = s.state.Insert(kr); err != nil {
			return fmt.Errorf(`%s: %s`, ErrMappingUniqKeyExists, err)
		}
	}

	return s.state.Insert(mapped)
}

func (s *Impl) List(entry interface{}, target ...interface{}) (interface{}, error) {
	if !s.mappings.Exists(entry) {
		return s.state.List(entry, target...)
	}

	m, err := s.mappings.Get(entry)
	if err != nil {
		return nil, errors.Wrap(err, `mapping`)
	}

	namespace := m.Namespace()
	s.Logger().Debugf(`state mapped LIST with namespace: %s`, namespace)

	return s.state.List(namespace, m.Schema(), m.List())
}

func (s *Impl) ListWith(entry interface{}, key state.Key) (result interface{}, err error) {
	if !s.mappings.Exists(entry) {
		return nil, ErrStateMappingNotFound
	}
	m, err := s.mappings.Get(entry)
	if err != nil {
		return nil, errors.Wrap(err, `mapping`)
	}

	namespace := m.Namespace()
	s.Logger().Debugf(`state mapped LIST with namespace: %s`, namespace, namespace.Append(key))

	return s.state.List(namespace.Append(key), m.Schema(), m.List())
}

func (s *Impl) GetByUniqKey(
	entry interface{}, idx string, idxVal []string, target ...interface{}) (result interface{}, err error) {

	if !s.mappings.Exists(entry) {
		return nil, ErrStateMappingNotFound
	}

	keyRef, err := s.state.Get(NewKeyRefIDMapped(entry, idx, idxVal), &schema.KeyRef{})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(`uniq index: {%s}.%s`, mapKey(entry), idx))
	}

	return s.state.Get(keyRef.(*schema.KeyRef).PKey, target...)
}

func (s *Impl) Delete(entry interface{}) error {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.Delete(entry) // return as is
	}

	// Entry can be record to delete or reference to record
	// If entry is keyer entity for another entry (reference)
	if mapped.Mapper().KeyerFor() != nil {
		referenceEntry, err := s.Get(entry)
		if err != nil {
			return err
		}

		mapped, err = s.mappings.Map(referenceEntry)
		if err != nil {
			return err
		}
	}

	keyRefs, err := mapped.Keys() // additional keys
	if err != nil {
		return err
	}

	// delete uniq key refs
	for _, kr := range keyRefs {
		if err = s.state.Delete(kr); err != nil {
			return errors.Wrap(err, `delete ref key`)
		}
	}

	return s.state.Delete(mapped)
}

func (s *Impl) Logger() *shim.ChaincodeLogger {
	return s.state.Logger()
}

func (s *Impl) UseKeyTransformer(kt state.KeyTransformer) state.State {
	return s.state.UseKeyTransformer(kt)
}

func (s *Impl) UseStateGetTransformer(fb state.FromBytesTransformer) state.State {
	return s.state.UseStateGetTransformer(fb)
}

func (s *Impl) UseStatePutTransformer(tb state.ToBytesTransformer) state.State {
	return s.state.UseStatePutTransformer(tb)
}

func (s *Impl) GetPrivate(collection string, entry interface{}, target ...interface{}) (result interface{}, err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.GetPrivate(collection, entry, target...) // return as is
	}

	return s.state.GetPrivate(collection, mapped, target...)
}

func (s *Impl) DeletePrivate(collection string, entry interface{}) (err error) {

	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.DeletePrivate(collection, entry) // return as is
	}

	return s.state.DeletePrivate(collection, mapped)
}

func (s *Impl) ListPrivate(collection string, usePrivateDataIterator bool, namespace interface{}, target ...interface{}) (result interface{}, err error) {
	if !s.mappings.Exists(namespace) {
		return s.state.ListPrivate(collection, usePrivateDataIterator, namespace, target...)
	}
	m, err := s.mappings.Get(namespace)
	if err != nil {
		return nil, errors.Wrap(err, `mapping`)
	}

	namespace = m.Namespace()
	s.Logger().Debugf(`private state mapped LIST with namespace: %s`, namespace)
	return s.state.ListPrivate(collection, usePrivateDataIterator, namespace, target[0], m.List())
}

func (s *Impl) InsertPrivate(collection string, entry interface{}, value ...interface{}) (err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.InsertPrivate(collection, entry, value...) // return as is
	}

	keyRefs, err := mapped.Keys() // additional keys
	if err != nil {
		return
	}

	// insert uniq key refs. if key already exists - error returned
	for _, kr := range keyRefs {
		if err = s.state.InsertPrivate(collection, kr); err != nil {
			return fmt.Errorf(`%s: %s`, ErrMappingUniqKeyExists, err)
		}
	}

	return s.state.InsertPrivate(collection, mapped)
}

func (s *Impl) PutPrivate(collection string, entry interface{}, value ...interface{}) (err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.PutPrivate(collection, entry, value...) // return as is
	}

	keyRefs, err := mapped.Keys() // additional keys
	if err != nil {
		return
	}

	// delete previous key refs if key exists

	// put uniq key refs. if key already exists - error returned
	for _, kr := range keyRefs {
		if err = s.state.PutPrivate(collection, kr); err != nil {
			return fmt.Errorf(`%s: %s`, ErrMappingUniqKeyExists, err)
		}
	}

	return s.state.PutPrivate(collection, mapped)
}

func (s *Impl) ExistsPrivate(collection string, entry interface{}) (exists bool, err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.state.ExistsPrivate(collection, entry) // return as is
	}

	return s.state.ExistsPrivate(collection, mapped)
}
