package service

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric-protos-go/peer"

	cckit_gateway "github.com/s7techlab/cckit/gateway"
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

// gateway/chaincode.go needs access to grpc stream
type (
	Chaincode             = ChaincodeServiceServer
	ChaincodeServer       = ChaincodeServiceServer
	ChaincodeEventsServer = chaincodeServiceEventsServer
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
func (cs *ChaincodeService) ServiceDef() cckit_gateway.ServiceDef {
	return cckit_gateway.ServiceDef{
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

func (cs *ChaincodeService) Events(in *ChaincodeEventsRequest, stream ChaincodeService_EventsServer) error {
	return cs.EventService.Events(in, stream)
}

// ServiceDef returns service definition
func (ce *ChaincodeEventService) ServiceDef() cckit_gateway.ServiceDef {
	return cckit_gateway.ServiceDef{
		Desc:                        &_ChaincodeEventsService_serviceDesc,
		Service:                     ce,
		HandlerFromEndpointRegister: RegisterChaincodeEventsServiceHandlerFromEndpoint,
	}
}

func (ce *ChaincodeEventService) Events(in *ChaincodeEventsRequest, stream ChaincodeService_EventsServer) error {
	signer, _ := SignerFromContext(stream.Context())

	var blockRange []int64
	if in.Block != nil {
		blockRange = []int64{in.Block.From, in.Block.To}
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
		if err = stream.Send(e); err != nil {
			return err
		}
	}
}
