package gateway

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric-protos-go/peer"
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

func (cs *ChaincodeService) EventsStream(in *ChaincodeEventsStreamRequest, stream ChaincodeService_EventsStreamServer) error {
	return cs.EventService.EventsStream(&ChaincodeEventsStreamRequest{}, stream)

}

func (cs *ChaincodeService) Events(ctx context.Context, request *ChaincodeEventsRequest) (*ChaincodeEvents, error) {
	return cs.EventService.Events(ctx, request)
}

// ServiceDef returns service definition
func (ce *ChaincodeEventService) ServiceDef() ServiceDef {
	return ServiceDef{
		Desc:                        &_ChaincodeEventsService_serviceDesc,
		Service:                     ce,
		HandlerFromEndpointRegister: RegisterChaincodeEventsServiceHandlerFromEndpoint,
	}
}

func (ce *ChaincodeEventService) Events(ctx context.Context, in *ChaincodeEventsRequest) (*ChaincodeEvents, error) {
	signer, _ := SignerFromContext(ctx)

	// by default from == 0, so we go stream with seek from oldest blocks
	// by default to == 0, so we go stream with seek until current channel height
	blockRange := []int64{0, 0}
	if in.Block != nil {
		blockRange = []int64{in.Block.From, in.Block.To} //
	}

	eventStream, err := ce.EventDelivery.Events(
		ctx,
		in.Chaincode.Channel,
		in.Chaincode.Chaincode,
		signer,
		blockRange...,
	)

	if err != nil {
		return nil, err
	}

	events := &ChaincodeEvents{
		Chaincode: in.Chaincode,
		Block:     in.Block,
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

func (ce *ChaincodeEventService) EventsStream(in *ChaincodeEventsStreamRequest, stream ChaincodeEventsService_EventsStreamServer) error {
	signer, _ := SignerFromContext(stream.Context())

	var blockRange []int64

	if in.Block != nil {
		blockRange = []int64{}
	}

	events, err := ce.EventDelivery.Events(
		stream.Context(),
		in.Chaincode.Channel,
		in.Chaincode.Chaincode,
		signer,
		blockRange...,
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
