package gateway

import (
	"context"

	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/gateway/service"
)

type Action string

const (
	Query  Action = `query`
	Invoke Action = `invoke`
)

type chaincode struct {
	Service   service.Chaincode
	Channel   string
	Chaincode string

	ContextOpts []ContextOpt
	InputOpts   []InputOpt
	OutputOpts  []OutputOpt
	EventOpts   []EventOpt
}

type Opt func(*chaincode)

type ContextOpt func(ctx context.Context) context.Context
type InputOpt func(action Action, input *service.ChaincodeInput) error
type OutputOpt func(action Action, response *peer.Response) error
type EventOpt func(event *peer.ChaincodeEvent) error

func NewChaincode(service service.Chaincode, channelName, chaincodeName string, opts ...Opt) *chaincode {
	c := &chaincode{
		Service:   service,
		Channel:   channelName,
		Chaincode: chaincodeName,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (g *chaincode) Events(ctx context.Context) (ChaincodeEventSub, error) {
	stream := NewChaincodeEventServerStream(ctx, g.EventOpts...)

	go func() {
		err := g.Service.Events(&service.ChaincodeLocator{
			Channel:   g.Channel,
			Chaincode: g.Chaincode,
		}, &service.ChaincodeEventsServer{ServerStream: stream})

		if err != nil {
			stream.Close()
		}
	}()

	return stream, nil
}

func (g *chaincode) Query(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {
	c := g.context(ctx)
	ccInput, err := g.ccInput(c, Query, fn, args)
	if err != nil {
		return nil, err
	}

	if response, err := g.Service.Query(c, ccInput); err != nil {
		return nil, err
	} else {
		return g.ccOutput(c, Query, response.Response, target)
	}
}

func (g *chaincode) Invoke(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {
	c := g.context(ctx)
	ccInput, err := g.ccInput(c, Invoke, fn, args)
	if err != nil {
		return nil, err
	}

	if response, err := g.Service.Invoke(c, ccInput); err != nil {
		return nil, err
	} else {
		return g.ccOutput(c, Invoke, response.Response, target)
	}
}

func (g *chaincode) context(ctx context.Context) context.Context {
	for _, c := range g.ContextOpts {
		ctx = c(ctx)
	}
	return ctx
}

func (g *chaincode) ccInput(ctx context.Context, action Action, fn string, args []interface{}) (ccInput *service.ChaincodeInput, err error) {
	var argsBytes [][]byte
	if argsBytes, err = convert.ArgsToBytes(args...); err != nil {
		return nil, err
	}

	ccInput = &service.ChaincodeInput{
		Channel:   g.Channel,
		Chaincode: g.Chaincode,
		Args:      append([][]byte{[]byte(fn)}, argsBytes...),
	}

	if ccInput.Transient, err = TransientFromContext(ctx); err != nil {
		return nil, err
	}

	for _, i := range g.InputOpts {
		if err = i(action, ccInput); err != nil {
			return nil, err
		}
	}

	return
}

func (g *chaincode) ccOutput(ctx context.Context, action Action, response *peer.Response, target interface{}) (res interface{}, err error) {
	for _, o := range g.OutputOpts {
		if err = o(action, response); err != nil {
			return nil, err
		}
	}
	return convert.FromBytes(response.Payload, target)
}
