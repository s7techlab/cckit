package gateway

import (
	"context"

	"github.com/s7techlab/cckit/sdk"
)

type (
	ChaincodeEventService struct {
		EventDelivery sdk.EventDelivery
	}

	// ChaincodeEventsServer  gateway/chaincode.go needs access to grpc stream
	ChaincodeEventsServer = chaincodeServiceEventsStreamServer
)

var _ ChaincodeEventsServiceServer = &ChaincodeEventService{}

// ServiceDef returns service definition
func (ce *ChaincodeEventService) ServiceDef() ServiceDef {
	return ServiceDef{
		Desc:                        &_ChaincodeEventsService_serviceDesc,
		Service:                     ce,
		HandlerFromEndpointRegister: RegisterChaincodeEventsServiceHandlerFromEndpoint,
	}
}

func (ce *ChaincodeEventService) Events(ctx context.Context, req *ChaincodeEventsRequest) (*ChaincodeEvents, error) {
	return NewChaincodeInstanceEventService(ce.EventDelivery, req.Locator.Channel, req.Locator.Chaincode).
		Events(ctx, &ChaincodeInstanceEventsRequest{
			FromBlock: req.FromBlock,
			ToBlock:   req.ToBlock,
			EventName: req.EventName,
			Limit:     req.Limit,
		})
}

func (ce *ChaincodeEventService) EventsStream(req *ChaincodeEventsStreamRequest, stream ChaincodeEventsService_EventsStreamServer) error {
	return NewChaincodeInstanceEventService(ce.EventDelivery, req.Locator.Channel, req.Locator.Chaincode).
		EventsStream(&ChaincodeInstanceEventsStreamRequest{
			FromBlock: req.FromBlock,
			ToBlock:   req.ToBlock,
			EventName: req.EventName,
		}, stream)
}
