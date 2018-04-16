package convert

import (
	"github.com/pkg/errors"
	"fmt"
)

type FromByter interface {
	FromBytes([]byte) (interface{}, error)
}

type ToByter interface {
	ToBytes() ([]byte)
}

func FromBytes(bb []byte, target interface{}) (result interface{}, err error) {

	switch target.(type) {
	case string:
		return string(bb), nil

	case FromByter:
		return target.(FromByter).FromBytes(bb)

	default:
		return nil, errors.New(`from []byte converting supports FromByter interface and string`)
	}

	return
}

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
