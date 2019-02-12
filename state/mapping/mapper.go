package mapping

import (
	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/state"
)

type (
	ProtoStateMapper struct {
		instance    interface{}
		stateMapper StateMapper
	}

	ProtoEventMapper struct {
		instance    interface{}
		eventMapper EventMapper
	}
)

func NewProtoStateMapper(instance interface{}, stateMapper StateMapper) (*ProtoStateMapper, error) {
	return &ProtoStateMapper{instance, stateMapper}, nil
}

func (pm *ProtoStateMapper) Key() (state.Key, error) {
	return pm.stateMapper.PrimaryKey(pm.instance)
}

func (pm *ProtoStateMapper) ToBytes() ([]byte, error) {
	return proto.Marshal(pm.instance.(proto.Message))
}

func NewProtoEventMapper(instance interface{}, eventMapper EventMapper) (*ProtoEventMapper, error) {
	return &ProtoEventMapper{instance, eventMapper}, nil
}

func (em *ProtoEventMapper) Name() (string, error) {
	return em.eventMapper.Name(em.instance)
}

func (em *ProtoEventMapper) ToBytes() ([]byte, error) {
	return proto.Marshal(em.instance.(proto.Message))
}
