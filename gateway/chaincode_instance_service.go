package gateway

import (
	"context"

	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/s7techlab/cckit/sdk"
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

func NewChaincodeInstanceService(sdk sdk.SDK, channel, chaincode string, opts ...Opt) *ChaincodeInstanceService {
	return &ChaincodeInstanceService{
		ChaincodeService: NewChaincodeService(sdk, opts...),
		Locator: &ChaincodeLocator{
			Channel:   channel,
			Chaincode: chaincode,
		},
	}
}

func (cis *ChaincodeInstanceService) EventService() *ChaincodeInstanceEventService {
	return &ChaincodeInstanceEventService{
		EventService: cis.ChaincodeService.EventService,
		Locator: &ChaincodeLocator{
			Channel:   cis.Locator.Channel,
			Chaincode: cis.Locator.Chaincode,
		},
	}
}

func NewChaincodeInstanceEventService(delivery sdk.EventDelivery, channel, chaincode string, opts ...Opt) *ChaincodeInstanceEventService {
	return &ChaincodeInstanceEventService{
		EventService: NewChaincodeEventService(delivery, opts...),
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
		FromBlock: request.FromBlock,
		ToBlock:   request.ToBlock,
	}, stream)
}

func (cis *ChaincodeInstanceService) Events(ctx context.Context, request *ChaincodeInstanceEventsRequest) (*ChaincodeEvents, error) {
	return cis.ChaincodeService.Events(ctx, &ChaincodeEventsRequest{
		Chaincode: cis.Locator,
		FromBlock: request.FromBlock,
		ToBlock:   request.ToBlock,
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

func (ces *ChaincodeInstanceEventService) EventsStream(request *ChaincodeInstanceEventsStreamRequest, stream ChaincodeInstanceEventsService_EventsStreamServer) error {
	return ces.EventService.EventsStream(&ChaincodeEventsStreamRequest{
		Chaincode: ces.Locator,
		FromBlock: request.FromBlock,
		ToBlock:   request.ToBlock,
	}, stream)
}

func (ces *ChaincodeInstanceEventService) Events(ctx context.Context, request *ChaincodeInstanceEventsRequest) (*ChaincodeEvents, error) {
	return ces.EventService.Events(ctx, &ChaincodeEventsRequest{
		Chaincode: ces.Locator,
		FromBlock: request.FromBlock,
		ToBlock:   request.ToBlock,
	})
}
