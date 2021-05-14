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
		// PrimaryKey returns primary key for entry
		PrimaryKey(instance interface{}) (key state.Key, err error)
		// Keys returns additional keys for
		Keys(instance interface{}) (key []state.KeyValue, err error)
		//KeyerFor returns target entity if mapper is key mapper
		KeyerFor() (schema interface{})
		Indexes() []*StateIndex
	}

	// InstanceKeyer returns key of an state entry instance
	InstanceKeyer      func(instance interface{}) (state.Key, error)
	InstanceMultiKeyer func(instance interface{}) ([]state.Key, error)

	StateMapped interface {
		state.KeyValue // entry key and value
		Mapper() StateMapper
		Keys() ([]state.KeyValue, error)
	}
	// StateMapping defines metadata for mapping from schema to state keys/values
	StateMapping struct {
		schema         interface{}
		namespace      state.Key     // prefix for primary key
		keyerForSchema interface{}   // schema is keyer for another schema ( for example *schema.StaffId for *schema.Staff )
		primaryKeyer   InstanceKeyer // primary key always one
		list           interface{}   // list schema
		indexes        []*StateIndex // additional keys
	}

	// StateIndex additional index of entity instance
	StateIndex struct {
		Name     string
		Uniq     bool
		Required bool
		Keyer    InstanceMultiKeyer // index can have multiple keys
	}

	StateIndexDef struct {
		Name     string
		Fields   []string
		Required bool
		Multi    bool
		Keyer    InstanceMultiKeyer
	}

	StateMappings map[string]*StateMapping

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
		sm.namespace = SchemaNamespace(sm.schema)
	}
}

func SchemaNamespace(schema interface{}) state.Key {
	t := reflect.TypeOf(schema).String()
	return state.Key{t[strings.Index(t, `.`)+1:]}
}

// Get mapper for mapped entry
func (smm StateMappings) Get(entry interface{}) (StateMapper, error) {
	switch id := entry.(type) {
	case []string:
		return smm.GetByNamespace(id[0:1])
	default:
		m, ok := smm[mapKey(entry)]
		if !ok {
			return nil, fmt.Errorf(`%s: %s`, ErrStateMappingNotFound, mapKey(entry))
		}
		return m, nil
	}
}

// Get mapper by string namespace. It can be used in block explorer: we know state key, but don't know
// type actually mapped to state
func (smm StateMappings) GetByNamespace(namespace state.Key) (StateMapper, error) {
	for _, m := range smm {
		if m.keyerForSchema == nil && reflect.DeepEqual(m.namespace, namespace) {
			return m, nil
		}
	}
	return nil, fmt.Errorf(`%s: %s`, ErrStateMappingNotFound, namespace)
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
	case proto.Message, []string:
		return NewProtoStateMapped(entry, mapper), nil
	default:
		return nil, ErrEntryTypeNotSupported
	}
}

//
func (smm *StateMappings) IdxKey(entity interface{}, idx string, idxVal state.Key) (state.Key, error) {
	keyMapped := NewKeyRefIDMapped(entity, idx, idxVal)
	return keyMapped.Key()
}

func (sm *StateMapping) Namespace() state.Key {
	return sm.namespace
}

func (sm *StateMapping) Indexes() []*StateIndex {
	return sm.indexes
}

func (sm *StateMapping) Schema() interface{} {
	return sm.schema
}

func (sm *StateMapping) List() interface{} {
	return sm.list
}

func (sm *StateMapping) PrimaryKey(entity interface{}) (state.Key, error) {
	if sm.primaryKeyer == nil {
		return nil, fmt.Errorf(`%s: schema "%s", namespace : "%s"`,
			ErrPrimaryKeyerNotDefined, sm.schema, sm.namespace)
	}
	key, err := sm.primaryKeyer(entity)
	if err != nil {
		return nil, err
	}
	return append(sm.namespace, key...), nil
}

// Indexes prepares primary and additional uniq/non-uniq keys for storage
func (sm *StateMapping) Keys(entity interface{}) ([]state.KeyValue, error) {
	if len(sm.indexes) == 0 {
		return nil, nil
	}

	pk, err := sm.PrimaryKey(entity) // primary key, all additional keys refers to primary key
	if err != nil {
		return nil, err
	}

	var stateKeys []state.KeyValue
	for _, idx := range sm.indexes {
		// uniq key attr values
		idxKeys, err := idx.Keyer(entity)
		if err != nil {
			return nil, errors.Errorf(`uniq key %s: %s`, idx.Name, err)
		}

		for _, key := range idxKeys {
			// key will be <`_idx`,{SchemaName},{idxName}, {Key[1]},... {Key[n}}>s
			stateKeys = append(stateKeys, NewKeyRefMapped(sm.schema, idx.Name, key, pk))
		}
	}

	return stateKeys, nil
}

func (sm *StateMapping) AddIndex(idx *StateIndex) error {
	if exists := sm.Index(idx.Name); exists != nil {
		return ErrIndexAlreadyExists
	}

	sm.indexes = append(sm.indexes, idx)
	return nil
}

func (sm *StateMapping) Index(name string) *StateIndex {
	for _, idx := range sm.indexes {
		if idx.Name == name {
			return idx
		}
	}

	return nil
}

func (sm *StateMapping) KeyerFor() interface{} {
	return sm.keyerForSchema
}

// KeyRefsDiff calculates diff between key reference set
func KeyRefsDiff(prevKeys []state.KeyValue, newKeys []state.KeyValue) (deleted, inserted []state.KeyValue, err error) {

	var (
		prevK = make(map[string]int)
		newK  = make(map[string]int)
	)
	for i, kv := range prevKeys {
		k, err := kv.Key()
		if err != nil {
			return nil, nil, errors.Wrap(err, `prev ref key`)
		}

		prevK[k.String()] = i
	}

	for i, kv := range newKeys {
		k, err := kv.Key()
		if err != nil {
			return nil, nil, errors.Wrap(err, `new ref key`)
		}

		newK[k.String()] = i
	}

	for k, i := range prevK {
		if _, ok := newK[k]; !ok {
			deleted = append(deleted, prevKeys[i])
		}
	}

	for k, i := range newK {
		if _, ok := prevK[k]; !ok {
			inserted = append(inserted, newKeys[i])
		}
	}

	return deleted, inserted, nil
}
