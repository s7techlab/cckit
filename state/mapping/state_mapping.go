package mapping

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/state"
)

type (
	// StateMappers interface for mappers collection
	StateMappers interface {
		Exists(schema interface{}) (exists bool)
		Map(schema interface{}) (keyValue state.KeyValue, err error)
		Get(schema interface{}) (stateMapper StateMapper, err error)
		PrimaryKey(schema interface{}) (key state.Key, err error)
	}

	// StateMapper
	StateMapper interface {
		Schema() interface{}
		List() interface{}
		Namespace() state.Key
		PrimaryKey(instance interface{}) (key state.Key, err error)
		Keys(instance interface{}) (key []state.KeyValue, err error)
	}

	InstanceKeyer func(instance interface{}) (key state.Key, err error)

	StateMapped interface {
		state.KeyValue // entry key and value
		Mapper() StateMapper
		Keys() ([]state.KeyValue, error)
	}
	// StateMapping defines metadata for mapping from schema to state keys/values
	StateMapping struct {
		schema       interface{}
		namespace    state.Key
		primaryKeyer InstanceKeyer
		list         interface{}
		uniqKeys     []*StateKeyDefinition
	}

	StateKeyDefinition struct {
		Name  string
		Attrs []string
	}

	StateMappings map[string]*StateMapping

	StateMappingOptions struct {
		namespace    state.Key
		primaryKeyer state.KeyTransformer
	}

	StateMappingOpt func(*StateMapping, StateMappings)
)

func mapKey(entry interface{}) string {
	return reflect.TypeOf(entry).String()
}

func (smm StateMappings) Add(schema interface{}, opts ...StateMappingOpt) StateMappings {
	sm := &StateMapping{
		schema: schema,
	}

	for _, opt := range opts {
		opt(sm, smm)
	}

	applyStateMappingDefaults(sm)
	smm[mapKey(schema)] = sm
	return smm
}

func applyStateMappingDefaults(sm *StateMapping) {
	// default namespace based on type name
	if len(sm.namespace) == 0 {
		sm.namespace = schemaNamespace(sm.schema)
	}
}

func schemaNamespace(schema interface{}) state.Key {
	t := reflect.TypeOf(schema).String()
	return state.Key{t[strings.Index(t, `.`)+1:]}
}

func (smm StateMappings) Get(entry interface{}) (StateMapper, error) {
	m, ok := smm[mapKey(entry)]
	if !ok {
		return nil, fmt.Errorf(`%s: %s`, ErrStateMappingNotFound, mapKey(entry))
	}
	return m, nil
}

func (smm StateMappings) Exists(entry interface{}) bool {
	_, err := smm.Get(entry)
	return err == nil
}

func (smm StateMappings) PrimaryKey(entry interface{}) (pkey state.Key, err error) {
	var m StateMapper
	if m, err = smm.Get(entry); err != nil {
		return nil, err
	}
	return m.PrimaryKey(entry)
}

func (smm StateMappings) Map(entry interface{}) (mapped StateMapped, err error) {
	mapper, err := smm.Get(entry)
	if err != nil {
		return nil, errors.Wrap(err, `mapping`)
	}

	switch entry.(type) {
	case proto.Message:
		return NewProtoStateMapped(entry, mapper), nil
	default:
		return nil, ErrEntryTypeNotSupported
	}
}

func (sm *StateMapping) Namespace() state.Key {
	return sm.namespace
}
func (sm *StateMapping) Schema() interface{} {
	return sm.schema
}

func (sm *StateMapping) List() interface{} {
	return sm.list
}

func (sm *StateMapping) PrimaryKey(entity interface{}) (state.Key, error) {
	if sm.primaryKeyer == nil {
		return nil, fmt.Errorf(`%s: schema "%s", namespace : "%s"`, ErrPrimaryKeyerNotDefined, sm.schema, sm.namespace)
	}
	key, err := sm.primaryKeyer(entity)
	if err != nil {
		return nil, err
	}
	return append(sm.namespace, key...), nil
}

func (sm *StateMapping) Keys(entity interface{}) (kv []state.KeyValue, err error) {
	if len(sm.uniqKeys) == 0 {
		return
	}

	pk, err := sm.PrimaryKey(entity)
	if err != nil {
		return
	}

	for _, k := range sm.uniqKeys {
		// uniq key attr values
		refKey, err := attrsPKeyer(k.Attrs)(entity)
		if err != nil {
			return nil, fmt.Errorf(`uniq key %s: %s`, k.Name, err)
		}

		// key will be <`_idx`,{SchemaName},{idxName}, {Key[1]},... {Key[n}}>s
		kv = append(kv, NewKeyRefMapped(sm.schema, k.Name, refKey, pk))
	}

	return
}
