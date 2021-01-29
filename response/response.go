package response

import (
	"errors"
	"fmt"

	"github.com/s7techlab/cckit/convert"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// Error returns shim.Error
func Error(err interface{}) peer.Response {
	return shim.Error(fmt.Sprintf("%s", err))
}

// Success returns shim.Success with serialized json if necessary
func Success(data interface{}) peer.Response {
	bb, err := convert.ToBytes(data)
	if err != nil {
		return shim.Success(nil)
	}
	return shim.Success(bb)
}

// Create returns peer.Response (Success or Error) depending on value of err
// if err is (bool) false or is error interface - returns shim.Error
func Create(data interface{}, err interface{}) peer.Response {
	var errObj error

	switch e := err.(type) {
	case nil:
		errObj = nil
	case bool:
		if !e {
			errObj = errors.New(`boolean error: false`)
		}
	case string:
		if e != `` {
			errObj = errors.New(e)
		}
	case error:
		errObj = e
	default:
		panic(fmt.Sprintf(`unknowm error type %s`, err))

	}

	if errObj != nil {
		return Error(errObj)
	}
	return Success(data)
}

// Transformer type transforms data
type Transformer struct {
	data interface{}
	err  error
}

// With func transformer
func (t Transformer) With(transfomer func(interface{}) interface{}) peer.Response {
	return Create(transfomer(t.data), t.err)
}

// Transform creates Transformer struct
func Transform(data interface{}, err error) *Transformer {
	return &Transformer{data, err}
}
