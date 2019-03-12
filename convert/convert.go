// Package convert for transforming  between json serialized  []byte and go structs
package convert

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
)

var (
	// ErrUnableToConvertNilToStruct - nil cannot be converted to struct
	ErrUnableToConvertNilToStruct = errors.New(`unable to convert nil to [struct,array,slice,ptr]`)
	// ErrUnableToConvertValueToStruct - value  cannot be converted to struct
	ErrUnableToConvertValueToStruct = errors.New(`unable to convert value to struct`)
)

const TypeInt = 1
const TypeString = ``
const TypeBool = true

type (
	// FromByter interface supports FromBytes func for converting from slice of bytes to target type
	FromByter interface {
		FromBytes([]byte) (interface{}, error)
	}

	// ToByter interface supports ToBytes func for converting to slice of bytes from source type
	ToByter interface {
		ToBytes() ([]byte, error)
	}
)

// TimestampToTime converts timestamp to time.Time
func TimestampToTime(ts *timestamp.Timestamp) time.Time {
	return time.Unix(ts.GetSeconds(), int64(ts.GetNanos()))
}
