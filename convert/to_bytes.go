package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"

	"github.com/golang/protobuf/proto"
)

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

// ToBytes converts inteface{} (string, []byte , struct to ToByter interface to []byte for storing in state
func ToBytes(value interface{}) ([]byte, error) {
	if value == nil {
		return nil, nil
	}

	switch value.(type) {

	// first priority if value implements ToByter interface
	case ToByter:
		return value.(ToByter).ToBytes()
	case proto.Message:
		return proto.Marshal(proto.Clone(value.(proto.Message)))
	case bool:
		return []byte(strconv.FormatBool(value.(bool))), nil
	case string:
		return []byte(value.(string)), nil
	case uint:
		return []byte(fmt.Sprint(value.(uint))), nil
	case int:
		return []byte(fmt.Sprint(value.(int))), nil
	case int32:
		return []byte(fmt.Sprint(value.(int32))), nil
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
