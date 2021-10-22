package service

import (
	"context"

	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
	"github.com/pkg/errors"
)

type ChaincodeService struct {
	sdk FabricAPI
}

// gateway/chaincode.go needds access to grpc stream
type (
	//Chaincode             = ChaincodeServer
	ChaincodeEventsServer = chaincodeEventsServer
)

type FabricAPI interface {
	// Channel returns channel instance by channel name
	Channel(name string) ChannelAPI
	DeliverClient(mspId string, identity msp.SigningIdentity) (DeliverClientAPI, error)
	CurrentIdentity() msp.SigningIdentity
}

/* channel api interfaces section*/
type ChannelAPI interface {
	// Chaincode returns chaincode instance by chaincode name
	Chaincode(ctx context.Context, name string) (ChaincodeAPI, error)
}

type ChaincodeAPI interface {
	// Invoke returns invoke builder for presented chaincode function
	Invoke(fn string) ChaincodeAPIInvokeBuilder
	// Query returns query builder for presented function and arguments
	Query(fn string, args ...string) ChaincodeAPIQueryBuilder
}

// ChaincodeAPIQueryBuilder describe possibilities how to get query results
type ChaincodeAPIQueryBuilder interface {
	// WithIdentity allows to invoke chaincode from custom identity
	WithIdentity(identity msp.SigningIdentity) ChaincodeAPIQueryBuilder
	// Transient allows to pass arguments to transient map
	Transient(args map[string][]byte) ChaincodeAPIQueryBuilder
	// AsProposalResponse allows to get raw peer response
	AsProposalResponse(ctx context.Context) (*peer.ProposalResponse, error)
}

type ChaincodeAPIInvokeBuilder interface {
	// WithIdentity allows to invoke chaincode from custom identity
	WithIdentity(identity msp.SigningIdentity) ChaincodeAPIInvokeBuilder
	// Transient allows to pass arguments to transient map
	Transient(args map[string][]byte) ChaincodeAPIInvokeBuilder
	// ArgBytes set slice of bytes as argument
	ArgBytes([][]byte) ChaincodeAPIInvokeBuilder
	// Do makes invoke with built arguments
	Do(ctx context.Context) (res *peer.Response, chaincodeTx string, err error)
}

/* deliver client interfaces section */

type DeliverClientAPI interface {
	// SubscribeCC allows to subscribe on chaincode events using name of channel, chaincode and block offset
	SubscribeCC(
		ctx context.Context,
		channelName string,
		ccName string,
		eventCCSeekOption ...func() (*orderer.SeekPosition, *orderer.SeekPosition),
	) (EventCCSubscriptionAPI, error)
}
type EventCCSubscriptionAPI interface {
	// Events initiates internal GRPC stream and returns channel on chaincode events
	Events() chan *peer.ChaincodeEvent
}

func New(sdk FabricAPI) *ChaincodeService {
	return &ChaincodeService{sdk: sdk}
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

	ccApi, err := cs.sdk.
		Channel(in.Channel).
		Chaincode(ctx, in.Chaincode)
	if err != nil {
		return nil, err
	}

	response, _, err := ccApi.Invoke(string(in.Args[0])).
		WithIdentity(signer).
		ArgBytes(in.Args[1:]).
		Transient(in.Transient).
		Do(ctx)

	if err != nil {
		return nil, err
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

	ccApi, err := cs.sdk.Channel(in.Channel).Chaincode(ctx, in.Chaincode)
	if err != nil {
		return nil, err
	}

	resp, err := ccApi.Query(argSs[0], argSs[1:]...).
		WithIdentity(signer).
		Transient(in.Transient).
		AsProposalResponse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, `failed to query chaincode`)
	}

	return resp, nil
}

func (cs *ChaincodeService) Events(in *ChaincodeLocator, stream Chaincode_EventsServer) error {
	deliver, err := cs.sdk.DeliverClient(cs.sdk.CurrentIdentity().GetMSPIdentifier(), cs.sdk.CurrentIdentity())
	if err != nil {
		return err
	}

	events, err := deliver.SubscribeCC(stream.Context(), in.Channel, in.Chaincode)
	if err != nil {
		return err
	}

	for {
		e, ok := <-events.Events()
		if !ok {
			return nil
		}
		if err = stream.Send(e); err != nil {
			return err
		}
	}
}
