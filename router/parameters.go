package router

import (
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/response"
)

type (
	// Parameters list of chain code function parameters
	Parameters []Parameter

	// Parameter of chain code function
	Parameter struct {
		Name   string
		Type   interface{}
		ArgPos int
	}

	// MiddlewareFuncMap named list of middleware functions
	MiddlewareFuncMap map[string]MiddlewareFunc
)

func (p Parameter) getArgfromStub(stub shim.ChaincodeStubInterface) (arg interface{}, err error) {
	args := stub.GetArgs()

	if p.ArgPos > len(args) {
		return nil, errors.New(`Arg pos out of range`)
	}
	return convert.FromBytes(args[p.ArgPos+1], p.Type) //first arg is function name
}

// ParameterBag builder for named middleware list
func ParameterBag() MiddlewareFuncMap {
	return MiddlewareFuncMap{}
}

// Add middleware function
func (pbag MiddlewareFuncMap) Add(name string, paramType interface{}) MiddlewareFuncMap {
	pbag[name] = Param(name, paramType)
	return pbag
}

// Param create middleware function for transforming stub arg to context arg
func Param(name string, paramType interface{}, argPoss ...int) MiddlewareFunc {
	var argPos int
	if len(argPoss) == 0 {
		argPos = 0
	} else {
		argPos = argPoss[0]
	}

	parameter := Parameter{name, paramType, argPos}

	return func(next HandlerFunc, pos ...int) HandlerFunc {
		return func(context Context) peer.Response {
			arg, err := parameter.getArgfromStub(context.Stub())
			if err != nil {
				return response.Error(err)
			}
			context.SetArg(name, arg)
			return next(context)
		}
	}
}

//if ph.Parameters.Length() != len(args) {
//return nil, ErrArgsNumMismatch
//}
