package testing

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

func MustProtoMarshal(pb proto.Message) []byte {
	bb, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	return bb
}

func MustProtoTimestamp(t time.Time) *timestamp.Timestamp {
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	return ts
}
