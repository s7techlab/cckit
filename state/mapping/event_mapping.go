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
	ErrEventNameNotFound = errors.New(`event name not found`)
)

type (
	Namer func(entity interface{}) string

	EventMapping struct {
		schema interface{}
		name   string
	}

	EventMappings map[string]*EventMapping

	EventMapped interface {
		state.NameValue
	}

	EventMappers interface {
		Exists(schema interface{}) (exists bool)
		Map(schema interface{}) (keyValue state.KeyValue, err error)
		Get(schema interface{}) (eventMapper EventMapper, err error)
	}

	EventMapper interface {
		Schema() interface{}
		Name(instance interface{}) (string, error)
	}

	Event struct {
		Name    string
		Payload interface{}
	}

	EventMappingOpt func(*EventMapping)

	EventResolver interface {
		Resolve(eventName string, payload []byte) (event interface{}, err error)
	}
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
		em.name = EventNameForPayload(em.schema)
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

func (emm EventMappings) Map(entry interface{}) (instance *EventInstance, err error) {
	mapping, err := emm.Get(entry)
	if err != nil {
		return nil, errors.Wrap(err, `mapping`)
	}

	switch entry.(type) {
	case proto.Message:
		return NewEventInstance(entry, mapping, DefaultSerializer)
	default:
		return nil, ErrEntryTypeNotSupported
	}
}

func (emm EventMappings) Resolve(eventName string, payload []byte) (event interface{}, err error) {
	for _, m := range emm {
		if m.name == eventName {
			return DefaultSerializer.FromBytes(payload, m.Schema())
		}
	}

	return nil, ErrEventNameNotFound
}

func (em EventMapping) Schema() interface{} {
	return em.schema
}

func (em EventMapping) Name(instance interface{}) (string, error) {
	return em.name, nil
}

func EventNameForPayload(payload interface{}) string {
	t := reflect.TypeOf(payload).String()
	return t[strings.Index(t, `.`)+1:]
}

func EventFromPayload(payload interface{}) *Event {
	return &Event{
		Name:    EventNameForPayload(payload),
		Payload: payload,
	}
}

func MergeEventMappings(one EventMappings, more ...EventMappings) EventMappings {
	out := make(EventMappings)
	for k, v := range one {
		out[k] = v
	}

	for _, m := range more {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}
