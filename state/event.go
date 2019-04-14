package state

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/convert"
)

type (
	// Event interface for working with events in chaincode
	Event interface {
		Set(entry interface{}, value ...interface{}) error
		UseSetTransformer(ToBytesTransformer) Event
		UseNameTransformer(StringTransformer) Event
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
		stub            shim.ChaincodeStubInterface
		NameTransformer StringTransformer
		SetTransformer  ToBytesTransformer
	}
)

// NewEvent creates wrapper on shim.ChaincodeStubInterface for working with events
func NewEvent(stub shim.ChaincodeStubInterface) *EventImpl {
	return &EventImpl{
		stub:            stub,
		NameTransformer: NameAsIs,
		SetTransformer:  ConvertToBytes,
	}
}

func (e *EventImpl) UseSetTransformer(tb ToBytesTransformer) Event {
	e.SetTransformer = tb
	return e
}

func (e *EventImpl) UseNameTransformer(nt StringTransformer) Event {
	e.NameTransformer = nt
	return e
}

func (e *EventImpl) Set(entry interface{}, values ...interface{}) error {
	name, value, err := e.ArgNameValue(entry, values)
	if err != nil {
		return err
	}

	nameStr, err := e.NameTransformer(name)
	if err != nil {
		return err
	}

	bb, err := e.SetTransformer(value)
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
	switch n := name.(type) {
	case Namer:
		return n.Name()
	case string:
		return n, nil
	}

	return ``, ErrUnableToCreateEventName
}
