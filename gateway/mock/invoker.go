package mock

import (
	"context"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/s7techlab/cckit/gateway"
	"github.com/s7techlab/cckit/testing"
)

var (
	InvokeErrorResponse = &peer.Response{Status: shim.ERROR, Message: `invoke failed`}
)

func FailInvokeChaincode(chaincodes ...string) ChaincodeInvoker {
	return func(ctx context.Context, mockstub *testing.MockStub, in *gateway.ChaincodeExec) *peer.Response {
		if in.Type == gateway.InvocationType_INVOKE {
			for _, c := range chaincodes {
				if in.Input.Chaincode.Chaincode == c {
					return InvokeErrorResponse
				}
			}
		}
		return DefaultInvoker(ctx, mockstub, in)
	}
}

func FailChaincode(chaincodes ...string) ChaincodeInvoker {
	return func(ctx context.Context, mockstub *testing.MockStub, in *gateway.ChaincodeExec) *peer.Response {
		for _, c := range chaincodes {
			if in.Input.Chaincode.Chaincode == c {
				return InvokeErrorResponse
			}
		}

		return DefaultInvoker(ctx, mockstub, in)
	}
}
