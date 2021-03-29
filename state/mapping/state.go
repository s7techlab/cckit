package mapping

import (
	"fmt"

	"go.uber.org/zap"

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
		// Deprecated: use GetByKey
		GetByUniqKey(schema interface{}, idx string, idxVal []string, target ...interface{}) (result interface{}, err error)

		// GetByKey
		GetByKey(schema interface{}, idx string, idxVal []string, target ...interface{}) (result interface{}, err error)
	}

	Impl struct {
		state.State
		mappings StateMappings
	}
)

func WrapState(s state.State, mappings StateMappings) *Impl {
	return &Impl{
		State:    s,
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
		return s.State.Get(entry, target...) // return as is
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

	return s.State.Get(mapped, target...)
}

func (s *Impl) GetHistory(entry interface{}, target interface{}) (state.HistoryEntryList, error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.State.GetHistory(entry, target) // return as is
	}

	return s.State.GetHistory(mapped, target)
}

func (s *Impl) Exists(entry interface{}) (bool, error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.State.Exists(entry) // return as is
	}

	return s.State.Exists(mapped)
}

func (s *Impl) Put(entry interface{}, value ...interface{}) error {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.State.Put(entry, value...) // return as is
	}

	// update ref keys
	if len(mapped.Mapper().Indexes()) > 0 {
		keyRefs, err := mapped.Keys() // key refs based on current entry value, defined by mapping indexes
		if err != nil {
			return errors.Wrap(err, `put mapping key refs`)
		}

		var insertKeyRefs, deleteKeyRefs []state.KeyValue
		//get previous entry value
		prevEntry, err := s.Get(entry)

		if err == nil { // prev exists

			// prev entry exists, calculate refs to delete and to insert
			prevMapped, err := s.mappings.Map(prevEntry)
			if err != nil {
				return errors.Wrap(err, `get prev`)
			}
			prevKeyRefs, err := prevMapped.Keys() // key refs based on current entry value, defined by mapping indexes
			if err != nil {
				return errors.Wrap(err, `previ keys`)
			}

			deleteKeyRefs, insertKeyRefs, err = KeyRefsDiff(prevKeyRefs, keyRefs)
			if err != nil {
				return errors.Wrap(err, `calculate ref keys diff`)
			}

		} else {
			// prev entry not exists, all current key refs should be inserted
			insertKeyRefs = keyRefs
		}

		// delete previous key refs if key exists
		for _, kr := range deleteKeyRefs {
			if err = s.State.Delete(kr); err != nil {
				return errors.Wrap(err, `delete previous mapping key ref`)
			}
		}

		// insert new key refs
		for _, kr := range insertKeyRefs {
			if err = s.State.Insert(kr); err != nil {
				return fmt.Errorf(`%s: %s`, ErrMappingUniqKeyExists, err)
			}
		}
	}

	return s.State.Put(mapped)
}

func (s *Impl) Insert(entry interface{}, value ...interface{}) error {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.State.Insert(entry, value...) // return as is
	}

	keyRefs, err := mapped.Keys() // key refs, defined by mapping indexes
	if err != nil {
		return err
	}

	// insert key refs, if key already exists - error returned
	for _, kr := range keyRefs {
		if err = s.State.Insert(kr); err != nil {
			return fmt.Errorf(`%s: %s`, ErrMappingUniqKeyExists, err)
		}
	}

	return s.State.Insert(mapped)
}

func (s *Impl) List(entry interface{}, target ...interface{}) (interface{}, error) {
	if !s.mappings.Exists(entry) {
		return s.State.List(entry, target...)
	}

	m, err := s.mappings.Get(entry)
	if err != nil {
		return nil, errors.Wrap(err, `mapping`)
	}

	namespace := m.Namespace()
	s.Logger().Debug(`state mapped LIST`, zap.String(`namespace`, namespace.String()))

	return s.State.List(namespace, m.Schema(), m.List())
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
	s.Logger().Debug(`state mapped LIST`, zap.String(`namespace`, namespace.String()), zap.String(`list`, namespace.Append(key).String()))

	return s.State.List(namespace.Append(key), m.Schema(), m.List())
}

func (s *Impl) GetByUniqKey(
	entry interface{}, idx string, idxVal []string, target ...interface{}) (result interface{}, err error) {
	return s.GetByKey(entry, idx, idxVal, target...)
}

func (s *Impl) GetByKey(
	entry interface{}, idx string, idxVal []string, target ...interface{}) (result interface{}, err error) {

	if !s.mappings.Exists(entry) {
		return nil, ErrStateMappingNotFound
	}

	keyRef, err := s.State.Get(NewKeyRefIDMapped(entry, idx, idxVal), &schema.KeyRef{})
	if err != nil {
		return nil, errors.Errorf(`%s: {%s}.%s: %s`, ErrIndexReferenceNotFound, mapKey(entry), idx, err)
	}

	return s.State.Get(keyRef.(*schema.KeyRef).PKey, target...)
}

func (s *Impl) Delete(entry interface{}) error {
	if !s.mappings.Exists(entry) {
		return s.State.Delete(entry) // return as is
	}

	// we need full entry data fro state
	// AND entry can be record to delete or reference to record
	// If entry is keyer entity for another entry (reference)
	entry, err := s.Get(entry)
	if err != nil {
		return err
	}

	mapped, err := s.mappings.Map(entry)
	if err != nil {
		return err
	}

	keyRefs, err := mapped.Keys() // additional keys
	if err != nil {
		return err
	}

	// delete uniq key refs
	for _, kr := range keyRefs {
		if err = s.State.Delete(kr); err != nil {
			return errors.Wrap(err, `delete ref key`)
		}
	}

	return s.State.Delete(mapped)
}

func (s *Impl) Logger() *zap.Logger {
	return s.State.Logger()
}

func (s *Impl) UseKeyTransformer(kt state.KeyTransformer) state.State {
	return s.State.UseKeyTransformer(kt)
}

func (s *Impl) GetPrivate(collection string, entry interface{}, target ...interface{}) (result interface{}, err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.State.GetPrivate(collection, entry, target...) // return as is
	}

	return s.State.GetPrivate(collection, mapped, target...)
}

func (s *Impl) DeletePrivate(collection string, entry interface{}) (err error) {

	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.State.DeletePrivate(collection, entry) // return as is
	}

	return s.State.DeletePrivate(collection, mapped)
}

func (s *Impl) ListPrivate(collection string, usePrivateDataIterator bool, namespace interface{}, target ...interface{}) (result interface{}, err error) {
	if !s.mappings.Exists(namespace) {
		return s.State.ListPrivate(collection, usePrivateDataIterator, namespace, target...)
	}
	m, err := s.mappings.Get(namespace)
	if err != nil {
		return nil, errors.Wrap(err, `mapping`)
	}

	namespace = m.Namespace()
	s.Logger().Debug(`private state mapped LIST`, zap.Reflect(`namespace`, namespace))
	return s.State.ListPrivate(collection, usePrivateDataIterator, namespace, target[0], m.List())
}

func (s *Impl) InsertPrivate(collection string, entry interface{}, value ...interface{}) (err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.State.InsertPrivate(collection, entry, value...) // return as is
	}

	keyRefs, err := mapped.Keys() // additional keys
	if err != nil {
		return
	}

	// insert uniq key refs. if key already exists - error returned
	for _, kr := range keyRefs {
		if err = s.State.InsertPrivate(collection, kr); err != nil {
			return fmt.Errorf(`%s: %s`, ErrMappingUniqKeyExists, err)
		}
	}

	return s.State.InsertPrivate(collection, mapped)
}

func (s *Impl) PutPrivate(collection string, entry interface{}, value ...interface{}) (err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.State.PutPrivate(collection, entry, value...) // return as is
	}

	keyRefs, err := mapped.Keys() // additional keys
	if err != nil {
		return
	}

	// delete previous key refs if key exists

	// put uniq key refs. if key already exists - error returned
	for _, kr := range keyRefs {
		if err = s.State.PutPrivate(collection, kr); err != nil {
			return fmt.Errorf(`%s: %s`, ErrMappingUniqKeyExists, err)
		}
	}

	return s.State.PutPrivate(collection, mapped)
}

func (s *Impl) ExistsPrivate(collection string, entry interface{}) (exists bool, err error) {
	mapped, err := s.mappings.Map(entry)
	if err != nil { // mapping is not exists
		return s.State.ExistsPrivate(collection, entry) // return as is
	}

	return s.State.ExistsPrivate(collection, mapped)
}
