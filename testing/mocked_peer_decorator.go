package testing

import (
	"context"
	"errors"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"

	"github.com/s7techlab/cckit/gateway"
)

var (
	ErrPeerInvoke                   = errors.New(`invoke failed`)
	ErrPeerQuery                    = errors.New(`query failed`)
	ChaincodeSimulationErrorReponse = &peer.Response{
		Status: shim.ERROR, Message: `chaincode simulation produced error response`}
)

type (
	MockedPeerDecorator struct {
		Peer gateway.Peer

		InvokeMutator InvokeMutator
		QueryMutator  QueryMutator
	}

	InvokeMutator func(
		peer gateway.Peer,
		ctx context.Context,
		channel string,
		chaincode string,
		args [][]byte,
		identity msp.SigningIdentity,
		transArgs map[string][]byte,
		txWaiterType string) (res *peer.Response, chaincodeTx string, err error)

	QueryMutator func(
		peer gateway.Peer,
		ctx context.Context,
		channel string,
		chaincode string,
		args [][]byte,
		identity msp.SigningIdentity,
		transArgs map[string][]byte) (*peer.Response, error)
)

func NewPeerDecorator(peer gateway.Peer) *MockedPeerDecorator {
	return &MockedPeerDecorator{
		Peer: peer,
	}
}

func (mpd *MockedPeerDecorator) Invoke(
	ctx context.Context,
	channel string,
	chaincode string,
	args [][]byte,
	identity msp.SigningIdentity,
	transArgs map[string][]byte,
	txWaiterType string) (res *peer.Response, chaincodeTx string, err error) {

	if mpd.InvokeMutator != nil {
		return mpd.InvokeMutator(mpd.Peer, ctx, channel, chaincode, args, identity, transArgs, txWaiterType)
	}

	return mpd.Peer.Invoke(ctx, channel, chaincode, args, identity, transArgs, txWaiterType)
}

func (mpd *MockedPeerDecorator) Query(
	ctx context.Context,
	channel string,
	chaincode string,
	args [][]byte,
	identity msp.SigningIdentity,
	transArgs map[string][]byte) (*peer.Response, error) {

	if mpd.QueryMutator != nil {
		return mpd.QueryMutator(mpd.Peer, ctx, channel, chaincode, args, identity, transArgs)
	}

	return mpd.Peer.Query(ctx, channel, chaincode, args, identity, transArgs)

}

func (mpd *MockedPeerDecorator) Events(
	ctx context.Context,
	channel string,
	chaincode string,
	identity msp.SigningIdentity,
	blockRange ...int64,
) (chan *peer.ChaincodeEvent, error) {
	return mpd.Peer.Events(ctx, channel, chaincode, identity, blockRange...)
}

func (mpd *MockedPeerDecorator) FailChaincode(chaincodes ...string) {
	mpd.FailInvoke(chaincodes...)
	mpd.FailQuery(chaincodes...)
}

func (mpd *MockedPeerDecorator) FailInvoke(chaincodes ...string) {

	mpd.InvokeMutator = func(
		peer gateway.Peer,
		ctx context.Context,
		channel string,
		chaincode string,
		args [][]byte,
		identity msp.SigningIdentity,
		transArgs map[string][]byte,
		txWaiterType string) (res *peer.Response, chaincodeTx string, err error) {

		for _, c := range chaincodes {
			if chaincode == c {
				return nil, ``, ErrPeerInvoke
			}
		}

		return peer.Invoke(ctx, channel, chaincode, args, identity, transArgs, txWaiterType)
	}
}

func (mpd *MockedPeerDecorator) FailQuery(chaincodes ...string) {
	mpd.QueryMutator = func(
		peer gateway.Peer,
		ctx context.Context,
		channel string,
		chaincode string,
		args [][]byte,
		identity msp.SigningIdentity,
		transArgs map[string][]byte) (res *peer.Response, err error) {

		for _, c := range chaincodes {
			if chaincode == c {
				return nil, ErrPeerQuery
			}
		}

		return peer.Query(ctx, channel, chaincode, args, identity, transArgs)
	}
}
