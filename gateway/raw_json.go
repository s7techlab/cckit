package gateway

import (
	"github.com/golang/protobuf/jsonpb"
)

func (x *RawJson) MarshalJSON() ([]byte, error) {
	if x.GetValue() == nil {
		return []byte("null"), nil
	}

	return x.GetValue(), nil
}

func (x *RawJson) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	if x.GetValue() == nil {
		return []byte("null"), nil
	}

	return x.GetValue(), nil
}
