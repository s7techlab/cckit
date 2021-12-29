package gateway

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

func NewChaincodeInstanceService(peer Peer, channel, chaincode string) *ChaincodeInstanceService {
	return &ChaincodeInstanceService{
		ChaincodeService: NewChaincodeService(peer),
		Locator: &ChaincodeLocator{
			Channel:   channel,
			Chaincode: chaincode,
		},
	}
}

func (cis *ChaincodeInstanceService) EventService() *ChaincodeInstanceEventService {
	return NewChaincodeInstanceEventService(cis.ChaincodeService.Peer, cis.Locator.Channel, cis.Locator.Chaincode)
}

func NewChaincodeInstanceEventService(delivery EventDelivery, channel, chaincode string) *ChaincodeInstanceEventService {
	return &ChaincodeInstanceEventService{
		EventService: NewChaincodeEventService(delivery),
		Locator: &ChaincodeLocator{
			Channel:   channel,
			Chaincode: chaincode,
		},
	}
}

func (cis *ChaincodeInstanceService) ServiceDef() ServiceDef {
	return ServiceDef{
		Desc:                        &_ChaincodeInstanceService_serviceDesc,
		Service:                     cis,
		HandlerFromEndpointRegister: RegisterChaincodeInstanceServiceHandlerFromEndpoint,
	}
}

func (cis *ChaincodeInstanceService) Exec(ctx context.Context, exec *ChaincodeInstanceExec) (*peer.Response, error) {
	return cis.ChaincodeService.Exec(ctx, &ChaincodeExec{
		Type: exec.Type,
		Input: &ChaincodeInput{
			Chaincode: cis.Locator,
			Args:      exec.Input.Args,
			Transient: exec.Input.Transient,
		},
	})
}

func (cis *ChaincodeInstanceService) Query(ctx context.Context, input *ChaincodeInstanceInput) (*peer.Response, error) {
	return cis.Exec(ctx, &ChaincodeInstanceExec{
		Type:  InvocationType_QUERY,
		Input: input,
	})
}

func (cis *ChaincodeInstanceService) Invoke(ctx context.Context, input *ChaincodeInstanceInput) (*peer.Response, error) {
	return cis.Exec(ctx, &ChaincodeInstanceExec{
		Type:  InvocationType_INVOKE,
		Input: input,
	})
}

func (cis *ChaincodeInstanceService) EventsStream(request *ChaincodeInstanceEventsStreamRequest, stream ChaincodeInstanceEventsService_EventsStreamServer) error {
	return cis.ChaincodeService.EventsStream(&ChaincodeEventsStreamRequest{
		Chaincode: cis.Locator,
		Block:     request.Block,
	}, stream)
}

func (cis *ChaincodeInstanceService) Events(ctx context.Context, request *ChaincodeInstanceEventsRequest) (*ChaincodeEvents, error) {
	return cis.ChaincodeService.Events(ctx, &ChaincodeEventsRequest{
		Chaincode: cis.Locator,
		Block:     request.Block,
		EventName: request.EventName,
	})
}

func (ces *ChaincodeInstanceEventService) ServiceDef() ServiceDef {
	return ServiceDef{
		Desc:                        &_ChaincodeInstanceEventsService_serviceDesc,
		Service:                     ces,
		HandlerFromEndpointRegister: RegisterChaincodeInstanceEventsServiceHandlerFromEndpoint,
	}
}
