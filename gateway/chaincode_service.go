package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/sdk"
)

const EventListStreamTimeout = 500 * time.Millisecond

type (
	ChaincodeService struct {
		SDK          sdk.SDK
		EventService *ChaincodeEventService
		Opts         *Opts
	}

	ChaincodeEventService struct {
		EventDelivery sdk.EventDelivery
		Opts          *Opts
	}
)

// ChaincodeEventsServer  gateway/chaincode.go needs access to grpc stream
type (
	ChaincodeEventsServer = chaincodeServiceEventsStreamServer
)

// Deprecated: use NewChaincodeService instead
func New(sdk sdk.SDK, opts ...Opt) *ChaincodeService {
	return NewChaincodeService(sdk, opts...)
}

func NewChaincodeService(sdk sdk.SDK, opts ...Opt) *ChaincodeService {

	ccService := &ChaincodeService{
		SDK:  sdk,
		Opts: &Opts{},
	}

	for _, o := range opts {
		o(ccService.Opts)
	}

	ccService.EventService = NewChaincodeEventService(sdk)
	ccService.EventService.Opts = ccService.Opts
	return ccService
}

func NewChaincodeEventService(eventDelivery sdk.EventDelivery, opts ...Opt) *ChaincodeEventService {
	eventService := &ChaincodeEventService{
		EventDelivery: eventDelivery,
		Opts:          &Opts{},
	}

	for _, o := range opts {
		o(eventService.Opts)
	}

	return eventService
}

// InstanceService returns ChaincodeInstanceService for current Peer interface and provided channel and chaincode name
func (cs *ChaincodeService) InstanceService(channel, chaincode string) *ChaincodeInstanceService {
	return &ChaincodeInstanceService{
		ChaincodeService: cs,
		Locator: &ChaincodeLocator{
			Channel:   channel,
			Chaincode: chaincode,
		},
	}
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

	response, _, err := cs.SDK.Invoke(
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

	resp, err := cs.SDK.Query(
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

	eventStream, closer, err := ce.EventDelivery.Events(
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
		// wait for timeout, return values if events are not streamed
		ticker := time.NewTicker(EventListStreamTimeout)

		select {

		case <-ctx.Done():
			_ = closer()
			return events, nil

		case <-ticker.C:
			_ = closer()
			return events, nil

		case event, hasMore := <-eventStream:
			if !hasMore {
				_ = closer()
				return events, nil
			}

			if !MatchEventName(event.Event().EventName, req.EventName) {
				continue
			}

			ccEvent := &ChaincodeEvent{
				Event:       event.Event(),
				Block:       event.Block(),
				TxTimestamp: event.TxTimestamp(),
			}
			for _, o := range ce.Opts.Event {
				_ = o(ccEvent)
			}

			events.Items = append(events.Items, ccEvent)
		}
	}

}

func (ce *ChaincodeEventService) EventsStream(req *ChaincodeEventsStreamRequest, stream ChaincodeEventsService_EventsStreamServer) error {
	if err := router.ValidateRequest(req); err != nil {
		return err
	}

	signer, _ := SignerFromContext(stream.Context())

	events, closer, err := ce.EventDelivery.Events(
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
		select {

		case <-stream.Context().Done():
			_ = closer()
			return nil

		case event, ok := <-events:
			if !ok {
				return nil
			}

			if !MatchEventName(event.Event().EventName, req.EventName) {
				continue
			}

			ccEvent := &ChaincodeEvent{
				Event:       event.Event(),
				Block:       event.Block(),
				TxTimestamp: event.TxTimestamp(),
			}
			for _, o := range ce.Opts.Event {
				_ = o(ccEvent)
			}

			if err = stream.Send(ccEvent); err != nil {
				return err
			}

		}
	}
}

func MatchEventName(eventName string, expectedNames []string) bool {
	if len(expectedNames) == 0 {
		return true
	}

	for _, expectedName := range expectedNames {
		if eventName == expectedName {
			return true
		}
	}

	return false
}
