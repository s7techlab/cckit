package mapping

import (
	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/state"
)

type (
	ProtoStateMapped struct {
		// instance can be instance itself or key for instance
		// key can be proto or Key ( []string )
		instance    interface{}
		stateMapper StateMapper
	}
)

func NewProtoStateMapped(instance interface{}, stateMapper StateMapper) *ProtoStateMapped {
	return &ProtoStateMapped{
		instance:    instance,
		stateMapper: stateMapper,
	}
}

func (pm *ProtoStateMapped) Key() (state.Key, error) {
	switch instance := pm.instance.(type) {
	case []string:
		return instance, nil
	default:
		return pm.stateMapper.PrimaryKey(instance)
	}
}

func (pm *ProtoStateMapped) Keys() ([]state.KeyValue, error) {
	return pm.stateMapper.Keys(pm.instance)
}

func (pm *ProtoStateMapped) ToBytes() ([]byte, error) {
	return proto.Marshal(pm.instance.(proto.Message))
}

func (pm *ProtoStateMapped) Mapper() StateMapper {
	return pm.stateMapper
}
