package gateway

import (
	"context"
	"time"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/sdk"
)

type ChaincodeInstanceEventService struct {
	EventDelivery sdk.EventDelivery
	Locator       *ChaincodeLocator
	Opts          *Opts
}

var _ ChaincodeInstanceEventsServiceServer = &ChaincodeInstanceEventService{}

func NewChaincodeInstanceEventService(delivery sdk.EventDelivery, channel, chaincode string, opts ...OptFunc) *ChaincodeInstanceEventService {
	cis := &ChaincodeInstanceEventService{
		EventDelivery: delivery,
		Locator: &ChaincodeLocator{
			Channel:   channel,
			Chaincode: chaincode,
		},
		Opts: &Opts{},
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

		case event, ok := <-events:
			if !ok {
				return nil
			}

			ccEvent := &ChaincodeEvent{
				Event:       event.Event(),
				Block:       event.Block(),
				TxTimestamp: event.TxTimestamp(),
			}
			for _, o := range ces.Opts.Event {
				_ = o(ccEvent)
			}

			if !MatchEventName(event.Event().EventName, req.EventName) {
				continue
			}

			if err = stream.Send(ccEvent); err != nil {
				return err
			}

		}
	}
}

func (ces *ChaincodeInstanceEventService) Events(ctx context.Context, req *ChaincodeInstanceEventsRequest) (*ChaincodeEvents, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	signer, _ := SignerFromContext(ctx)

	// for event _list_ default block range is up to current channel height
	if req.ToBlock == nil {
		req.ToBlock = &BlockLimit{Num: 0}
	}

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

		case event, hasMore := <-eventStream:
			if !hasMore {
				_ = closer()
				return events, nil
			}

			ccEvent := &ChaincodeEvent{
				Event:       event.Event(),
				Block:       event.Block(),
				TxTimestamp: event.TxTimestamp(),
			}

			for _, o := range ces.Opts.Event {
				_ = o(ccEvent)
			}

			if !MatchEventName(ccEvent.Event.EventName, req.EventName) {
				continue
			}

			events.Items = append(events.Items, ccEvent)
			ticker.Reset(EventListStreamTimeout)
		}
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
