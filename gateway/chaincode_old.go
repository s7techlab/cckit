package gateway

import (
	"context"

	"github.com/s7techlab/cckit/convert"
)

// All this components are deprecated and will be removed,  use ChaincodeInvoker
// ====================================================================

type Action string

const (
	Query  Action = `query`
	Invoke Action = `invoke`
)

// Deprecated: use ChaincodeInvoker
type Chaincode interface {
	Query(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
	Invoke(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
	Events(ctx context.Context, r ...*ChaincodeInstanceEventsStreamRequest) (ChaincodeEventSub, error)
}

// Deprecated: use ChaincodeInstanceEventDelivery
type ChaincodeEventSub interface {
	Context() context.Context
	Events() <-chan *ChaincodeEvent
	Recv(*ChaincodeEvent) error
	Close()
}

type chaincode struct {
	Service   ChaincodeServiceServer
	Channel   string
	Chaincode string
}

// Deprecated: use NewChaincodeInvoker
func NewChaincode(service ChaincodeServiceServer, channelName, chaincodeName string, opts ...OptFunc) *chaincode {
	c := &chaincode{
		Service:   service,
		Channel:   channelName,
		Chaincode: chaincodeName,
	}

	return c
}

// Deprecated: use ChaincodeInstanceEventDelivery
func (c *chaincode) Events(ctx context.Context, r ...*ChaincodeInstanceEventsStreamRequest) (ChaincodeEventSub, error) {
	panic(`use ChaincodeInstanceEventDelivery`)
}

func (c *chaincode) Locator() *ChaincodeLocator {
	return &ChaincodeLocator{
		Channel:   c.Channel,
		Chaincode: c.Chaincode,
	}
}
func (c *chaincode) Query(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {
	ccInput, err := ccInput(ctx, fn, args)
	if err != nil {
		return nil, err
	}

	response, err := c.Service.Query(ctx, &ChaincodeQueryRequest{
		Locator: c.Locator(),
		Input:   ccInput,
	})
	if err != nil {
		return nil, err
	}

	return convert.FromBytes(response.Payload, target)
}

func (c *chaincode) Invoke(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {
	ccInput, err := ccInput(ctx, fn, args)
	if err != nil {
		return nil, err
	}

	response, err := c.Service.Invoke(ctx, &ChaincodeInvokeRequest{
		Locator: c.Locator(),
		Input:   ccInput,
	})
	if err != nil {
		return nil, err
	}

	return convert.FromBytes(response.Payload, target)
}
