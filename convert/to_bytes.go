package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
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

	switch v := value.(type) {

	// first priority if value implements ToByter interface
	case ToByter:
		return v.ToBytes()
	case proto.Message:
		return proto.Marshal(proto.Clone(v))
	case bool:
		return []byte(strconv.FormatBool(v)), nil
	case string:
		return []byte(v), nil
	case uint, int, int32:
		return []byte(fmt.Sprint(v)), nil
	case []byte:
		return v, nil

	default:
		valueType := reflect.TypeOf(value).Kind()

		switch valueType {
		case reflect.Ptr, reflect.Struct, reflect.Array, reflect.Map, reflect.Slice:
			return json.Marshal(value)
			// used when type based on string
		case reflect.String:
			return []byte(reflect.ValueOf(value).String()), nil

		default:
			return nil, fmt.Errorf(
				`toBytes converting supports ToByter interface,struct,array,slice,bool and string, current type is %s`,
				valueType)
		}

	}
}
