package testing

import (
	"encoding/json"
	"time"

	"github.com/s7techlab/cckit/convert"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

// PanicIfError
func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// MustProtoMarshal marshals proto.Message, panics if error
func MustProtoMarshal(pb proto.Message) []byte {
	bb, err := proto.Marshal(pb)
	PanicIfError(err)

	return bb
}

func MustJSONMarshal(val interface{}) []byte {
	bb, err := json.Marshal(val)
	PanicIfError(err)
	return bb
}

// MustProtoUnmarshal unmarshals proto.Message, panics if error
func MustProtoUnmarshal(bb []byte, pm proto.Message) proto.Message {
	p := proto.Clone(pm)
	PanicIfError(proto.Unmarshal(bb, p))
	return p
}

// MustProtoTimestamp, creates proto.Timestamp, panics if error
func MustProtoTimestamp(t time.Time) *timestamp.Timestamp {
	ts, err := ptypes.TimestampProto(t)
	PanicIfError(err)
	return ts
}

func MustConvertFromBytes(bb []byte, target interface{}) interface{} {
	v, err := convert.FromBytes(bb, target)
	PanicIfError(err)
	return v
}
