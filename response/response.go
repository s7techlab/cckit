package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Success returns shim.Error
func Error(err interface{}) peer.Response {
	return shim.Error(fmt.Sprintf("%s", err))
}

// Success returns shim.Success with serialized json if necessary
func Success(data interface{}) peer.Response {
	switch data.(type) {
	case string:
		return shim.Success([]byte(data.(string)))
	case []byte:
		return shim.Success(data.([]byte))
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return shim.Success(nil)
		} else {
			return shim.Success(b)
		}
	}
}

// Create peer.Response (Success or Error) depending on value of err
// if err is (bool) false or is error interface - returns shim.Error
func Create(data interface{}, err interface{}) peer.Response {
	var errObj error = nil

	switch err.(type) {

	case nil:
		errObj = nil
	case bool:
		if !err.(bool) {
			errObj = errors.New(`boolean error: false`)
		}
	case string:
		if err.(string) != `` {
			errObj = errors.New(err.(string))
		}
	case error:
		errObj = err.(error)
	default:
		panic(fmt.Sprintf(`unknowm error type %s`, err))

	}

	if errObj != nil {
		return Error(errObj)
	} else {
		return Success(data)
	}
}

type Transformer struct {
	data interface{}
	err  error
}

func (t Transformer) With(transfomer func(interface{}) interface{}) peer.Response {
	return Create(transfomer(t.data), t.err)
}

func Transform(data interface{}, err error) *Transformer {
	return &Transformer{data, err}
}
