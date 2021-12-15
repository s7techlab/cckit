package gateway

import (
	"context"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
)

type PeerReader interface {
	ChaincodeQuerier
	EventDelivery
}

type Peer interface {
	PeerReader
	ChaincodeInvoker
}

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

type ChaincodeQuerier interface {
	Query(
		ctx context.Context,
		chanName string,
		ccName string,
		args [][]byte,
		identity msp.SigningIdentity,
		transient map[string][]byte,
	) (*peer.Response, error)
}

type EventDelivery interface {
	Events(
		ctx context.Context,
		channelName string,
		ccName string,
		identity msp.SigningIdentity,
		blockRange ...int64,
	) (chan *peer.ChaincodeEvent, error)
}
