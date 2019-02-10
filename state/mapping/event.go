package mapping

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
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

func (e *EventImpl) mapIfMappingExists(entry interface{}) (mapped interface{}, err error) {
	if !e.mappings.Exists(entry) {
		return entry, nil
	}
	return e.mappings.Map(entry)
}

func (e *EventImpl) Set(entry interface{}, value ...interface{}) (err error) {
	if entry, err = e.mapIfMappingExists(entry); err != nil {
		return err
	}
	return e.event.Set(entry, value...)
}
