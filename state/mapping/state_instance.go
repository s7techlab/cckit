package mapping

import (
	"github.com/s7techlab/cckit/state"
)

type (
	StateInstance struct {
		// instance can be instance itself or key for instance
		// key can be proto or Key ( []string )
		instance    interface{}
		stateMapper StateMapper
		serializer  state.Serializer
	}
)

func NewStateInstance(instance interface{}, stateMapper StateMapper, serializer state.Serializer) *StateInstance {
	return &StateInstance{
		instance:    instance,
		stateMapper: stateMapper,
		serializer:  serializer,
	}
}

func (si *StateInstance) Key() (state.Key, error) {
	switch instance := si.instance.(type) {
	case []string:
		return instance, nil
	default:
		return si.stateMapper.PrimaryKey(instance)
	}
}

func (si *StateInstance) Keys() ([]state.KeyValue, error) {
	return si.stateMapper.Keys(si.instance)
}

func (si *StateInstance) ToBytes() ([]byte, error) {
	return si.serializer.ToBytes(si.instance)
}

func (si *StateInstance) Mapper() StateMapper {
	return si.stateMapper
}
