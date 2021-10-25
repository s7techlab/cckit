package service

import (
	"context"

	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
	"github.com/pkg/errors"
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
		identity msp.SigningIdentity,
		fnName string,
		args [][]byte,
		transient map[string][]byte,
	) (res *peer.Response, chaincodeTx string, err error)

	Query(
		ctx context.Context,
		chanName string,
		ccName string,
		identity msp.SigningIdentity,
		fnName string,
		args []string,
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
	signer, err := SignerFromContext(ctx)
	if err != nil {
		return nil, err
	}

	response, _, err := cs.peerSDK.Invoke(
		ctx,
		in.Channel,
		in.Chaincode,
		signer,
		string(in.Args[0]),
		in.Args[1:],
		in.Transient,
	)
	if err != nil {
		return nil, errors.Wrap(err, `failed to invoke chaincode`)
	}

	// todo: add to hlf-sdk-go method returning ProposalResponse
	proposalResponse := &peer.ProposalResponse{
		Response: response,
	}
	return proposalResponse, nil
}

func (cs *ChaincodeService) Query(ctx context.Context, in *ChaincodeInput) (*peer.ProposalResponse, error) {
	argSs := make([]string, 0)
	for _, arg := range in.Args {
		argSs = append(argSs, string(arg))
	}

	signer, err := SignerFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cs.peerSDK.Query(
		ctx,
		in.Channel,
		in.Chaincode,
		signer,
		argSs[0],
		argSs[1:],
		in.Transient,
	)
	if err != nil {
		return nil, errors.Wrap(err, `failed to query chaincode`)
	}

	return resp, nil
}

func (cs *ChaincodeService) Events(in *ChaincodeLocator, stream Chaincode_EventsServer) error {
	signer, err := SignerFromContext(stream.Context())
	if err != nil {
		return err
	}

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
