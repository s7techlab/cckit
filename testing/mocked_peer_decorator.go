package testing

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"

	"github.com/s7techlab/cckit/sdk"
)

var (
	ErrPeerInvoke                   = errors.New(`invoke failed`)
	ErrPeerQuery                    = errors.New(`query failed`)
	ChaincodeSimulationErrorReponse = &peer.Response{
		Status: shim.ERROR, Message: `chaincode simulation produced error response`}
)

type (
	MockedPeerDecorator struct {
		SDK sdk.SDK

		InvokeMutator InvokeMutator
		QueryMutator  QueryMutator
	}

	InvokeMutator func(
		sdk sdk.SDK,
		ctx context.Context,
		channel string,
		chaincode string,
		args [][]byte,
		identity msp.SigningIdentity,
		transArgs map[string][]byte,
		txWaiterType string) (res *peer.Response, chaincodeTx string, err error)

	QueryMutator func(
		sdk sdk.SDK,
		ctx context.Context,
		channel string,
		chaincode string,
		args [][]byte,
		identity msp.SigningIdentity,
		transArgs map[string][]byte) (*peer.Response, error)
)

func NewPeerDecorator(sdk sdk.SDK) *MockedPeerDecorator {
	return &MockedPeerDecorator{
		SDK: sdk,
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
		return mpd.InvokeMutator(mpd.SDK, ctx, channel, chaincode, args, identity, transArgs, txWaiterType)
	}

	return mpd.SDK.Invoke(ctx, channel, chaincode, args, identity, transArgs, txWaiterType)
}

func (mpd *MockedPeerDecorator) Query(
	ctx context.Context,
	channel string,
	chaincode string,
	args [][]byte,
	identity msp.SigningIdentity,
	transArgs map[string][]byte) (*peer.Response, error) {

	if mpd.QueryMutator != nil {
		return mpd.QueryMutator(mpd.SDK, ctx, channel, chaincode, args, identity, transArgs)
	}

	return mpd.SDK.Query(ctx, channel, chaincode, args, identity, transArgs)

}

func (mpd *MockedPeerDecorator) Events(
	ctx context.Context,
	channel string,
	chaincode string,
	identity msp.SigningIdentity,
	blockRange ...int64,
) (events chan interface {
	Event() *peer.ChaincodeEvent
	Block() uint64
	TxTimestamp() *timestamp.Timestamp
}, closer func() error, err error) {
	return mpd.SDK.Events(ctx, channel, chaincode, identity, blockRange...)
}

func (mpd *MockedPeerDecorator) FailChaincode(chaincodes ...string) {
	mpd.FailInvoke(chaincodes...)
	mpd.FailQuery(chaincodes...)
}

func (mpd *MockedPeerDecorator) FailInvoke(chaincodes ...string) {

	mpd.InvokeMutator = func(
		sdk sdk.SDK,
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

		return sdk.Invoke(ctx, channel, chaincode, args, identity, transArgs, txWaiterType)
	}
}

func (mpd *MockedPeerDecorator) FailQuery(chaincodes ...string) {
	mpd.QueryMutator = func(
		sdk sdk.SDK,
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

		return sdk.Query(ctx, channel, chaincode, args, identity, transArgs)
	}
}
