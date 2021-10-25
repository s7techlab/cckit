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
	//Chaincode             = ChaincodeServer
	ChaincodeEventsServer = chaincodeEventsServer
)

type Peer interface {
	ChannelChaincode(ctx context.Context, chanName string, ccName string) (Chaincode, error)
	Events(
		ctx context.Context,
		channelName string,
		ccName string,
		eventCCSeekOption ...func() (*orderer.SeekPosition, *orderer.SeekPosition),
	) (chan *peer.ChaincodeEvent, error)
	CurrentIdentity() msp.SigningIdentity
}

type Chaincode interface {
	// Invoke returns invoke builder for presented chaincode function
	Invoke(fn string) ChaincodeInvokeBuilder
	// Query returns query builder for presented function and arguments
	Query(fn string, args ...string) ChaincodeQueryBuilder
}

// ChaincodeQueryBuilder describe possibilities how to get query results
type ChaincodeQueryBuilder interface {
	// WithIdentity allows to invoke chaincode from custom identity
	WithIdentity(identity msp.SigningIdentity) ChaincodeQueryBuilder
	// Transient allows to pass arguments to transient map
	Transient(args map[string][]byte) ChaincodeQueryBuilder
	// AsProposalResponse allows to get raw peer response
	AsProposalResponse(ctx context.Context) (*peer.ProposalResponse, error)
}

type ChaincodeInvokeBuilder interface {
	// WithIdentity allows to invoke chaincode from custom identity
	WithIdentity(identity msp.SigningIdentity) ChaincodeInvokeBuilder
	// Transient allows to pass arguments to transient map
	Transient(args map[string][]byte) ChaincodeInvokeBuilder
	// ArgBytes set slice of bytes as argument
	ArgBytes([][]byte) ChaincodeInvokeBuilder
	// Do makes invoke with built arguments
	Do(ctx context.Context) (res *peer.Response, chaincodeTx string, err error)
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

	ccAPI, err := cs.peerSDK.ChannelChaincode(ctx, in.Channel, in.Chaincode)
	if err != nil {
		return nil, err
	}

	response, _, err := ccAPI.Invoke(string(in.Args[0])).
		WithIdentity(signer).
		ArgBytes(in.Args[1:]).
		Transient(in.Transient).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	// todo: add to hlf-peerSDK-go method returning ProposalResponse
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

	ccAPI, err := cs.peerSDK.ChannelChaincode(ctx, in.Channel, in.Chaincode)
	if err != nil {
		return nil, err
	}

	resp, err := ccAPI.Query(argSs[0], argSs[1:]...).
		WithIdentity(signer).
		Transient(in.Transient).
		AsProposalResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, `failed to query chaincode`)
	}

	return resp, nil
}

func (cs *ChaincodeService) Events(in *ChaincodeLocator, stream Chaincode_EventsServer) error {
	events, err := cs.peerSDK.Events(stream.Context(), in.Channel, in.Chaincode)
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
