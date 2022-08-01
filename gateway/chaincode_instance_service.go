package gateway

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/sdk"
)

type (
	ChaincodeInstanceService struct {
		SDK     sdk.SDK
		Locator *ChaincodeLocator
		Opts    *Opts
	}

	// ChaincodeInstance base implementation ChaincodeInstanceService  - via SDK
	// but also possible another implementation (via REST, another SDK etc.)
	ChaincodeInstance interface {
		ChaincodeInstanceServiceServer

		EventsChan(ctx context.Context, rr ...*ChaincodeInstanceEventsStreamRequest) (
			_ chan *ChaincodeEvent, closer func() error, _ error)
	}
)

var _ ChaincodeInstance = &ChaincodeInstanceService{}

func NewChaincodeInstanceService(sdk sdk.SDK, locator *ChaincodeLocator, opts ...Opt) *ChaincodeInstanceService {
	ccInstanceService := &ChaincodeInstanceService{
		SDK:     sdk,
		Locator: locator,
		Opts:    &Opts{},
	}

	for _, o := range opts {
		o(ccInstanceService.Opts)
	}

	return ccInstanceService
}

func (cis *ChaincodeInstanceService) EventService() *ChaincodeInstanceEventService {
	return &ChaincodeInstanceEventService{
		EventDelivery: cis.SDK,
		Locator: &ChaincodeLocator{
			Channel:   cis.Locator.Channel,
			Chaincode: cis.Locator.Chaincode,
		},
		Opts: cis.Opts,
	}
}

func (cis *ChaincodeInstanceService) ServiceDef() ServiceDef {
	return ServiceDef{
		Desc:                        &_ChaincodeInstanceService_serviceDesc,
		Service:                     cis,
		HandlerFromEndpointRegister: RegisterChaincodeInstanceServiceHandlerFromEndpoint,
	}
}

func (cis *ChaincodeInstanceService) Exec(ctx context.Context, req *ChaincodeInstanceExecRequest) (*peer.Response, error) {
	switch req.Type {
	case InvocationType_INVOCATION_TYPE_QUERY:
		return cis.Query(ctx, &ChaincodeInstanceQueryRequest{Input: req.Input})
	case InvocationType_INVOCATION_TYPE_INVOKE:
		return cis.Invoke(ctx, &ChaincodeInstanceInvokeRequest{Input: req.Input})
	default:
		return nil, ErrUnknownInvocationType
	}
}

func (cis *ChaincodeInstanceService) Query(ctx context.Context, req *ChaincodeInstanceQueryRequest) (*peer.Response, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	for _, c := range cis.Opts.Context {
		ctx = c(ctx)
	}

	signer, _ := SignerFromContext(ctx)

	for _, i := range cis.Opts.Input {
		if err := i(req.Input); err != nil {
			return nil, err
		}
	}

	response, err := cis.SDK.Query(
		ctx,
		cis.Locator.Channel,
		cis.Locator.Chaincode,
		req.Input.Args,
		signer,
		req.Input.Transient,
	)
	if err != nil {
		return nil, fmt.Errorf("query chaincode: %w", err)
	}
	for _, o := range cis.Opts.Output {
		if err = o(InvocationType_INVOCATION_TYPE_QUERY, response); err != nil {
			return nil, fmt.Errorf(`output opt: %w`, err)
		}
	}

	return response, nil

}

func (cis *ChaincodeInstanceService) Invoke(ctx context.Context, req *ChaincodeInstanceInvokeRequest) (*peer.Response, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	for _, c := range cis.Opts.Context {
		ctx = c(ctx)
	}

	// underlying hlf-sdk(or your implementation must handle it) should handle 'nil' identity cases
	// and set default if identity wasn't provided here
	// if smth goes wrong we'll see it on the step below
	signer, _ := SignerFromContext(ctx)

	for _, i := range cis.Opts.Input {
		if err := i(req.Input); err != nil {
			return nil, err
		}
	}

	response, _, err := cis.SDK.Invoke(
		ctx,
		cis.Locator.Channel,
		cis.Locator.Chaincode,
		req.Input.Args,
		signer,
		req.Input.Transient,
		TxWaiterFromContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("invoke chaincode: %w", err)
	}

	for _, o := range cis.Opts.Output {
		if err = o(InvocationType_INVOCATION_TYPE_INVOKE, response); err != nil {
			return nil, fmt.Errorf(`output opt: %w`, err)
		}
	}

	return response, nil
}

func (cis *ChaincodeInstanceService) EventsStream(
	req *ChaincodeInstanceEventsStreamRequest, stream ChaincodeInstanceService_EventsStreamServer) error {
	return cis.EventService().EventsStream(req, stream)
}

func (cis *ChaincodeInstanceService) Events(
	ctx context.Context, req *ChaincodeInstanceEventsRequest) (*ChaincodeEvents, error) {
	return cis.EventService().Events(ctx, req)
}

func (cis *ChaincodeInstanceService) EventsChan(
	ctx context.Context, rr ...*ChaincodeInstanceEventsStreamRequest) (_ chan *ChaincodeEvent, closer func() error, _ error) {
	return cis.EventService().EventsChan(ctx, rr...)
}
