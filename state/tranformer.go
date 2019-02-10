package state

import "github.com/s7techlab/cckit/convert"

type (
	// ToBytesTransformer is used after getState operation for convert value
	FromBytesTransformer func(bb []byte, config ...interface{}) (interface{}, error)

	// ToBytesTransformer is used before putState operation for convert payload
	ToBytesTransformer func(v interface{}, config ...interface{}) ([]byte, error)

	// KeyTransformer is used before putState operation for convert key
	KeyTransformer func(key interface{}) ([]string, error)

	// NameTransformer is used before setEvent operation for convert name
	NameTransformer func(name interface{}) (string, error)
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

	return nil, ErrUnableToCreateStateKey
}

func ConvertName(name interface{}) (string, error) {
	switch name.(type) {
	case Namer:
		return name.(Namer).Name()
	case string:
		return name.(string), nil
	}

	return ``, ErrUnableToCreateEventName
}
