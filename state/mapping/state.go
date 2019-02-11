package mapping

import (
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/state"
)

type (
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

func (s *StateImpl) mapIfMappingExists(entry interface{}) (mapped interface{}, err error) {
	if !s.mappings.Exists(entry) {
		return entry, nil
	}
	return s.mappings.Map(entry)
}

func (s *StateImpl) Get(entry interface{}, target ...interface{}) (result interface{}, err error) {
	if entry, err = s.mapIfMappingExists(entry); err != nil {
		return nil, err
	}
	return s.state.Get(entry, target...)
}

func (s *StateImpl) GetInt(key interface{}, defaultValue int) (result int, err error) {
	return s.state.GetInt(key, defaultValue)
}

func (s *StateImpl) GetHistory(key interface{}, target interface{}) (result state.HistoryEntryList, err error) {
	return s.state.GetHistory(key, target)
}

func (s *StateImpl) Exists(entry interface{}) (exists bool, err error) {
	if entry, err = s.mapIfMappingExists(entry); err != nil {
		return false, err
	}
	return s.state.Exists(entry)
}

func (s *StateImpl) Put(entry interface{}, value ...interface{}) (err error) {
	if entry, err = s.mapIfMappingExists(entry); err != nil {
		return err
	}
	return s.state.Put(entry, value...)
}

func (s *StateImpl) Insert(entry interface{}, value ...interface{}) (err error) {
	if entry, err = s.mapIfMappingExists(entry); err != nil {
		return err
	}
	return s.state.Insert(entry, value...)
}

func (s *StateImpl) List(namespace interface{}, target ...interface{}) (result []interface{}, err error) {
	if s.mappings.Exists(namespace) {
		m, err := s.mappings.Get(namespace)
		if err != nil {
			return nil, errors.Wrap(err, `mapping`)
		}

		namespace = m.Namespace()
		target = []interface{}{m.Schema()}
	}

	return s.state.List(namespace, target...)
}

func (s *StateImpl) Delete(entry interface{}) (err error) {
	if entry, err = s.mapIfMappingExists(entry); err != nil {
		return err
	}
	return s.state.Delete(entry)
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
