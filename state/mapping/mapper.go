package mapping

import (
	"github.com/golang/protobuf/proto"
)

type (
	ProtoMapper struct {
		instance interface{}
		mapping  Mapping
	}
)

func NewProtoMapper(instance interface{}, mapping Mapping) (*ProtoMapper, error) {
	return &ProtoMapper{instance, mapping}, nil
}

func (pm *ProtoMapper) Key() ([]string, error) {
	return pm.mapping.PrimaryKey(pm.instance)
}

func (mp *ProtoMapper) ToBytes() ([]byte, error) {
	return proto.Marshal(mp.instance.(proto.Message))
}
