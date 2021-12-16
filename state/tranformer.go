package state

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"

	"github.com/s7techlab/cckit/convert"
)

type (
	// FromBytesTransformer is used after getState operation for convert value
	FromBytesTransformer func(bb []byte, config ...interface{}) (interface{}, error)

	// ToBytesTransformer is used before putState operation for convert payload
	ToBytesTransformer func(v interface{}, config ...interface{}) ([]byte, error)

	// KeyTransformer is used before putState operation for convert key
	KeyTransformer func(Key) (Key, error)

	// StringTransformer is used before setEvent operation for convert name
	StringTransformer func(string) (string, error)

	Serializer interface {
		ToBytes(interface{}) ([]byte, error)
		FromBytes(serialized []byte, target interface{}) (interface{}, error)
	}

	ProtoSerializer struct {
	}

	JSONSerializer struct {
	}
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

// KeyAsIs returns string parts of composite key
func KeyAsIs(key Key) (Key, error) {
	return key, nil
}

func NameAsIs(name string) (string, error) {
	return name, nil
}

func (ps *ProtoSerializer) ToBytes(entry interface{}) ([]byte, error) {
	return proto.Marshal(entry.(proto.Message))
}

func (ps *ProtoSerializer) FromBytes(serialized []byte, target interface{}) (interface{}, error) {
	return convert.FromBytes(serialized, target)
}

func (js *JSONSerializer) ToBytes(entry interface{}) ([]byte, error) {
	return json.Marshal(entry)
}

func (js *JSONSerializer) FromBytes(serialized []byte, target interface{}) (interface{}, error) {
	return convert.JSONUnmarshalPtr(serialized, target)
}
