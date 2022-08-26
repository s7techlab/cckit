package mapping

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/s7techlab/cckit/state"
)

const (
	TimestampKeyLayout = `2006-01-02`
)

// WithNamespace sets namespace for mapping
func WithNamespace(namespace state.Key) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.namespace = namespace
	}
}

// WithConstPKey set static key for all instances of mapped entry
func WithConstPKey(keys ...state.Key) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		key := state.Key{}
		for _, k := range keys {
			key = key.Append(k)
		}

		sm.primaryKeyer = func(_ interface{}) (state.Key, error) {
			return key, nil
		}
	}
}

func KeyerFor(schema interface{}) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.keyerForSchema = schema
	}
}

// List defined list container, it must have `Items` attr
func List(list proto.Message) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.list = list
	}
}

// UniqKey defined uniq key in entity
func UniqKey(name string, fields ...[]string) StateMappingOpt {
	var ff []string
	if len(fields) > 0 {
		ff = fields[0]
	}
	return WithIndex(&StateIndexDef{
		Name:     name,
		Fields:   ff,
		Required: true,
		Multi:    false,
	})
}

func WithIndex(idx *StateIndexDef) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		if idx.Name == `` {
			return
		}

		var keyer InstanceMultiKeyer
		if idx.Keyer != nil {
			keyer = idx.Keyer
		} else {
			aa := []string{idx.Name}
			if len(idx.Fields) > 0 {
				aa = idx.Fields
			}

			// multiple external ids refers to one entry
			if idx.Multi {
				keyer = attrMultiKeyer(aa[0])
			} else {
				keyer = keyerAsMulti(attrsKeyer(aa))
			}
		}

		_ = sm.AddIndex(&StateIndex{
			Name:     idx.Name,
			Uniq:     true,
			Required: idx.Required,
			Keyer:    keyer,
		})
	}
}

// PKeySchema registers all fields from pkeySchema as part of primary key.
// Same fields should exist in mapped entity.
// Also register keyer for pkeySchema with namespace from current schema.
func PKeySchema(pkeySchema interface{}) StateMappingOpt {
	attrs := attrsFrom(pkeySchema)

	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = attrsKeyer(attrs)

		// inherit namespace from "parent" mapping
		namespace := sm.namespace
		if len(namespace) == 0 {
			namespace = sm.DefaultNamespace()
		}

		//add mapping for schema identifier
		smm.Add(
			pkeySchema,
			WithNamespace(namespace),
			PKeyAttr(attrs...),
			KeyerFor(sm.schema))
	}
}

func PKeyAttr(attrs ...string) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = attrsKeyer(attrs)
	}
}

// PKeyId use ID attr as source for mapped state entry key
func PKeyId() StateMappingOpt {
	return PKeyAttr(`Id`)
}

// PKeyComplexId sets ID as key field, also adds mapping for pkeySchema
// with namespace from mapping schema
func PKeyComplexId(pkeySchema interface{}) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = attrsKeyer([]string{`Id`})
		smm.Add(pkeySchema,
			WithNamespace(SchemaNamespace(sm.schema)),
			PKeyAttr(attrsFrom(pkeySchema)...),
			KeyerFor(sm.schema))
	}
}

func PKeyer(pKeyer InstanceKeyer) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = pKeyer
	}
}

func skipField(name string, field reflect.Value) bool {
	if strings.HasPrefix(name, `XXX_`) || !field.CanSet() {
		return true
	}
	return false
}

// attrFrom extracts list of field names from struct
func attrsFrom(schema interface{}) (attrs []string) {
	// fields from schema
	s := reflect.ValueOf(schema).Elem()
	fs := s.Type()
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		if skipField(fs.Field(i).Name, field) {
			continue
		}
		attrs = append(attrs, fs.Field(i).Name)
	}
	return
}

// attrsKeyer creates instance keyer
func attrsKeyer(attrs []string) InstanceKeyer {
	return func(instance interface{}) (state.Key, error) {
		var key = state.Key{}
		inst := reflect.Indirect(reflect.ValueOf(instance))

		for _, attr := range attrs {

			v := inst.FieldByName(attr)
			if !v.IsValid() {
				return nil, fmt.Errorf(`%s: %s`, ErrFieldNotExists, attr)
			}

			keyPart, err := keyFromValue(v)
			if err != nil {
				return nil, fmt.Errorf(`key from field %s.%s: %s`, mapKey(instance), attr, err)
			}
			key = key.Append(keyPart)
		}
		return key, nil
	}
}

// attrMultiKeyer creates keyer based of one field and can return multiple keys
func attrMultiKeyer(attr string) InstanceMultiKeyer {
	return func(instance interface{}) ([]state.Key, error) {
		inst := reflect.Indirect(reflect.ValueOf(instance))

		v := inst.FieldByName(attr)
		if !v.IsValid() {
			return nil, fmt.Errorf(`%s: %s`, ErrFieldNotExists, attr)
		}

		return keysFromValue(v)
	}
}

// keyerAsMulti adapter keyer to multiKeyer
func keyerAsMulti(keyer InstanceKeyer) InstanceMultiKeyer {
	return func(instance interface{}) (key []state.Key, err error) {
		k, err := keyer(instance)
		if err != nil {
			return nil, err
		}

		return []state.Key{k}, nil
	}
}

// multi - returns multiple key if value type allows it
func keysFromValue(v reflect.Value) ([]state.Key, error) {
	var keys []state.Key

	switch v.Type().String() {
	case `[]string`:
		for i := 0; i < v.Len(); i++ {
			keys = append(keys, state.Key{v.Index(i).String()})
		}

	default:
		return nil, ErrFieldTypeNotSupportedForKeyExtraction
	}

	return keys, nil
}

// keyFromValue creates string representation of value for state key
func keyFromValue(v reflect.Value) (state.Key, error) {
	switch v.Kind() {

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return state.Key{strconv.Itoa(int(v.Uint()))}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// if it is enum in protobuf
		if stringer, ok := v.Interface().(fmt.Stringer); ok {
			return state.Key{stringer.String()}, nil
		}

		return state.Key{strconv.Itoa(int(v.Int()))}, nil

	case reflect.Ptr:
		// todo: extract key producer and add custom serializers
		switch val := v.Interface().(type) {

		case *timestamp.Timestamp:
			t, err := ptypes.Timestamp(val)
			if err != nil {
				return nil, fmt.Errorf(`timestamp key to time: %w`, err)
			}
			return state.Key{t.Format(TimestampKeyLayout)}, nil

		default:
			key := state.Key{}
			s := reflect.ValueOf(v.Interface()).Elem()
			fs := s.Type()
			// get all field values from struct
			for i := 0; i < s.NumField(); i++ {
				field := s.Field(i)
				if skipField(fs.Field(i).Name, field) {
					continue
				} else {
					subKey, err := keyFromValue(reflect.Indirect(v).Field(i))
					if err != nil {
						return nil, fmt.Errorf(`sub key=%s: %w`, fs.Field(i).Name, err)
					}
					key = key.Append(subKey)
				}
			}

			return key, nil
		}
	}

	switch v.Type().String() {

	case `string`, `int32`, `uint32`, `bool`:
		// multi key possible
		return state.Key{v.String()}, nil

	case `[]string`:
		key := state.Key{}
		// every slice element is a part of one key
		for i := 0; i < v.Len(); i++ {
			key = append(key, v.Index(i).String())
		}
		return key, nil

	default:
		return nil, ErrFieldTypeNotSupportedForKeyExtraction
	}
}
