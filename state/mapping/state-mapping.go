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
	}

	StateMapper interface {
		Schema() interface{}
		Namespace() state.Key
		PrimaryKey(instance interface{}) (key []string, err error)
	}

	InstanceKeyer func(instance interface{}) (key []string, err error)

	StateMapping struct {
		schema       interface{}
		namespace    state.Key
		primaryKeyer InstanceKeyer
		//PKStringer KeyerFunc
		//PKToString KeyerFunc
		//niqKey []KeyTransformer
		//Key     []KeyTransformer
	}

	StateMappings map[string]*StateMapping

	StateMappingOptions struct {
		namespace    []string
		primaryKeyer state.KeyTransformer
	}

	StateMappingOpt func(*StateMapping)
)

func mapKey(entry interface{}) string {
	return reflect.TypeOf(entry).String()
}

func (smm StateMappings) Add(schema interface{}, opts ...StateMappingOpt) StateMappings {
	sm := &StateMapping{
		schema: schema,
	}

	for _, opt := range opts {
		opt(sm)
	}

	applyStateMappingDefaults(sm)
	smm[mapKey(schema)] = sm
	return smm
}

func applyStateMappingDefaults(sm *StateMapping) {
	// default namespace based on type name
	if len(sm.namespace) == 0 {
		t := reflect.TypeOf(sm.schema).String()
		sm.namespace = []string{t[strings.Index(t, `.`)+1:]}
	}
}

func (smm StateMappings) Get(entry interface{}) (StateMapper, error) {
	m, ok := smm[mapKey(entry)]
	if !ok {
		return nil, fmt.Errorf(`%s: %s`, ErrEntryTypeNotDefined, mapKey(entry))
	}
	return m, nil
}

func (smm StateMappings) Exists(entry interface{}) bool {
	_, err := smm.Get(entry)
	return err == nil
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

func (s *StateMapping) PrimaryKey(entity interface{}) ([]string, error) {
	key, err := s.primaryKeyer(entity)
	if err != nil {
		return nil, err
	}
	return append(s.namespace, key...), nil
}

func UseStatePKeyer(pkeyer InstanceKeyer) StateMappingOpt {
	return func(sm *StateMapping) {
		sm.primaryKeyer = pkeyer
	}
}
func UseStateNamespace(namespace state.Key) StateMappingOpt {
	return func(sm *StateMapping) {
		sm.namespace = namespace
	}
}
