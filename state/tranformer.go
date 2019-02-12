package state

import "github.com/s7techlab/cckit/convert"

type (
	// ToBytesTransformer is used after getState operation for convert value
	FromBytesTransformer func(bb []byte, config ...interface{}) (interface{}, error)

	// ToBytesTransformer is used before putState operation for convert payload
	ToBytesTransformer func(v interface{}, config ...interface{}) ([]byte, error)

	// KeyTransformer is used before putState operation for convert key
	KeyTransformer func(key []string) ([]string, error)

	// NameTransformer is used before setEvent operation for convert name
	NameTransformer func(name string) (string, error)
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
func KeyAsIs(key []string) ([]string, error) {
	return key, nil
}

func NameAsIs(name string) (string, error) {
	return name, nil
}
