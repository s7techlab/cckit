package convert

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	// ErrUnsupportedType - type cannot be translate to or from []byte
	ErrUnsupportedType = errors.New(`from []byte converting supports FromByter interface and string`)
)

// FromByter interface supports FromBytes func for converting to structure
type FromByter interface {
	FromBytes([]byte) (interface{}, error)
}

// ToByter interface supports ToBytes func, marshalling to []byte (json.Marshall)
type ToByter interface {
	ToBytes() []byte
}

// FromBytes converts []byte to target interface
func FromBytes(bb []byte, target interface{}) (result interface{}, err error) {
	switch target.(type) {
	case string:
		return string(bb), nil
	case FromByter:
		return target.(FromByter).FromBytes(bb)
	}

	return nil, ErrUnsupportedType
}

// ToBytes converts inteface{} (string, []byte , struct to ToByter interface to []byte for storing in state
func ToBytes(value interface{}) ([]byte, error) {

	switch value.(type) {
	case string:
		return []byte(value.(string)), nil
	case []byte:
		return value.([]byte), nil
	case ToByter:
		return value.(ToByter).ToBytes(), nil

		//case reflect.Struct:
		//	//TODO:
		//	marshalled, _ := json.Marshal(reflect.ValueOf(arg))
		//	args[i] = marshalled

	default:
		return nil, errors.New(`to []byte converting supports ToByter interface and string`)
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
