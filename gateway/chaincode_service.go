package gateway

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/s7techlab/cckit/router"
)

type (
	ChaincodeService struct {
		Peer         Peer
		EventService *ChaincodeEventService
	}

	ChaincodeEventService struct {
		EventDelivery EventDelivery
	}
)

// ChaincodeEventsServer  gateway/chaincode.go needs access to grpc stream
type (
	ChaincodeEventsServer = chaincodeServiceEventsStreamServer
)

// Deprecated: use NewChaincodeService instead
func New(peer Peer) *ChaincodeService {
	return NewChaincodeService(peer)
}

func NewChaincodeService(peer Peer) *ChaincodeService {
	return &ChaincodeService{
		Peer:         peer,
		EventService: NewChaincodeEventService(peer),
	}
}

func NewChaincodeEventService(eventDelivery EventDelivery) *ChaincodeEventService {
	return &ChaincodeEventService{EventDelivery: eventDelivery}
}

// InstanceService returns ChaincodeInstanceService for current Peer interface and provided channel and chaincode name
func (cs *ChaincodeService) InstanceService(channel, chaincode string) *ChaincodeInstanceService {
	return NewChaincodeInstanceService(cs.Peer, channel, chaincode)
}

// ServiceDef returns service definition
func (cs *ChaincodeService) ServiceDef() ServiceDef {
	return ServiceDef{
		Desc:                        &_ChaincodeService_serviceDesc,
		Service:                     cs,
		HandlerFromEndpointRegister: RegisterChaincodeServiceHandlerFromEndpoint,
	}
}

func (cs *ChaincodeService) Exec(ctx context.Context, in *ChaincodeExec) (*peer.Response, error) {
	switch in.Type {
	case InvocationType_QUERY:
		return cs.Query(ctx, in.Input)
	case InvocationType_INVOKE:
		return cs.Invoke(ctx, in.Input)
	default:
		return nil, ErrUnknownInvocationType
	}
}

func (cs *ChaincodeService) Invoke(ctx context.Context, in *ChaincodeInput) (*peer.Response, error) {
	// underlying hlf-sdk(or your implementation must handle it) can handle 'nil' identity cases and set default if call identity wasn't provided
	// if smth goes wrong we'll see it on the step below
	signer, _ := SignerFromContext(ctx)

	response, _, err := cs.Peer.Invoke(
		ctx,
		in.Chaincode.Channel,
		in.Chaincode.Chaincode,
		in.Args,
		signer,
		in.Transient,
		TxWaiterFromContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("invoke chaincode: %w", err)
	}

	return response, nil
}

func (cs *ChaincodeService) Query(ctx context.Context, in *ChaincodeInput) (*peer.Response, error) {
	if err := router.ValidateRequest(in); err != nil {
		return nil, err
	}

	signer, _ := SignerFromContext(ctx)

	resp, err := cs.Peer.Query(
		ctx,
		in.Chaincode.Channel,
		in.Chaincode.Chaincode,
		in.Args,
		signer,
		in.Transient,
	)
	if err != nil {
		return nil, fmt.Errorf("query chaincode: %w", err)
	}

	return resp, nil
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

// ServiceDef returns service definition
func (ce *ChaincodeEventService) ServiceDef() ServiceDef {
	return ServiceDef{
		Desc:                        &_ChaincodeEventsService_serviceDesc,
		Service:                     ce,
		HandlerFromEndpointRegister: RegisterChaincodeEventsServiceHandlerFromEndpoint,
	}
}

func BlockRange(from, to *BlockLimit) []int64 {
	var blockRange []int64

	switch {
	case from != nil && to != nil:
		blockRange = []int64{from.Num, to.Num}

	case from != nil:
		blockRange = []int64{from.Num}

	case to != nil:
		blockRange = []int64{0, to.Num}
	}

	return blockRange
}

func (ce *ChaincodeEventService) Events(ctx context.Context, req *ChaincodeEventsRequest) (*ChaincodeEvents, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	signer, _ := SignerFromContext(ctx)

	// for event _list_ default block range is up to current channel height
	if req.ToBlock == nil {
		req.ToBlock = &BlockLimit{Num: 0}
	}

	eventStream, err := ce.EventDelivery.Events(
		ctx,
		req.Chaincode.Channel,
		req.Chaincode.Chaincode,
		signer,
		BlockRange(req.FromBlock, req.ToBlock)...,
	)

	if err != nil {
		return nil, err
	}

	events := &ChaincodeEvents{
		Chaincode: req.Chaincode,
		FromBlock: req.FromBlock,
		ToBlock:   req.ToBlock,
		Items:     []*ChaincodeEvent{},
	}

	for {
		event, hasMore := <-eventStream
		if !hasMore {
			break
		}

		events.Items = append(events.Items, &ChaincodeEvent{
			Event: event,
		})
	}

	return events, nil
}

func (ce *ChaincodeEventService) EventsStream(req *ChaincodeEventsStreamRequest, stream ChaincodeEventsService_EventsStreamServer) error {
	if err := router.ValidateRequest(req); err != nil {
		return err
	}

	signer, _ := SignerFromContext(stream.Context())

	events, err := ce.EventDelivery.Events(
		stream.Context(),
		req.Chaincode.Channel,
		req.Chaincode.Chaincode,
		signer,
		BlockRange(req.FromBlock, req.ToBlock)...,
	)
	if err != nil {
		return err
	}

	for {
		e, ok := <-events
		if !ok {
			return nil
		}
		if err = stream.Send(&ChaincodeEvent{Event: e}); err != nil {
			return err
		}
	}
}
