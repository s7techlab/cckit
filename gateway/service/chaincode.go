package service

import (
	"context"
	"log"

	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/s7techlab/hlf-sdk-go/api"
)

type (
	// Chaincode service interface
	Chaincode             = ChaincodeServer
	ChaincodeEventsServer = chaincodeEventsServer
)

// ChaincodeService implementation based of hlf-sdk-go
type ChaincodeService struct {
	sdk api.Core
}

func New(sdk api.Core) *ChaincodeService {
	return &ChaincodeService{sdk: sdk}
}

func (cs *ChaincodeService) Exec(ctx context.Context, in *ChaincodeExec) (*peer.ProposalResponse, error) {
	if in.Type == InvocationType_QUERY {
		return cs.Query(ctx, in.Input)
	} else if in.Type == InvocationType_INVOKE {
		return cs.Invoke(ctx, in.Input)
	} else {
		return nil, ErrUnknownInvocationType
	}
}

func (cs *ChaincodeService) Invoke(ctx context.Context, in *ChaincodeInput) (*peer.ProposalResponse, error) {
	signer, err := SignerFromContext(ctx)
	if err != nil {
		return nil, err
	}

	response, _, err := cs.sdk.
		Channel(in.Channel).
		Chaincode(in.Chaincode).
		Invoke(string(in.Args[0])).
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

	if resp, err := cs.sdk.Channel(in.Channel).Chaincode(in.Chaincode).Query(argSs[0], argSs[1:]...).WithIdentity(signer).Transient(in.Transient).AsProposalResponse(ctx); err != nil {
		return nil, errors.Wrap(err, `failed to query chaincode`)
	} else {
		return resp, nil
	}
}

func (cs *ChaincodeService) Events(in *ChaincodeLocator, stream Chaincode_EventsServer) error {

	deliver, err := cs.sdk.PeerPool().DeliverClient(cs.sdk.CurrentIdentity().GetMSPIdentifier(), cs.sdk.CurrentIdentity())
	if err != nil {
		return err
	}

	events, err := deliver.SubscribeCC(stream.Context(), in.Channel, in.Chaincode)
	if err != nil {
		return err
	}

	for {
		e, ok := <-events.Events()

		log.Println(`event received`, e.EventName)
		if !ok {
			return nil
		}
		if err = stream.Send(e); err != nil {
			return err
		}
	}
}
