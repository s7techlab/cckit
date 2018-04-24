package convert

import (
	"fmt"

	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

var (
	// ErrUnsupportedType - type cannot be translate to or from []byte
	ErrUnsupportedType              = errors.New(`fromBytes converting supports targets: FromByter interface, string,struct,array,slice,ptr`)
	ErrUnableToConvertNilToStruct   = errors.New(`unable to convert nil to [struct,array,slice,ptr]`)
	ErrUnableToConvertValueToStruct = errors.New(`unable to convert value to [struct,array,slice,ptr]`)
)

// FromByter interface supports FromBytes func for converting to structure
type FromByter interface {
	FromBytes([]byte) (interface{}, error)
}

// ToByter interface supports ToBytes func, marshalling to []byte (json.Marshall)
type ToByter interface {
	ToBytes() ([]byte, error)
}

// FromBytes converts []byte to target interface
func FromBytes(bb []byte, target interface{}) (result interface{}, err error) {
	switch target.(type) {
	case string:
		return string(bb), nil
	case FromByter:
		return target.(FromByter).FromBytes(bb)

	default:
		return FromBytesToStruct(bb, target)
	}

}

// FromBytesToStruct converts []byte to struct,array,slice depending on target type
func FromBytesToStruct(bb []byte, target interface{}) (result interface{}, err error) {

	if bb == nil {
		return nil, ErrUnableToConvertNilToStruct
	}

	targetType := reflect.TypeOf(target).Kind()
	var targetPtr interface{}

	switch targetType {

	case reflect.Struct:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		targetPtr = &target // will be map[string]interface{}
		fallthrough
	case reflect.Ptr:
		targetPtr = target // will be ptr to target struct, array ot slice
		err = json.Unmarshal(bb, targetPtr)

		if err != nil {
			err = errors.Wrap(err, ErrUnableToConvertValueToStruct.Error())
		}
		return targetPtr, err
	default:
		return nil, ErrUnsupportedType
	}
}

// ToBytes converts inteface{} (string, []byte , struct to ToByter interface to []byte for storing in state
func ToBytes(value interface{}) ([]byte, error) {
	switch value.(type) {

	// first priority if value implements ToByter interface
	case ToByter:
		return value.(ToByter).ToBytes()
	case string:
		return []byte(value.(string)), nil
	case []byte:
		return value.([]byte), nil

	default:
		valueType := reflect.TypeOf(value).Kind()

		switch valueType {

		case reflect.Ptr:
			fallthrough
		case reflect.Struct:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Slice:

			//TODO: problems when struct includes anonymous struct - anonymous struct is in separate field
			return json.Marshal(value)
		default:
			return nil, fmt.Errorf(`to []byte converting supports ToByter interface,struct,array,slice  and string, current type is %s`, valueType)
		}

	}

}

// ArgsToBytes converts func arguments to bytes
func ArgsToBytes(iargs ...interface{}) (aa [][]byte, err error) {
	args := make([][]byte, len(iargs))

	for i, arg := range iargs {
		val, err := ToBytes(arg)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf(`unable to convert invoke arg[%d]`, i))
		}
		args[i] = val
	}

	return args, nil
}
