package state

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/convert"
)

type (
	EntryMapper interface {
		Keyer
		convert.ToByter
	}
	Mapping struct {
		Prefix string
		Type   EntryMappings
	}

	EntryMapping struct {
		PrimaryKey           KeyTransformer
		PrimaryKeyFromString KeyerFunc
		PrimaryKeyToString   KeyerFunc
		UniqKey              []KeyTransformer
		Key                  []KeyTransformer
	}

	EntryMappings map[string]EntryMapping

	ProtoMapper struct {
		schema interface{}
		keyer  KeyFunc
	}
)

func (em Mapping) Exists(entry interface{}) bool {
	_, ok := em.Type[reflect.TypeOf(entry).String()]
	return ok
}

func (em Mapping) Apply(entry interface{}) (mapper EntryMapper, err error) {

	entryType := reflect.TypeOf(entry).String()
	mapping, ok := em.Type[entryType]
	if !ok {
		return nil, fmt.Errorf(`%s: %s`, ErrEntryMappingNotDefined, entryType)
	}

	switch entry.(type) {
	case proto.Message:
		return NewProtoMapper(entry, mapping)
	default:
		return nil, ErrEntryTypeMappingNotSupported
	}
}

func NewProtoMapper(schema interface{}, mapping EntryMapping) (*ProtoMapper, error) {

	return &ProtoMapper{schema, func() ([]string, error) {
		return mapping.PrimaryKey(schema)
	}}, nil
}

func (mp *ProtoMapper) Key() ([]string, error) {
	return mp.keyer()
}

func (mp *ProtoMapper) ToBytes() ([]byte, error) {
	return proto.Marshal(mp.schema.(proto.Message))
}
