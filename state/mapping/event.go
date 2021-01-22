package mapping

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/s7techlab/cckit/state"
)

type (
	EventImpl struct {
		event    state.Event
		mappings EventMappings
	}
)

func NewEvent(stub shim.ChaincodeStubInterface, mappings EventMappings) *EventImpl {
	return &EventImpl{
		event:    state.NewEvent(stub),
		mappings: mappings,
	}
}

func WrapEvent(event state.Event, mappings EventMappings) *EventImpl {
	return &EventImpl{
		event:    event,
		mappings: mappings,
	}
}

func (e *EventImpl) mapIfMappingExists(entry interface{}) (mapped interface{}, exists bool, err error) {
	if !e.mappings.Exists(entry) {
		return entry, false, nil
	}
	mapped, err = e.mappings.Map(entry)
	return mapped, true, err
}

func (e *EventImpl) Set(entry interface{}, value ...interface{}) error {
	mapped, exists, err := e.mapIfMappingExists(entry)
	if err != nil {
		return err
	}

	if !exists && len(value) == 0 {
		return fmt.Errorf(`%s: %s`, ErrEventMappingNotFound, mapKey(entry))
	}
	return e.event.Set(mapped, value...)
}

func (e *EventImpl) UseNameTransformer(nt state.StringTransformer) state.Event {
	return e.event.UseNameTransformer(nt)
}

func (e *EventImpl) UseSetTransformer(tb state.ToBytesTransformer) state.Event {
	return e.event.UseSetTransformer(tb)
}
