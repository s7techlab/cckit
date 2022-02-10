package gateway

import (
	"context"
	"time"

	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/sdk"
)

const EventListStreamTimeout = 500 * time.Millisecond

var _ ChaincodeServiceServer = &ChaincodeService{}

type ChaincodeService struct {
	SDK          sdk.SDK
	EventService *ChaincodeEventService
}

// Deprecated: use NewChaincodeService instead
func New(sdk sdk.SDK) *ChaincodeService {
	return NewChaincodeService(sdk)
}

func NewChaincodeService(sdk sdk.SDK) *ChaincodeService {
	ccService := &ChaincodeService{
		SDK:          sdk,
		EventService: NewChaincodeEventService(sdk),
	}

	return ccService
}

func NewChaincodeEventService(eventDelivery sdk.EventDelivery) *ChaincodeEventService {
	eventService := &ChaincodeEventService{
		EventDelivery: eventDelivery,
	}

	return eventService
}

// InstanceService returns ChaincodeInstanceService for current Peer interface and provided channel and chaincode name
func (cs *ChaincodeService) InstanceService(locator *ChaincodeLocator, opts ...OptFunc) *ChaincodeInstanceService {
	return NewChaincodeInstanceService(cs.SDK, locator, opts...)
}

// ServiceDef returns service definition
func (cs *ChaincodeService) ServiceDef() ServiceDef {
	return ServiceDef{
		Desc:                        &_ChaincodeService_serviceDesc,
		Service:                     cs,
		HandlerFromEndpointRegister: RegisterChaincodeServiceHandlerFromEndpoint,
	}
}

func (cs *ChaincodeService) Exec(ctx context.Context, req *ChaincodeExecRequest) (*peer.Response, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	return cs.InstanceService(req.Locator).Exec(ctx, &ChaincodeInstanceExecRequest{
		Type: req.Type, Input: req.Input})
}

func (cs *ChaincodeService) Invoke(ctx context.Context, req *ChaincodeInvokeRequest) (*peer.Response, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	return cs.InstanceService(req.Locator).Invoke(ctx, &ChaincodeInstanceInvokeRequest{Input: req.Input})
}

func (cs *ChaincodeService) Query(ctx context.Context, req *ChaincodeQueryRequest) (*peer.Response, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	return cs.InstanceService(req.Locator).Query(ctx, &ChaincodeInstanceQueryRequest{Input: req.Input})
}

func (cs *ChaincodeService) EventsStream(req *ChaincodeEventsStreamRequest, stream ChaincodeService_EventsStreamServer) error {
	if err := router.ValidateRequest(req); err != nil {
		return err
	}
	return cs.EventService.EventsStream(req, stream)

}

func (cs *ChaincodeService) Events(ctx context.Context, req *ChaincodeEventsRequest) (*ChaincodeEvents, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}
	return cs.EventService.Events(ctx, req)
}
