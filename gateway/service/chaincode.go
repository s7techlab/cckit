package service

import (
	"context"
	fmt "fmt"

	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
)

type ChaincodeService struct {
	peerSDK Peer
}

// gateway/chaincode.go needds access to grpc stream
type (
	Chaincode             = ChaincodeServer
	ChaincodeEventsServer = chaincodeEventsServer
)

type Peer interface {
	Invoke(
		ctx context.Context,
		chanName string,
		ccName string,
		args [][]byte,
		identity msp.SigningIdentity,
		transient map[string][]byte,
	) (res *peer.Response, chaincodeTx string, err error)

	Query(
		ctx context.Context,
		chanName string,
		ccName string,
		args [][]byte,
		identity msp.SigningIdentity,
		transient map[string][]byte,
	) (*peer.ProposalResponse, error)

	Events(
		ctx context.Context,
		channelName string,
		ccName string,
		identity msp.SigningIdentity,
		eventCCSeekOption ...func() (*orderer.SeekPosition, *orderer.SeekPosition),
	) (chan *peer.ChaincodeEvent, error)
}

func New(peerSDK Peer) *ChaincodeService {
	return &ChaincodeService{peerSDK: peerSDK}
}

func (cs *ChaincodeService) Exec(ctx context.Context, in *ChaincodeExec) (*peer.ProposalResponse, error) {
	switch in.Type {
	case InvocationType_QUERY:
		return cs.Query(ctx, in.Input)
	case InvocationType_INVOKE:
		return cs.Invoke(ctx, in.Input)
	default:
		return nil, ErrUnknownInvocationType
	}
}

func (cs *ChaincodeService) Invoke(ctx context.Context, in *ChaincodeInput) (*peer.ProposalResponse, error) {
	// underlying hlf-sdk can handle 'nil' identity cases and set default if call identity wasn't provided
	// if smth goes wrong we'll see it on step below
	signer, _ := SignerFromContext(ctx)

	response, _, err := cs.peerSDK.Invoke(
		ctx,
		in.Channel,
		in.Chaincode,
		in.Args,
		signer,
		in.Transient,
	)
	if err != nil {
		return nil, fmt.Errorf("invoke chaincode: %w", err)
	}

	// todo: add to hlf-sdk-go method returning ProposalResponse
	proposalResponse := &peer.ProposalResponse{
		Response: response,
	}
	return proposalResponse, nil
}

func (cs *ChaincodeService) Query(ctx context.Context, in *ChaincodeInput) (*peer.ProposalResponse, error) {
	signer, _ := SignerFromContext(ctx)

	resp, err := cs.peerSDK.Query(
		ctx,
		in.Channel,
		in.Chaincode,
		in.Args,
		signer,
		in.Transient,
	)
	if err != nil {
		return nil, fmt.Errorf("query chaincode: %w", err)
	}

	return resp, nil
}

func (cs *ChaincodeService) Events(in *ChaincodeLocator, stream Chaincode_EventsServer) error {
	signer, _ := SignerFromContext(stream.Context())

	events, err := cs.peerSDK.Events(
		stream.Context(),
		in.Channel,
		in.Chaincode,
		signer,
	)
	if err != nil {
		return err
	}

	for {
		e, ok := <-events
		if !ok {
			return nil
		}
		if err = stream.Send(e); err != nil {
			return err
		}
	}
}
