package param

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/router"
)

var (
	ErrProtoExpected = errors.New(`protobuf expected`)
)

// String creates middleware for converting to string chaincode method parameter
func String(name string, argPoss ...int) router.MiddlewareFunc {
	return Param(name, convert.TypeString, argPoss...)
}

func Strings(name string, argPoss ...int) router.MiddlewareFunc {
	return Param(name, []string{}, argPoss...)
}

// Int creates middleware for converting to integer chaincode method parameter
func Int(name string, argPoss ...int) router.MiddlewareFunc {
	return Param(name, convert.TypeInt, argPoss...)
}

// Bool creates middleware for converting to bool chaincode method parameter
func Bool(name string, argPoss ...int) router.MiddlewareFunc {
	return Param(name, convert.TypeBool, argPoss...)
}

// Struct creates middleware for converting to struct chaincode method parameter
func Struct(name string, target interface{}, argPoss ...int) router.MiddlewareFunc {
	return Param(name, target, argPoss...)
}

// Bytes creates middleware for converting to []byte chaincode method parameter
func Bytes(name string, argPoss ...int) router.MiddlewareFunc {
	return Param(name, []byte{}, argPoss...)
}

// Proto creates middleware for converting to protobuf chaincode method parameter
func Proto(name string, target interface{}, argPoss ...int) router.MiddlewareFunc {
	if _, ok := target.(proto.Message); !ok {
		TypeErrorMiddleware(name, ErrProtoExpected)
	}
	return Param(name, target, argPoss...)
}

func TypeErrorMiddleware(name string, err error) router.MiddlewareFunc {
	return func(next router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(c router.Context) (interface{}, error) {
			return nil, fmt.Errorf(`%s: %s`, err, name)
		}
	}
}
