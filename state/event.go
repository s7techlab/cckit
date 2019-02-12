package state

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/convert"
)

type (
	// Event interface for working with events in chaincode
	Event interface {
		Set(entry interface{}, value ...interface{}) error
	}

	Namer interface {
		Name() (string, error)
	}
	// NameValue interface combines Name() as ToByter methods - event representation
	NameValue interface {
		Namer
		convert.ToByter
	}

	EventImpl struct {
		stub                shim.ChaincodeStubInterface
		NameTransformer     NameTransformer
		EventSetTransformer ToBytesTransformer
	}
)

// NewEvent creates wrapper on shim.ChaincodeStubInterface for working with events
func NewEvent(stub shim.ChaincodeStubInterface) *EventImpl {
	return &EventImpl{
		stub:                stub,
		NameTransformer:     NameAsIs,
		EventSetTransformer: ConvertToBytes,
	}
}

func (e *EventImpl) Set(entry interface{}, values ...interface{}) error {
	name, value, err := e.ArgNameValue(entry, values)

	nameStr, err := e.NameTransformer(name)
	if err != nil {
		return err
	}

	bb, err := e.EventSetTransformer(value)
	if err != nil {
		return err
	}

	return e.stub.SetEvent(nameStr, bb)
}

func (e *EventImpl) ArgNameValue(arg interface{}, values []interface{}) (name string, value interface{}, err error) {
	// name must be
	name, err = NormalizeEventName(arg)
	if err != nil {
		return
	}

	switch len(values) {
	// arg is name and  value
	case 0:
		return name, arg, nil
	case 1:
		return name, values[0], nil
	default:
		return ``, nil, ErrAllowOnlyOneValue
	}
}

func NormalizeEventName(name interface{}) (string, error) {
	switch name.(type) {
	case Namer:
		return name.(Namer).Name()
	case string:
		return name.(string), nil
	}

	return ``, ErrUnableToCreateEventName
}
