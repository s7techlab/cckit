package mapping

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/state"
)

type (
	StateImpl struct {
		state    state.State
		mappings Mappings
	}
)

func NewState(stub shim.ChaincodeStubInterface, mappings Mappings) *StateImpl {
	return &StateImpl{
		state:    state.New(stub),
		mappings: mappings,
	}
}

func (s *StateImpl) Get(key interface{}, target ...interface{}) (result interface{}, err error) {
	return s.state.Get(key, target...)
}

func (s *StateImpl) GetInt(key interface{}, defaultValue int) (result int, err error) {
	return s.state.GetInt(key, defaultValue)
}

func (s *StateImpl) GetHistory(key interface{}, target interface{}) (result state.HistoryEntryList, err error) {
	return s.state.GetHistory(key, target)
}

func (s *StateImpl) Exists(entry interface{}) (exists bool, err error) {
	if s.mappings.Exists(entry) {
		entry, err = s.mappings.Map(entry)
		if err != nil {
			return false, errors.Wrap(err, `mapping`)
		}
	}
	return s.state.Exists(entry)
}

func (s *StateImpl) Put(entry interface{}, value ...interface{}) (err error) {
	if s.mappings.Exists(entry) {
		entry, err = s.mappings.Map(entry)
		if err != nil {
			return errors.Wrap(err, `mapping`)
		}
	}
	return s.state.Put(entry, value...)
}

func (s *StateImpl) Insert(entry interface{}, value ...interface{}) (err error) {
	if s.mappings.Exists(entry) {
		entry, err = s.mappings.Map(entry)
		if err != nil {
			return errors.Wrap(err, `mapping`)
		}
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
	if s.mappings.Exists(entry) {
		entry, err = s.mappings.Map(entry)
		if err != nil {
			return errors.Wrap(err, `mapping`)
		}
	}
	return s.state.Delete(entry)
}
