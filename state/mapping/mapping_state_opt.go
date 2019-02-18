package mapping

import (
	"fmt"
	"reflect"

	"github.com/s7techlab/cckit/state"
)

func StateNamespace(namespace state.Key) StateMappingOpt {
	return func(sm *StateMapping, smm StateMappings) {
		sm.namespace = namespace
	}
}

func PKeySchema(pkeySchema interface{}) StateMappingOpt {

	var attrs []string
	// fields from pkey schema
	s := reflect.ValueOf(pkeySchema).Elem().Type()
	for i := 0; i < s.NumField(); i++ {
		attrs = append(attrs, s.Field(i).Name)
	}

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

func attrsPKeyer(attrs []string) InstanceKeyer {
	return func(instance interface{}) (state.Key, error) {
		r := reflect.ValueOf(instance)
		var k state.Key
		for _, attr := range attrs {
			f := reflect.Indirect(r).FieldByName(attr)

			if !f.IsValid() {
				return nil, fmt.Errorf(`%s: %s`, ErrFieldNotExists, attr)
			}
			k = append(k, f.String())
		}
		return k, nil
	}
}
