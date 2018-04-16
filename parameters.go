package cckit

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/vitiko/cckit/response"
	"github.com/vitiko/cckit/convert"
)

type (

	FromBytes interface {
		FromBytes([]byte) (result interface{}, err error)
	}


	Parameters []Parameter

	Parameter struct {
		Name   string
		Type   interface{}
		ArgPos int
	}

	parameterBag map[string]MiddlewareFunc
)

func (p Parameter) getArgfromStub(stub shim.ChaincodeStubInterface) (arg interface{}, err error) {
	args := stub.GetArgs()

	if p.ArgPos > len(args) {
		return nil, errors.New(`Arg pos out of range`)
	}
	return convert.FromBytes(args[p.ArgPos+1],p.Type) //first arg is function name
}

func ParameterBag() parameterBag {
	return parameterBag{}
}
func (pbag parameterBag) Add(name string, paramType interface{}) parameterBag {
	pbag[name] = Param(name, paramType)
	return pbag
}


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
