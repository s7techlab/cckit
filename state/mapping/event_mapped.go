package mapping

import "github.com/golang/protobuf/proto"

type (
	ProtoEventMapped struct {
		instance    interface{}
		eventMapper EventMapper
	}
)

func NewProtoEventMapped(instance interface{}, eventMapper EventMapper) (*ProtoEventMapped, error) {
	return &ProtoEventMapped{instance, eventMapper}, nil
}

func (em *ProtoEventMapped) Name() (string, error) {
	return em.eventMapper.Name(em.instance)
}

func (em *ProtoEventMapped) ToBytes() ([]byte, error) {
	return proto.Marshal(em.instance.(proto.Message))
}
