// Package convert for transforming  between json serialized  []byte and go structs
package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
)

var (
	// ErrUnableToConvertNilToStruct - nil cannot be converted to struct
	ErrUnableToConvertNilToStruct = errors.New(`unable to convert nil to [struct,array,slice,ptr]`)
	// ErrUnableToConvertValueToStruct - value  cannot be converted to struct
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

	// create copy

	switch target.(type) {
	case string:
		return string(bb), nil
	case []byte:
		return bb, nil
	case int:
		return strconv.Atoi(string(bb))

	case []string:
		arrInterface, err := UnmarshallPtr(bb, &target)

		if err != nil {
			return nil, err
		}
		arrString := []string{}
		for _, v := range arrInterface.([]interface{}) {
			arrString = append(arrString, v.(string))
		}
		return arrString, nil
	case FromByter:
		return target.(FromByter).FromBytes(bb)

	default:
		return FromBytesToStruct(bb, target)
	}

}

// FromResponse converts response.Payload to target
func FromResponse(response peer.Response, target interface{}) (result interface{}, err error) {
	if response.Status == shim.ERROR {
		return nil, errors.New(response.Message)
	}
	return FromBytes(response.Payload, target)
}

// FromBytesToStruct converts []byte to struct,array,slice depending on target type
func FromBytesToStruct(bb []byte, target interface{}) (result interface{}, err error) {
	if bb == nil {
		return nil, ErrUnableToConvertNilToStruct
	}
	targetType := reflect.TypeOf(target).Kind()

	switch targetType {
	case reflect.Struct:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		// will be map[string]interface{}
		return UnmarshallPtr(bb, &target)
	case reflect.Ptr:
		return UnmarshallPtr(bb, target)

	default:
		return nil, fmt.Errorf(
			`toBytes converting supports ToByter interface,struct,array,slice and string, current type is %s`,
			targetType)
	}
}

// UnmarshallPtr unmarshalls  [] byte to pointer, and returns value pointed to
func UnmarshallPtr(bb []byte, to interface{}) (result interface{}, err error) {

	targetPtr := reflect.New(reflect.ValueOf(to).Elem().Type()).Interface()

	err = json.Unmarshal(bb, targetPtr)
	if err != nil {
		return nil, errors.Wrap(err, ErrUnableToConvertValueToStruct.Error())
	}
	return reflect.Indirect(reflect.ValueOf(targetPtr)).Interface(), nil
}

// ToBytes converts inteface{} (string, []byte , struct to ToByter interface to []byte for storing in state
func ToBytes(value interface{}) ([]byte, error) {
	switch value.(type) {

	// first priority if value implements ToByter interface
	case ToByter:
		return value.(ToByter).ToBytes()
	case string:
		return []byte(value.(string)), nil
	case uint:
		return []byte(fmt.Sprint(value.(uint))), nil
	case int:
		return []byte(fmt.Sprint(value.(int))), nil
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

			// used when type based on string
		case reflect.String:
			return []byte(reflect.ValueOf(value).String()), nil

		default:
			return nil, fmt.Errorf(
				`toBytes converting supports ToByter interface,struct,array,slice and string, current type is %s`,
				valueType)
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

func TimestampToTime(ts *timestamp.Timestamp) time.Time {
	return time.Unix(ts.GetSeconds(), int64(ts.GetNanos()))
}
