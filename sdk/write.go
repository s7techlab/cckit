package sdk

import (
	"context"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
)

type ChaincodeInvoker interface {
	Invoke(
		ctx context.Context,
		chanName string,
		ccName string,
		args [][]byte,
		identity msp.SigningIdentity,
		transient map[string][]byte,
		txWaiterType string,
	) (res *peer.Response, chaincodeTx string, err error)
}
