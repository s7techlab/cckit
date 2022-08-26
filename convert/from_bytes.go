package convert

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// FromBytes converts []byte to target interface
func FromBytes(bb []byte, target interface{}) (result interface{}, err error) {
	// create copy

	switch t := target.(type) {
	case string:
		return string(bb), nil
	case []byte:
		return bb, nil
	case int:
		return strconv.Atoi(string(bb))
	case bool:
		return strconv.ParseBool(string(bb))
	case []string:
		arrInterface, err := JSONUnmarshalPtr(bb, &target)

		if err != nil {
			return nil, err
		}
		var arrString []string
		for _, v := range arrInterface.([]interface{}) {
			arrString = append(arrString, v.(string))
		}
		return arrString, nil

	case FromByter:
		return t.FromBytes(bb)

	case proto.Message:
		return ProtoUnmarshal(bb, t)

	default:
		return FromBytesToStruct(bb, target)
	}

}

// FromBytesToStruct converts []byte to struct,array,slice depending on target type
func FromBytesToStruct(bb []byte, target interface{}) (result interface{}, err error) {
	if bb == nil {
		return nil, ErrUnableToConvertNilToStruct
	}
	targetType := reflect.TypeOf(target).Kind()

	switch targetType {
	case reflect.Struct:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		// will be map[string]interface{}
		return JSONUnmarshalPtr(bb, &target)
	case reflect.Ptr:
		return JSONUnmarshalPtr(bb, target)

	default:
		return nil, fmt.Errorf(
			`fromBytes converting supports ToByter interface,struct,array,slice and string, current type is %s`,
			targetType)
	}
}

// JSONUnmarshalPtr unmarshalls []byte as json to pointer, and returns value pointed to
func JSONUnmarshalPtr(bb []byte, to interface{}) (result interface{}, err error) {
	targetPtr := reflect.New(reflect.ValueOf(to).Elem().Type()).Interface()
	err = json.Unmarshal(bb, targetPtr)
	if err != nil {
		return nil, ErrUnableToConvertValueToStruct
	}
	return reflect.Indirect(reflect.ValueOf(targetPtr)).Interface(), nil
}

// ProtoUnmarshal r unmarshalls []byte as proto.Message to pointer, and returns value pointed to
func ProtoUnmarshal(bb []byte, messageType proto.Message) (message proto.Message, err error) {
	msg := proto.Clone(messageType)
	err = proto.Unmarshal(bb, msg)
	if err != nil {
		return nil, ErrUnableToConvertValueToStruct
	}
	return msg, nil
}

// FromResponse converts response.Payload to target
func FromResponse(response peer.Response, target interface{}) (result interface{}, err error) {
	if response.Status == shim.ERROR {
		return nil, errors.New(response.Message)
	}
	return FromBytes(response.Payload, target)
}
