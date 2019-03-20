package testing

import (
	"encoding/json"
	"time"

	"github.com/s7techlab/cckit/convert"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

// MustProtoMarshal marshals proto.Message, panics if error
func MustProtoMarshal(pb proto.Message) []byte {
	bb, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	return bb
}

func MustJsonMarshal(val interface{}) []byte {
	bb, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	return bb
}

// MustProtoUnmarshal unmarshals proto.Message, panics if error
func MustProtoUnmarshal(bb []byte, pm proto.Message) proto.Message {
	p := proto.Clone(pm)
	if err := proto.Unmarshal(bb, p); err != nil {
		panic(err)
	}
	return p
}

// MustProtoTimestamp, creates proto.Timestamp, panics if error
func MustProtoTimestamp(t time.Time) *timestamp.Timestamp {
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	return ts
}

func MustConvertFromBytes(bb []byte, target interface{}) interface{} {
	v, err := convert.FromBytes(bb, target)
	if err != nil {
		panic(err)
	}
	return v
}
