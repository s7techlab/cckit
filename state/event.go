package state

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/convert"
)

type (
	Event interface {
		Set(name string, payload interface{}) error
	}

	EventImpl struct {
		stub shim.ChaincodeStubInterface
	}
)

// NewEvent creates wrapper on shim.ChaincodeStubInterface for working with events
func NewEvent(stub shim.ChaincodeStubInterface) *EventImpl {
	return &EventImpl{
		stub: stub,
	}
}

func (e *EventImpl) Set(name string, payload interface{}) error {
	bb, err := convert.ToBytes(payload)
	if err != nil {
		return err
	}
	return e.stub.SetEvent(name, bb)
}
