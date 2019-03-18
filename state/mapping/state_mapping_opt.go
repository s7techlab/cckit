package mapping

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/state"
)

// StateNamespace sets namespace for mapping
func StateNamespace(namespace state.Key) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.namespace = namespace
	}
}

// List defined list container, it must have `Items` attr
func List(list proto.Message) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.list = list
	}
}

func UniqKey(name string, attrs ...[]string) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		aa := []string{name}
		if len(attrs) > 0 {
			aa = attrs[0]
		}
		sm.uniqKeys = append(sm.uniqKeys, &StateKeyDefinition{Name: name, Attrs: aa})
	}
}

// PKeySchema registers all fields from pkeySchema as part of primary key
// also register keyer for pkeySchema with with namespace from current schema
func PKeySchema(pkeySchema interface{}) StateMappingOpt {
	attrs := attrsFrom(pkeySchema)

	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = attrsPKeyer(attrs)

		//add mapping namespace for id schema same as schema
		smm.Add(pkeySchema, StateNamespace(schemaNamespace(sm.schema)), PKeyAttr(attrs...))
	}
}

func PKeyAttr(attrs ...string) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = attrsPKeyer(attrs)
	}
}

func PKeyId() StateMappingOpt {
	return PKeyAttr(`Id`)
}

func PKeyComplexId(pkeySchema interface{}) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = attrsPKeyer([]string{`Id`})
		smm.Add(pkeySchema, StateNamespace(schemaNamespace(sm.schema)), PKeyAttr(attrsFrom(pkeySchema)...))
	}
}

func PKeyer(pkeyer InstanceKeyer) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = pkeyer
	}
}

func attrsFrom(schema interface{}) (attrs []string) {
	// fields from schema
	s := reflect.ValueOf(schema).Elem().Type()
	for i := 0; i < s.NumField(); i++ {
		attrs = append(attrs, s.Field(i).Name)
	}
	return
}

func attrsPKeyer(attrs []string) InstanceKeyer {
	return func(instance interface{}) (state.Key, error) {
		inst := reflect.Indirect(reflect.ValueOf(instance))
		var pkey state.Key
		for _, attr := range attrs {
			v := inst.FieldByName(attr)
			if !v.IsValid() {
				return nil, fmt.Errorf(`%s: %s`, ErrFieldNotExists, attr)
			}

			if key, err := keyFromValue(v); err != nil {
				return nil, fmt.Errorf(`key from field %s.%s: %s`, mapKey(instance), attr, err)
			} else {
				pkey = pkey.Append(key)
			}
		}
		return pkey, nil
	}
}

func keyFromValue(v reflect.Value) (key state.Key, err error) {
	switch v.Type().String() {
	case `string`, `int32`, `bool`:
		return state.Key{v.String()}, nil
	case `[]string`:
		for i := 0; i < v.Len(); i++ {
			key = append(key, v.Index(i).String())
		}
		return key, nil
	}

	switch v.Kind() {
	case reflect.Ptr:
		s := reflect.ValueOf(v.Interface()).Elem().Type()
		// get all field values from struct
		for i := 0; i < s.NumField(); i++ {
			key = append(key, reflect.Indirect(v).Field(i).String())
		}
		return key, nil
	}

	return nil, ErrFieldTypeNotSupportedForKeyExtraction
}
