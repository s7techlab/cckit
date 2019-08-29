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
		smm.Add(pkeySchema, StateNamespace(SchemaNamespace(sm.schema)), PKeyAttr(attrs...), KeyerFor(sm.schema))
	}
}

func PKeyAttr(attrs ...string) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.primaryKeyer = attrsPKeyer(attrs)
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
		sm.primaryKeyer = attrsPKeyer([]string{`Id`})
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

func attrsFrom(schema interface{}) (attrs []string) {
	// fields from schema
	s := reflect.ValueOf(schema).Elem().Type()
	for i := 0; i < s.NumField(); i++ {

		name := s.Field(i).Name
		if strings.HasPrefix(name, `XXX_`) {
			continue
		}

		attrs = append(attrs, name)
	}
	return
}

func attrsPKeyer(attrs []string) InstanceKeyer {
	return func(instance interface{}) (key state.Key, err error) {
		inst := reflect.Indirect(reflect.ValueOf(instance))
		var pkey state.Key
		for _, attr := range attrs {
			v := inst.FieldByName(attr)
			if !v.IsValid() {
				return nil, fmt.Errorf(`%s: %s`, ErrFieldNotExists, attr)
			}

			if key, err = keyFromValue(v); err != nil {
				return nil, fmt.Errorf(`key from field %s.%s: %s`, mapKey(instance), attr, err)
			}
			pkey = pkey.Append(key)
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

	if v.Kind() == reflect.Ptr {
		s := reflect.ValueOf(v.Interface()).Elem().Type()
		// get all field values from struct
		for i := 0; i < s.NumField(); i++ {
			if !strings.HasPrefix(s.Field(i).Name, `XXX_`) {
				key = append(key, reflect.Indirect(v).Field(i).String())
			}
		}
		return key, nil
	}

	return nil, ErrFieldTypeNotSupportedForKeyExtraction
}
