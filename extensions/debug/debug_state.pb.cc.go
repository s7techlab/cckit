// Code generated by protoc-gen-cc-gateway. DO NOT EDIT.
// source: debug/debug_state.proto

/*
Package debug contains
  *   chaincode interface definition
  *   chaincode gateway definition
  *   chaincode service to cckit router registration func
*/
package debug

import (
	context "context"

	cckit_gateway "github.com/s7techlab/cckit/gateway"
	cckit_ccservice "github.com/s7techlab/cckit/gateway/service"
	cckit_router "github.com/s7techlab/cckit/router"
	cckit_param "github.com/s7techlab/cckit/router/param"
	cckit_defparam "github.com/s7techlab/cckit/router/param/defparam"
)

// DebugStateChaincode  method names
const (
	DebugStateChaincode_StateClean = "StateClean"

	DebugStateChaincode_StateKeys = "StateKeys"

	DebugStateChaincode_StateGet = "StateGet"

	DebugStateChaincode_StatePut = "StatePut"

	DebugStateChaincode_StateDelete = "StateDelete"
)

// DebugStateChaincodeResolver interface for service resolver
type DebugStateChaincodeResolver interface {
	DebugStateChaincode(ctx cckit_router.Context) (DebugStateChaincode, error)
}

// DebugStateChaincode chaincode methods interface
type DebugStateChaincode interface {
	StateClean(cckit_router.Context, *Prefixes) (*PrefixesMatchCount, error)

	StateKeys(cckit_router.Context, *Prefix) (*CompositeKeys, error)

	StateGet(cckit_router.Context, *CompositeKey) (*Value, error)

	StatePut(cckit_router.Context, *Value) (*Value, error)

	StateDelete(cckit_router.Context, *CompositeKey) (*Value, error)
}

// RegisterDebugStateChaincode registers service methods as chaincode router handlers
func RegisterDebugStateChaincode(r *cckit_router.Group, cc DebugStateChaincode) error {

	r.Invoke(DebugStateChaincode_StateClean,
		func(ctx cckit_router.Context) (interface{}, error) {
			if v, ok := ctx.Param().(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return nil, cckit_param.PayloadValidationError(err)
				}
			}
			return cc.StateClean(ctx, ctx.Param().(*Prefixes))
		},
		cckit_defparam.Proto(&Prefixes{}))

	r.Query(DebugStateChaincode_StateKeys,
		func(ctx cckit_router.Context) (interface{}, error) {
			if v, ok := ctx.Param().(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return nil, cckit_param.PayloadValidationError(err)
				}
			}
			return cc.StateKeys(ctx, ctx.Param().(*Prefix))
		},
		cckit_defparam.Proto(&Prefix{}))

	r.Query(DebugStateChaincode_StateGet,
		func(ctx cckit_router.Context) (interface{}, error) {
			if v, ok := ctx.Param().(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return nil, cckit_param.PayloadValidationError(err)
				}
			}
			return cc.StateGet(ctx, ctx.Param().(*CompositeKey))
		},
		cckit_defparam.Proto(&CompositeKey{}))

	r.Invoke(DebugStateChaincode_StatePut,
		func(ctx cckit_router.Context) (interface{}, error) {
			if v, ok := ctx.Param().(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return nil, cckit_param.PayloadValidationError(err)
				}
			}
			return cc.StatePut(ctx, ctx.Param().(*Value))
		},
		cckit_defparam.Proto(&Value{}))

	r.Invoke(DebugStateChaincode_StateDelete,
		func(ctx cckit_router.Context) (interface{}, error) {
			if v, ok := ctx.Param().(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return nil, cckit_param.PayloadValidationError(err)
				}
			}
			return cc.StateDelete(ctx, ctx.Param().(*CompositeKey))
		},
		cckit_defparam.Proto(&CompositeKey{}))

	return nil
}

// NewDebugStateGateway creates gateway to access chaincode method via chaincode service
func NewDebugStateGateway(ccService cckit_ccservice.Chaincode, channel, chaincode string, opts ...cckit_gateway.Opt) *DebugStateGateway {
	return &DebugStateGateway{Gateway: cckit_gateway.NewChaincode(ccService, channel, chaincode, opts...)}
}

// gateway implementation
// gateway can be used as kind of SDK, GRPC or REST server ( via grpc-gateway or clay )
type DebugStateGateway struct {
	Gateway cckit_gateway.Chaincode
}

// ServiceDef returns service definition
func (c *DebugStateGateway) ServiceDef() cckit_gateway.ServiceDef {
	return cckit_gateway.ServiceDef{
		Desc:                        &_DebugState_serviceDesc,
		Service:                     c,
		HandlerFromEndpointRegister: RegisterDebugStateHandlerFromEndpoint,
	}
}

// ApiDef deprecated, use ServiceDef
func (c *DebugStateGateway) ApiDef() cckit_gateway.ServiceDef {
	return c.ServiceDef()
}

// Events returns events subscription
func (c *DebugStateGateway) Events(ctx context.Context) (cckit_gateway.ChaincodeEventSub, error) {
	return c.Gateway.Events(ctx)
}

func (c *DebugStateGateway) StateClean(ctx context.Context, in *Prefixes) (*PrefixesMatchCount, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Invoke(ctx, DebugStateChaincode_StateClean, []interface{}{in}, &PrefixesMatchCount{}); err != nil {
		return nil, err
	} else {
		return res.(*PrefixesMatchCount), nil
	}
}

func (c *DebugStateGateway) StateKeys(ctx context.Context, in *Prefix) (*CompositeKeys, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Query(ctx, DebugStateChaincode_StateKeys, []interface{}{in}, &CompositeKeys{}); err != nil {
		return nil, err
	} else {
		return res.(*CompositeKeys), nil
	}
}

func (c *DebugStateGateway) StateGet(ctx context.Context, in *CompositeKey) (*Value, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Query(ctx, DebugStateChaincode_StateGet, []interface{}{in}, &Value{}); err != nil {
		return nil, err
	} else {
		return res.(*Value), nil
	}
}

func (c *DebugStateGateway) StatePut(ctx context.Context, in *Value) (*Value, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Invoke(ctx, DebugStateChaincode_StatePut, []interface{}{in}, &Value{}); err != nil {
		return nil, err
	} else {
		return res.(*Value), nil
	}
}

func (c *DebugStateGateway) StateDelete(ctx context.Context, in *CompositeKey) (*Value, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Invoke(ctx, DebugStateChaincode_StateDelete, []interface{}{in}, &Value{}); err != nil {
		return nil, err
	} else {
		return res.(*Value), nil
	}
}
