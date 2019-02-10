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
		NameTransformer:     ConvertName,
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

func (s *EventImpl) ArgNameValue(arg interface{}, values []interface{}) (name interface{}, value interface{}, err error) {
	switch len(values) {

	// key is struct implementing keyer interface or has mapping instructions
	case 0:

		switch arg.(type) {

		case NameValue:
			name, err = arg.(NameValue).Name()
			if err != nil {
				return nil, nil, err
			}

			value, err := arg.(NameValue).ToBytes()
			if err != nil {
				return nil, nil, err
			}

			return name, value, nil

		default:
			return nil, nil, ErrEventEntryNotSupportNamerInterface
		}

	case 1:
		return arg, values[0], nil
	default:
		return nil, nil, ErrAllowOnlyOneValue
	}
}
