package sdk

import (
	"context"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
)

type (
	Reader interface {
		Evaluator
		EventDelivery
	}

	Evaluator interface {
		Query(
			ctx context.Context,
			chanName string,
			ccName string,
			args [][]byte,
			identity msp.SigningIdentity,
			transient map[string][]byte,
		) (*peer.Response, error)
	}

	EventDelivery interface {

		// Events returns peer chaincode events stream
		// blockRange:
		// * [] -> api.SeekNewest()},
		// *  [0] - from oldest block to maxBlock
		// * [{blockFrom}] - from specified block to maxBlock
		// * [-{blockFrom}] - specified blocks back from channel height
		// * [0,0] from oldest block to current channel height
		// * [-{blockFrom},0] - specified blocks back from channel height to current channel height
		// * [-{blockFrom}, -{blockTo}} -{blockFrom} blocks back from channel height to -{blockTo} block from channel height
		// * [-{blockFrom}, {blockTo}} -{blockFrom} blocks back from channel height to block {blockTo}
		Events(
			ctx context.Context,
			channelName string,
			ccName string,
			identity msp.SigningIdentity,
			blockRange ...int64,
		) (events chan interface {
			Event() *peer.ChaincodeEvent
			Block() uint64
			TxTimestamp() *timestamp.Timestamp
		}, closer func() error, err error)
	}
)
