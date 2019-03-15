package mapping

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/state"
)

var (
	ErrFieldNotExists         = errors.New(`field is not exists`)
	ErrPrimaryKeyerNotDefined = errors.New(`primary keyer is not defined`)
)

type (
	// StateMappers interface for mappers collection
	StateMappers interface {
		Exists(schema interface{}) (exists bool)
		Map(schema interface{}) (keyValue state.KeyValue, err error)
		Get(schema interface{}) (stateMapper StateMapper, err error)
		PrimaryKey(schema interface{}) (key state.Key, err error)
	}

	StateMapper interface {
		Schema() interface{}
		List() interface{}
		Namespace() state.Key
		PrimaryKey(instance interface{}) (key state.Key, err error)
	}

	InstanceKeyer func(instance interface{}) (key state.Key, err error)

	StateMapping struct {
		schema       interface{}
		namespace    state.Key
		primaryKeyer InstanceKeyer
		list         interface{}
		//niqKey []KeyTransformer
		//Key     []KeyTransformer
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

func (smm StateMappings) Map(entry interface{}) (mapped state.KeyValue, err error) {
	mapping, err := smm.Get(entry)
	if err != nil {
		return nil, errors.Wrap(err, `mapping`)
	}

	switch entry.(type) {
	case proto.Message:
		return NewProtoStateMapper(entry, mapping)
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

func PKeyer(pkeyer InstanceKeyer) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = pkeyer
	}
}
