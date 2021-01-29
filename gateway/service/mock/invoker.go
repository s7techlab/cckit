package mock

import (
	"context"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/s7techlab/cckit/gateway/service"
	"github.com/s7techlab/cckit/testing"
)

var (
	InvokeErrorResponse = &peer.Response{Status: shim.ERROR, Message: `invoke failed`}
)

func FailInvokeChaincode(chaincodes ...string) ChaincodeInvoker {
	return func(ctx context.Context, mockstub *testing.MockStub, in *service.ChaincodeExec) *peer.Response {
		if in.Type == service.InvocationType_INVOKE {
			for _, c := range chaincodes {
				if in.Input.Chaincode == c {
					return InvokeErrorResponse
				}
			}
		}
		return DefaultInvoker(ctx, mockstub, in)
	}
}

func FailChaincode(chaincodes ...string) ChaincodeInvoker {
	return func(ctx context.Context, mockstub *testing.MockStub, in *service.ChaincodeExec) *peer.Response {
		for _, c := range chaincodes {
			if in.Input.Chaincode == c {
				return InvokeErrorResponse
			}
		}

		return DefaultInvoker(ctx, mockstub, in)
	}
}
