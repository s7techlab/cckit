package service

import (
	"context"

	"github.com/hyperledger/fabric-protos-go/peer"
)

type (
	ChaincodeInstanceService struct {
		ChaincodeService *ChaincodeService
		Locator          *ChaincodeLocator
	}

	ChaincodeInstanceEventService struct {
		EventService *ChaincodeEventService
		Locator      *ChaincodeLocator
	}
)

func (cis *ChaincodeInstanceService) Exec(ctx context.Context, exec *ChaincodeInstanceExec) (*peer.ProposalResponse, error) {
	return cis.ChaincodeService.Exec(ctx, &ChaincodeExec{
		Type: exec.Type,
		Input: &ChaincodeInput{
			Chaincode: cis.Locator,
			Args:      exec.Input.Args,
			Transient: exec.Input.Transient,
		},
	})
}

func (cis *ChaincodeInstanceService) Query(ctx context.Context, input *ChaincodeInstanceInput) (*peer.ProposalResponse, error) {
	return cis.Exec(ctx, &ChaincodeInstanceExec{
		Type:  InvocationType_QUERY,
		Input: input,
	})
}

func (cis *ChaincodeInstanceService) Invoke(ctx context.Context, input *ChaincodeInstanceInput) (*peer.ProposalResponse, error) {
	return cis.Exec(ctx, &ChaincodeInstanceExec{
		Type:  InvocationType_INVOKE,
		Input: input,
	})
}

func (cis ChaincodeInstanceService) Events(request *ChaincodeInstanceEventsRequest, stream ChaincodeInstanceService_EventsServer) error {
	return cis.ChaincodeService.Events(&ChaincodeEventsRequest{
		Chaincode: cis.Locator,
		Block:     request.Block,
	}, stream)
}

func (ces ChaincodeInstanceEventService) Events(request *ChaincodeInstanceEventsRequest, stream ChaincodeInstanceService_EventsServer) error {
	return ces.EventService.Events(&ChaincodeEventsRequest{
		Chaincode: ces.Locator,
		Block:     request.Block,
	}, stream)
}