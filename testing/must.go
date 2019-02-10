package testing

import (
	"github.com/golang/protobuf/proto"
)

func MustProtoMarshal(pb proto.Message) []byte {
	bb, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	return bb
}
