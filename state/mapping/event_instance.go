package mapping

import (
	"github.com/s7techlab/cckit/state"
)

type (
	EventInstance struct {
		instance    interface{}
		eventMapper EventMapper
		serializer  state.Serializer
	}
)

func NewEventInstance(instance interface{}, eventMapper EventMapper, serializer state.Serializer) (*EventInstance, error) {
	return &EventInstance{
		instance:    instance,
		eventMapper: eventMapper,
		serializer:  serializer,
	}, nil
}

func (ei EventInstance) Name() (string, error) {
	return ei.eventMapper.Name(ei.instance)
}

func (ei EventInstance) ToBytes() ([]byte, error) {
	return ei.serializer.ToBytes(ei.instance)
}
