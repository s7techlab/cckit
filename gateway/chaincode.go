package gateway

import (
	"context"
	"errors"

	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/s7techlab/cckit/convert"
)

type Action string

const (
	Query  Action = `query`
	Invoke Action = `invoke`
)

// Chaincode interface for work with chaincode
type Chaincode interface {
	Query(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
	Invoke(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
	Events(ctx context.Context, r ...*ChaincodeInstanceEventsStreamRequest) (ChaincodeEventSub, error)
}

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

	Opts *Opts
}

func NewChaincode(service ChaincodeServiceServer, channelName, chaincodeName string, opts ...Opt) *chaincode {
	c := &chaincode{
		Service:   service,
		Channel:   channelName,
		Chaincode: chaincodeName,
		Opts:      &Opts{},
	}

	for _, opt := range opts {
		opt(c.Opts)
	}

	return c
}

func (g *chaincode) Events(ctx context.Context, r ...*ChaincodeInstanceEventsStreamRequest) (ChaincodeEventSub, error) {
	stream := NewChaincodeEventServerStream(ctx, g.Opts.Event...)

	req := &ChaincodeEventsStreamRequest{
		Chaincode: &ChaincodeLocator{
			Channel:   g.Channel,
			Chaincode: g.Chaincode,
		},
	}

	switch {
	case len(r) == 1:
		req.FromBlock = r[0].FromBlock
		req.ToBlock = r[0].ToBlock
		req.EventName = r[0].EventName

	case len(r) > 1:
		return nil, errors.New(`zero or one stream request allowed`)
	}

	go func() {
		err := g.Service.EventsStream(req, &ChaincodeEventsServer{ServerStream: stream})

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
		return g.ccOutput(c, Query, response, target)
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
		return g.ccOutput(c, Invoke, response, target)
	}
}

func (g *chaincode) context(ctx context.Context) context.Context {
	for _, c := range g.Opts.Context {
		ctx = c(ctx)
	}
	return ctx
}

func (g *chaincode) ccInput(ctx context.Context, action Action, fn string, args []interface{}) (ccInput *ChaincodeInput, err error) {
	var argsBytes [][]byte
	if argsBytes, err = convert.ArgsToBytes(args...); err != nil {
		return nil, err
	}

	ccInput = &ChaincodeInput{
		Chaincode: &ChaincodeLocator{
			Channel:   g.Channel,
			Chaincode: g.Chaincode,
		},
		Args: append([][]byte{[]byte(fn)}, argsBytes...),
	}

	if ccInput.Transient, err = TransientFromContext(ctx); err != nil {
		return nil, err
	}

	for _, i := range g.Opts.Input {
		if err = i(action, ccInput); err != nil {
			return nil, err
		}
	}

	return
}

func (g *chaincode) ccOutput(ctx context.Context, action Action, response *peer.Response, target interface{}) (res interface{}, err error) {
	for _, o := range g.Opts.Output {
		if err = o(action, response); err != nil {
			return nil, err
		}
	}
	return convert.FromBytes(response.Payload, target)
}
