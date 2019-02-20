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
	Namer func(entity interface{}) string

	EventMapping struct {
		schema interface{}
		name   string
	}

	EventMappings map[string]*EventMapping

	EventMappers interface {
		Exists(schema interface{}) (exists bool)
		Map(schema interface{}) (keyValue state.KeyValue, err error)
		Get(schema interface{}) (eventMapper EventMapper, err error)
	}

	EventMapper interface {
		Schema() interface{}
		Name(instance interface{}) (string, error)
	}

	EventMappingOpt func(*EventMapping)
)

func (emm EventMappings) Add(schema interface{}, opts ...EventMappingOpt) EventMappings {
	em := &EventMapping{
		schema: schema,
	}

	for _, opt := range opts {
		opt(em)
	}

	applyEventMappingDefaults(em)
	emm[mapKey(schema)] = em
	return emm
}

func applyEventMappingDefaults(em *EventMapping) {
	// default namespace based on type names
	if len(em.name) == 0 {
		t := reflect.TypeOf(em.schema).String()
		em.name = t[strings.Index(t, `.`)+1:]
	}
}

func (emm EventMappings) Get(entry interface{}) (EventMapper, error) {
	m, ok := emm[mapKey(entry)]
	if !ok {
		return nil, fmt.Errorf(`%s: %s`, ErrEventMappingNotFound, mapKey(entry))
	}
	return m, nil
}

func (emm EventMappings) Exists(entry interface{}) bool {
	_, err := emm.Get(entry)
	return err == nil
}

func (emm EventMappings) Map(entry interface{}) (mapped state.NameValue, err error) {
	mapping, err := emm.Get(entry)
	if err != nil {
		return nil, errors.Wrap(err, `mapping`)
	}

	switch entry.(type) {
	case proto.Message:
		return NewProtoEventMapper(entry, mapping)
	default:
		return nil, ErrEntryTypeNotSupported
	}
}

func (em EventMapping) Schema() interface{} {
	return em.schema
}

func (em EventMapping) Name(instance interface{}) (string, error) {
	return em.name, nil
}
