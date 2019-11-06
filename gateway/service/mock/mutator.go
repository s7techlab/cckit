package mock

import (
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/gateway/service"
)

func ResponseError(in *service.ChaincodeExec, r peer.Response) (peer.Response, error) {
	return peer.Response{
		Status: shim.ERROR,
	}, errors.New(`peer is down`)
}

func ResponseInvokeError(in *service.ChaincodeExec, r peer.Response) (peer.Response, error) {
	if in.Type == service.InvocationType_INVOKE {
		return peer.Response{
			Status: shim.ERROR,
		}, errors.New(`peer is down`)
	}

	return r, nil
}
