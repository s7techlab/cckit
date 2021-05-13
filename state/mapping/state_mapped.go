package mapping

import (
	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/state"
)

type (
	ProtoStateMapped struct {
		instanceKey interface{}
		stateMapper StateMapper
	}
)

func NewProtoStateMapped(instanceKey interface{}, stateMapper StateMapper) *ProtoStateMapped {
	return &ProtoStateMapped{
		instanceKey: instanceKey,
		stateMapper: stateMapper,
	}
}

func (pm *ProtoStateMapped) Key() (state.Key, error) {
	switch instanceKey := pm.instanceKey.(type) {
	case []string:
		return instanceKey, nil
	default:
		return pm.stateMapper.PrimaryKey(instanceKey)
	}
}

func (pm *ProtoStateMapped) Keys() ([]state.KeyValue, error) {
	return pm.stateMapper.Keys(pm.instanceKey)
}

func (pm *ProtoStateMapped) ToBytes() ([]byte, error) {
	return proto.Marshal(pm.instanceKey.(proto.Message))
}

func (pm *ProtoStateMapped) Mapper() StateMapper {
	return pm.stateMapper
}
