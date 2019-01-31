package state

import "github.com/s7techlab/cckit/convert"

type (
	FromBytesTransformer func(bb []byte, config ...interface{}) (interface{}, error)
	ToBytesTransformer   func(v interface{}, config ...interface{}) ([]byte, error)
	KeyTransformer       func(key interface{}) ([]string, error)
)

func ConvertFromBytes(bb []byte, config ...interface{}) (interface{}, error) {
	// convertation not needed
	if len(config) == 0 {
		return bb, nil
	}
	return convert.FromBytes(bb, config[0])
}

func ConvertToBytes(v interface{}, config ...interface{}) ([]byte, error) {
	return convert.ToBytes(v)
}

//  ConvertKey returns string parts of composite key
func ConvertKey(key interface{}) ([]string, error) {
	switch key.(type) {
	case Keyer:
		return key.(Keyer).Key()
	case string:
		return []string{key.(string)}, nil
	case []string:
		return key.([]string), nil
	}

	return nil, ErrUnableToCreateKey
}
