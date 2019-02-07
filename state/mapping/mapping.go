package mapping

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/state"
)

type (
	// Mappings interface for mapping collection
	Mappings interface {
		Exists(entity interface{}) (exists bool)
		Map(entity interface{}) (entry state.KeyValue, err error)
		Get(entity interface{}) (mapping Mapping, err error)
	}

	Mapping interface {
		Schema() interface{}
		Namespace() []string
		PrimaryKey(instance interface{}) ([]string, error)
	}

	SchemaMapping struct {
		schema       interface{}
		namespace    []string
		primaryKeyer state.KeyTransformer
		//PKStringer KeyerFunc
		//PKToString KeyerFunc
		//niqKey []KeyTransformer
		//Key     []KeyTransformer
	}

	SchemaMappings map[string]*SchemaMapping
)

func mapKey(entry interface{}) string {
	return reflect.TypeOf(entry).String()
}

func (smm SchemaMappings) Add(schema interface{}, namespace []string, primaryKeyer state.KeyTransformer) SchemaMappings {
	smm[mapKey(schema)] = &SchemaMapping{
		schema:       schema,
		namespace:    namespace,
		primaryKeyer: primaryKeyer,
	}
	return smm
}

func (smm SchemaMappings) Get(entry interface{}) (Mapping, error) {
	m, ok := smm[mapKey(entry)]
	if !ok {
		return nil, fmt.Errorf(`%s: %s`, ErrEntryTypeNotDefined, mapKey(entry))
	}
	return m, nil
}

func (smm SchemaMappings) Exists(entry interface{}) bool {
	_, err := smm.Get(entry)
	return err == nil
}

func (smm SchemaMappings) Map(entry interface{}) (mapped state.KeyValue, err error) {
	mapping, err := smm.Get(entry)
	if err != nil {
		return nil, errors.Wrap(err, `mapping`)
	}

	switch entry.(type) {
	case proto.Message:
		return NewProtoMapper(entry, mapping)
	default:
		return nil, ErrEntryTypeNotSupported
	}
}

func (sm *SchemaMapping) Namespace() []string {
	return sm.namespace
}
func (sm *SchemaMapping) Schema() interface{} {
	return sm.schema
}

func (sm *SchemaMapping) PrimaryKey(entity interface{}) ([]string, error) {
	key, err := sm.primaryKeyer(entity)
	if err != nil {
		return nil, err
	}
	return append(sm.namespace, key...), nil
}
