package mapping

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/state"
)

// StateNamespace sets namespace for mapping
func StateNamespace(namespace state.Key) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.namespace = namespace
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

// PKeySchema registers all fields from pkeySchema as part of primary key
// also register keyer for pkeySchema with with namespace from current schema
func PKeySchema(pkeySchema interface{}) StateMappingOpt {
	attrs := attrsFrom(pkeySchema)

	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = attrsKeyer(attrs)

		//add mapping namespace for id schema same as schema
		smm.Add(pkeySchema, StateNamespace(SchemaNamespace(sm.schema)), PKeyAttr(attrs...), KeyerFor(sm.schema))
	}
}

func PKeyAttr(attrs ...string) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = attrsKeyer(attrs)
	}
}

// PKeyId use Id attr as source for mapped state entry key
func PKeyId() StateMappingOpt {
	return PKeyAttr(`Id`)
}

// PKeyConst use constant as state entry key
func PKeyConst(key state.Key) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = func(instance interface{}) (state.Key, error) {
			return key, nil
		}
	}
}

// PKeyComplexId sets Id as key field, also adds mapping for pkeySchema
// with namespace from mapping schema
func PKeyComplexId(pkeySchema interface{}) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = attrsKeyer([]string{`Id`})
		smm.Add(pkeySchema,
			StateNamespace(SchemaNamespace(sm.schema)),
			PKeyAttr(attrsFrom(pkeySchema)...),
			KeyerFor(sm.schema))
	}
}

func PKeyer(pkeyer InstanceKeyer) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = pkeyer
	}
}

func skipField(field reflect.Value) bool {
	if strings.HasPrefix(field.Type().Name(), `XXX_`) || !field.CanSet() {
		return true
	}
	return false
}

func attrsFrom(schema interface{}) (attrs []string) {
	// fields from schema
	s := reflect.ValueOf(schema).Elem()
	fs := s.Type()
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		if skipField(field) {
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

// attrMultiKeyer creates keyer based of one field and can return multiple keyss
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

func keyFromValue(v reflect.Value) (state.Key, error) {
	var key state.Key

	if v.Kind() == reflect.Ptr {
		s := reflect.ValueOf(v.Interface()).Elem()
		// get all field values from struct
		for i := 0; i < s.NumField(); i++ {
			field := s.Field(i)
			if skipField(field) {
				continue
			} else {
				key = append(key, reflect.Indirect(v).Field(i).String())
			}
		}

		return key, nil
	}

	switch v.Type().String() {

	case `string`, `int32`, `bool`:
		// multi key possible
		key = state.Key{v.String()}

	case `[]string`:
		// every slice element is a part of one key
		for i := 0; i < v.Len(); i++ {
			key = append(key, v.Index(i).String())
		}

	default:
		return nil, ErrFieldTypeNotSupportedForKeyExtraction
	}

	return key, nil
}
