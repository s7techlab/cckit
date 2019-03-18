package mapping

import (
	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/state"
)

type (
	ProtoStateMapped struct {
		instance    interface{}
		stateMapper StateMapper
	}
)

func NewProtoStateMapped(instance interface{}, stateMapper StateMapper) *ProtoStateMapped {
	return &ProtoStateMapped{instance, stateMapper}
}

func (pm *ProtoStateMapped) Key() (state.Key, error) {
	return pm.stateMapper.PrimaryKey(pm.instance)
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
