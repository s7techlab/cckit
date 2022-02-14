package gateway

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"go.uber.org/zap"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/sdk"
)

type (
	ChaincodeInstanceEventService struct {
		EventDelivery sdk.EventDelivery
		Locator       *ChaincodeLocator
		Opts          *Opts
		Logger        *zap.Logger
	}

	ChaincodeInstanceEvents interface {
		ChaincodeInstanceEventsServiceServer

		EventsChan(ctx context.Context, rr ...*ChaincodeInstanceEventsStreamRequest) (
			_ chan *ChaincodeEvent, closer func() error, _ error)
	}
)

var _ ChaincodeInstanceEventsServiceServer = &ChaincodeInstanceEventService{}

func NewChaincodeInstanceEventService(delivery sdk.EventDelivery, channel, chaincode string, opts ...Opt) *ChaincodeInstanceEventService {
	cis := &ChaincodeInstanceEventService{
		EventDelivery: delivery,
		Locator: &ChaincodeLocator{
			Channel:   channel,
			Chaincode: chaincode,
		},
		Opts:   &Opts{},
		Logger: zap.NewNop(),
	}

	for _, o := range opts {
		o(cis.Opts)
	}

	return cis
}

func (ces *ChaincodeInstanceEventService) ServiceDef() ServiceDef {
	return ServiceDef{
		Desc:                        &_ChaincodeInstanceEventsService_serviceDesc,
		Service:                     ces,
		HandlerFromEndpointRegister: RegisterChaincodeInstanceEventsServiceHandlerFromEndpoint,
	}
}

func (ces *ChaincodeInstanceEventService) EventsStream(
	req *ChaincodeInstanceEventsStreamRequest, stream ChaincodeInstanceEventsService_EventsStreamServer) error {
	if err := router.ValidateRequest(req); err != nil {
		return err
	}

	ctx := stream.Context()
	for _, c := range ces.Opts.Context {
		ctx = c(ctx)
	}

	signer, _ := SignerFromContext(stream.Context())
	events, closer, err := ces.EventDelivery.Events(
		stream.Context(),
		ces.Locator.Channel,
		ces.Locator.Chaincode,
		signer,
		BlockRange(req.FromBlock, req.ToBlock)...,
	)
	if err != nil {
		return err
	}

	for {
		select {

		case <-stream.Context().Done():
			return closer()

		case e, ok := <-events:
			if !ok {
				return nil
			}

			processedEvent, err := ProcessEvent(e, ces.Opts.Event, req.EventName)
			if err != nil {
				ces.Logger.Warn(`event processing`, zap.Error(err))
			}

			if processedEvent != nil {
				if err = stream.Send(processedEvent); err != nil {
					return err
				}
			}
		}
	}
}

func (ces *ChaincodeInstanceEventService) Events(ctx context.Context, req *ChaincodeInstanceEventsRequest) (*ChaincodeEvents, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	// for event _list_ default block range is up to current channel height
	if req.ToBlock == nil {
		req.ToBlock = &BlockLimit{Num: 0}
	}

	for _, c := range ces.Opts.Context {
		ctx = c(ctx)
	}
	signer, _ := SignerFromContext(ctx)
	eventStream, closer, err := ces.EventDelivery.Events(
		ctx,
		ces.Locator.Channel,
		ces.Locator.Chaincode,
		signer,
		BlockRange(req.FromBlock, req.ToBlock)...,
	)

	if err != nil {
		return nil, err
	}

	events := &ChaincodeEvents{
		Locator:   ces.Locator,
		FromBlock: req.FromBlock,
		ToBlock:   req.ToBlock,
		Items:     []*ChaincodeEvent{},
	}

	// timeout for receiving events
	ticker := time.NewTicker(EventListStreamTimeout)
	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			_ = closer()
			return events, nil

		case <-ticker.C:
			_ = closer()
			// return values if events are not streamed
			return events, nil

		case e, hasMore := <-eventStream:
			if !hasMore {
				_ = closer()
				return events, nil
			}

			processedEvent, err := ProcessEvent(e, ces.Opts.Event, req.EventName)
			if err != nil {
				ces.Logger.Warn(`event processing`, zap.Error(err))
			}

			if processedEvent != nil {
				events.Items = append(events.Items, processedEvent)
			}

			ticker.Reset(EventListStreamTimeout)
		}
	}
}

// EventsChan is not part of ChaincodeInstanceEventService interface, for calls as component
func (ces *ChaincodeInstanceEventService) EventsChan(ctx context.Context, rr ...*ChaincodeInstanceEventsStreamRequest) (_ chan *ChaincodeEvent, closer func() error, _ error) {
	req := &ChaincodeInstanceEventsStreamRequest{}
	if len(rr) == 1 {
		req = rr[0]
		if err := router.ValidateRequest(req); err != nil {
			return nil, nil, err
		}
	}

	for _, c := range ces.Opts.Context {
		ctx = c(ctx)
	}
	signer, _ := SignerFromContext(ctx)
	events, deliveryCloser, err := ces.EventDelivery.Events(
		ctx,
		ces.Locator.Channel,
		ces.Locator.Chaincode,
		signer,
		BlockRange(req.FromBlock, req.ToBlock)...,
	)

	if err != nil {
		return nil, nil, err
	}

	eventsProcessed := make(chan *ChaincodeEvent)

	closer = func() error {
		return deliveryCloser()
	}

	go func() {
		for e := range events {
			eventProcessed, err := ProcessEvent(e, ces.Opts.Event, req.EventName)

			if err != nil {
				ces.Logger.Warn(`event processing`, zap.Error(err))
			}
			if eventProcessed != nil {
				eventsProcessed <- eventProcessed
			}
		}
		close(eventsProcessed)
	}()

	return eventsProcessed, closer, nil
}

func ProcessEvent(event interface {
	Event() *peer.ChaincodeEvent
	Block() uint64
	TxTimestamp() *timestamp.Timestamp
}, opts []EventOpt, matchName []string) (*ChaincodeEvent, error) {

	processedEvent := &ChaincodeEvent{
		Event:       event.Event(),
		Block:       event.Block(),
		TxTimestamp: event.TxTimestamp(),
	}
	for _, o := range opts {
		if err := o(processedEvent); err != nil {
			return processedEvent, err
		}
	}

	if !MatchEventName(event.Event().EventName, matchName) {
		return nil, nil
	}

	return processedEvent, nil
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
