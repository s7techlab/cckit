package state

import (
	"fmt"
	"reflect"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/s7techlab/cckit/state/schema"
)

type (
	StateList struct {
		itemTarget interface{}
		listTarget interface{}
		list       []interface{}
	}
)

func NewStateList(config ...interface{}) (sl *StateList, err error) {
	var (
		itemTarget, listTarget interface{}
	)
	if len(config) > 0 {
		itemTarget = config[0]
	}
	if len(config) > 1 {
		listTarget = config[1]
	}

	return &StateList{itemTarget: itemTarget, listTarget: listTarget}, nil
}

func (sl *StateList) Fill(iter shim.StateQueryIteratorInterface, fromBytes FromBytesTransformer) (list interface{}, err error) {

	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return nil, err
		}
		item, err := fromBytes(kv.Value, sl.itemTarget)
		if err != nil {
			return nil, fmt.Errorf(`transform list entry: %w`, err)
		}
		sl.list = append(sl.list, item)
	}
	return sl.Get()
}

func (sl *StateList) Get() (list interface{}, err error) {

	// custom list proto.Message
	if _, isListProto := sl.listTarget.(proto.Message); isListProto {

		customList := proto.Clone(sl.listTarget.(proto.Message))
		items := reflect.ValueOf(customList).Elem().FieldByName(`Items`)
		for _, v := range sl.list {
			items.Set(reflect.Append(items, reflect.ValueOf(v)))
		}
		return customList, nil

		// default list proto.Message (with repeated Any)
	} else if _, isItemProto := sl.itemTarget.(proto.Message); isItemProto {
		defList := &schema.List{}

		for _, item := range sl.list {
			newAny, err := anypb.New(item.(proto.Message))
			if err != nil {
				return nil, err
			}
			defList.Items = append(defList.Items, newAny)
		}
		return defList, nil
	}

	return sl.list, nil
}

func (sl *StateList) AddElementToList(elem interface{}) {
	sl.list = append(sl.list, elem)
}
