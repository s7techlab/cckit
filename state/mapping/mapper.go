package mapping

import (
	"github.com/golang/protobuf/proto"
)

type (
	ProtoStateMapper struct {
		instance    interface{}
		stateMapper StateMapper
	}
)

func NewProtoStateMapper(instance interface{}, stateMapper StateMapper) (*ProtoStateMapper, error) {
	return &ProtoStateMapper{instance, stateMapper}, nil
}

func (pm *ProtoStateMapper) Key() ([]string, error) {
	return pm.stateMapper.PrimaryKey(pm.instance)
}

func (pm *ProtoStateMapper) ToBytes() ([]byte, error) {
	return proto.Marshal(pm.instance.(proto.Message))
}
