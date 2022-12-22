// Code generated by protoc-gen-cc-gateway. DO NOT EDIT.
// source: erc20_service/erc20.proto

/*
Package erc20_service contains
  *   chaincode methods names {service_name}Chaincode_{method_name}
  *   chaincode interface definition {service_name}Chaincode
  *   chaincode gateway definition {service_name}}Gateway
  *   chaincode service to cckit router registration func
*/
package erc20_service

import (
	context "context"
	_ "embed"
	errors "errors"

	cckit_gateway "github.com/s7techlab/cckit/gateway"
	cckit_router "github.com/s7techlab/cckit/router"
	cckit_defparam "github.com/s7techlab/cckit/router/param/defparam"
	cckit_sdk "github.com/s7techlab/cckit/sdk"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ERC20Chaincode method names
const (

	// ERC20ChaincodeMethodPrefix allows to use multiple services with same method names in one chaincode
	ERC20ChaincodeMethodPrefix = ""

	ERC20Chaincode_Name = ERC20ChaincodeMethodPrefix + "Name"

	ERC20Chaincode_Symbol = ERC20ChaincodeMethodPrefix + "Symbol"

	ERC20Chaincode_Decimals = ERC20ChaincodeMethodPrefix + "Decimals"

	ERC20Chaincode_TotalSupply = ERC20ChaincodeMethodPrefix + "TotalSupply"

	ERC20Chaincode_BalanceOf = ERC20ChaincodeMethodPrefix + "BalanceOf"

	ERC20Chaincode_Transfer = ERC20ChaincodeMethodPrefix + "Transfer"

	ERC20Chaincode_Allowance = ERC20ChaincodeMethodPrefix + "Allowance"

	ERC20Chaincode_Approve = ERC20ChaincodeMethodPrefix + "Approve"

	ERC20Chaincode_TransferFrom = ERC20ChaincodeMethodPrefix + "TransferFrom"
)

// ERC20Chaincode chaincode methods interface
type ERC20Chaincode interface {
	Name(cckit_router.Context, *emptypb.Empty) (*NameResponse, error)

	Symbol(cckit_router.Context, *emptypb.Empty) (*SymbolResponse, error)

	Decimals(cckit_router.Context, *emptypb.Empty) (*DecimalsResponse, error)

	TotalSupply(cckit_router.Context, *emptypb.Empty) (*TotalSupplyResponse, error)

	BalanceOf(cckit_router.Context, *BalanceOfRequest) (*BalanceOfResponse, error)

	Transfer(cckit_router.Context, *TransferRequest) (*TransferResponse, error)

	Allowance(cckit_router.Context, *AllowanceRequest) (*AllowanceResponse, error)

	Approve(cckit_router.Context, *ApproveRequest) (*ApproveResponse, error)

	TransferFrom(cckit_router.Context, *TransferFromRequest) (*TransferResponse, error)
}

// RegisterERC20Chaincode registers service methods as chaincode router handlers
func RegisterERC20Chaincode(r *cckit_router.Group, cc ERC20Chaincode) error {

	r.Query(ERC20Chaincode_Name,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Name(ctx, ctx.Param().(*emptypb.Empty))
		},
		cckit_defparam.Proto(&emptypb.Empty{}))

	r.Query(ERC20Chaincode_Symbol,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Symbol(ctx, ctx.Param().(*emptypb.Empty))
		},
		cckit_defparam.Proto(&emptypb.Empty{}))

	r.Query(ERC20Chaincode_Decimals,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Decimals(ctx, ctx.Param().(*emptypb.Empty))
		},
		cckit_defparam.Proto(&emptypb.Empty{}))

	r.Query(ERC20Chaincode_TotalSupply,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.TotalSupply(ctx, ctx.Param().(*emptypb.Empty))
		},
		cckit_defparam.Proto(&emptypb.Empty{}))

	r.Query(ERC20Chaincode_BalanceOf,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.BalanceOf(ctx, ctx.Param().(*BalanceOfRequest))
		},
		cckit_defparam.Proto(&BalanceOfRequest{}))

	r.Invoke(ERC20Chaincode_Transfer,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Transfer(ctx, ctx.Param().(*TransferRequest))
		},
		cckit_defparam.Proto(&TransferRequest{}))

	r.Query(ERC20Chaincode_Allowance,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Allowance(ctx, ctx.Param().(*AllowanceRequest))
		},
		cckit_defparam.Proto(&AllowanceRequest{}))

	r.Invoke(ERC20Chaincode_Approve,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Approve(ctx, ctx.Param().(*ApproveRequest))
		},
		cckit_defparam.Proto(&ApproveRequest{}))

	r.Invoke(ERC20Chaincode_TransferFrom,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.TransferFrom(ctx, ctx.Param().(*TransferFromRequest))
		},
		cckit_defparam.Proto(&TransferFromRequest{}))

	return nil
}

//go:embed erc20.swagger.json
var ERC20Swagger []byte

// NewERC20Gateway creates gateway to access chaincode method via chaincode service
func NewERC20Gateway(sdk cckit_sdk.SDK, channel, chaincode string, opts ...cckit_gateway.Opt) *ERC20Gateway {
	return NewERC20GatewayFromInstance(
		cckit_gateway.NewChaincodeInstanceService(
			sdk,
			&cckit_gateway.ChaincodeLocator{Channel: channel, Chaincode: chaincode},
			opts...,
		))
}

func NewERC20GatewayFromInstance(chaincodeInstance cckit_gateway.ChaincodeInstance) *ERC20Gateway {
	return &ERC20Gateway{
		ChaincodeInstance: chaincodeInstance,
	}
}

// gateway implementation
// gateway can be used as kind of SDK, GRPC or REST server ( via grpc-gateway or clay )
type ERC20Gateway struct {
	ChaincodeInstance cckit_gateway.ChaincodeInstance
}

func (c *ERC20Gateway) Invoker() cckit_gateway.ChaincodeInstanceInvoker {
	return cckit_gateway.NewChaincodeInstanceServiceInvoker(c.ChaincodeInstance)
}

// ServiceDef returns service definition
func (c *ERC20Gateway) ServiceDef() cckit_gateway.ServiceDef {
	return cckit_gateway.NewServiceDef(
		_ERC20_serviceDesc.ServiceName,
		ERC20Swagger,
		&_ERC20_serviceDesc,
		c,
		RegisterERC20HandlerFromEndpoint,
	)
}

func (c *ERC20Gateway) Name(ctx context.Context, in *emptypb.Empty) (*NameResponse, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker().Query(ctx, ERC20Chaincode_Name, []interface{}{in}, &NameResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*NameResponse), nil
	}
}

func (c *ERC20Gateway) Symbol(ctx context.Context, in *emptypb.Empty) (*SymbolResponse, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker().Query(ctx, ERC20Chaincode_Symbol, []interface{}{in}, &SymbolResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*SymbolResponse), nil
	}
}

func (c *ERC20Gateway) Decimals(ctx context.Context, in *emptypb.Empty) (*DecimalsResponse, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker().Query(ctx, ERC20Chaincode_Decimals, []interface{}{in}, &DecimalsResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*DecimalsResponse), nil
	}
}

func (c *ERC20Gateway) TotalSupply(ctx context.Context, in *emptypb.Empty) (*TotalSupplyResponse, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker().Query(ctx, ERC20Chaincode_TotalSupply, []interface{}{in}, &TotalSupplyResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*TotalSupplyResponse), nil
	}
}

func (c *ERC20Gateway) BalanceOf(ctx context.Context, in *BalanceOfRequest) (*BalanceOfResponse, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker().Query(ctx, ERC20Chaincode_BalanceOf, []interface{}{in}, &BalanceOfResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*BalanceOfResponse), nil
	}
}

func (c *ERC20Gateway) Transfer(ctx context.Context, in *TransferRequest) (*TransferResponse, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker().Invoke(ctx, ERC20Chaincode_Transfer, []interface{}{in}, &TransferResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*TransferResponse), nil
	}
}

func (c *ERC20Gateway) Allowance(ctx context.Context, in *AllowanceRequest) (*AllowanceResponse, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker().Query(ctx, ERC20Chaincode_Allowance, []interface{}{in}, &AllowanceResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*AllowanceResponse), nil
	}
}

func (c *ERC20Gateway) Approve(ctx context.Context, in *ApproveRequest) (*ApproveResponse, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker().Invoke(ctx, ERC20Chaincode_Approve, []interface{}{in}, &ApproveResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*ApproveResponse), nil
	}
}

func (c *ERC20Gateway) TransferFrom(ctx context.Context, in *TransferFromRequest) (*TransferResponse, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker().Invoke(ctx, ERC20Chaincode_TransferFrom, []interface{}{in}, &TransferResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*TransferResponse), nil
	}
}

// ERC20ChaincodeResolver interface for service resolver
type (
	ERC20ChaincodeResolver interface {
		Resolve(ctx cckit_router.Context) (ERC20Chaincode, error)
	}

	ERC20ChaincodeLocalResolver struct {
		service ERC20Chaincode
	}

	ERC20ChaincodeLocatorResolver struct {
		locatorResolver cckit_gateway.ChaincodeLocatorResolver
		service         ERC20Chaincode
	}
)

func NewERC20ChaincodeLocalResolver(service ERC20Chaincode) *ERC20ChaincodeLocalResolver {
	return &ERC20ChaincodeLocalResolver{
		service: service,
	}
}

func (r *ERC20ChaincodeLocalResolver) Resolve(ctx cckit_router.Context) (ERC20Chaincode, error) {
	if r.service == nil {
		return nil, errors.New("service not set for local chaincode resolver")
	}

	return r.service, nil
}

func NewERC20ChaincodeResolver(locatorResolver cckit_gateway.ChaincodeLocatorResolver) *ERC20ChaincodeLocatorResolver {
	return &ERC20ChaincodeLocatorResolver{
		locatorResolver: locatorResolver,
	}
}

func (r *ERC20ChaincodeLocatorResolver) Resolve(ctx cckit_router.Context) (ERC20Chaincode, error) {
	if r.service != nil {
		return r.service, nil
	}

	locator, err := r.locatorResolver(ctx, _ERC20_serviceDesc.ServiceName)
	if err != nil {
		return nil, err
	}

	r.service = NewERC20ChaincodeStubInvoker(locator)
	return r.service, nil
}

type ERC20ChaincodeStubInvoker struct {
	Invoker cckit_gateway.ChaincodeStubInvoker
}

func NewERC20ChaincodeStubInvoker(locator *cckit_gateway.ChaincodeLocator) *ERC20ChaincodeStubInvoker {
	return &ERC20ChaincodeStubInvoker{
		Invoker: &cckit_gateway.LocatorChaincodeStubInvoker{Locator: locator},
	}
}

func (c *ERC20ChaincodeStubInvoker) Name(ctx cckit_router.Context, in *emptypb.Empty) (*NameResponse, error) {

	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx.Stub(), ERC20Chaincode_Name, []interface{}{in}, &NameResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*NameResponse), nil
	}

}

func (c *ERC20ChaincodeStubInvoker) Symbol(ctx cckit_router.Context, in *emptypb.Empty) (*SymbolResponse, error) {

	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx.Stub(), ERC20Chaincode_Symbol, []interface{}{in}, &SymbolResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*SymbolResponse), nil
	}

}

func (c *ERC20ChaincodeStubInvoker) Decimals(ctx cckit_router.Context, in *emptypb.Empty) (*DecimalsResponse, error) {

	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx.Stub(), ERC20Chaincode_Decimals, []interface{}{in}, &DecimalsResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*DecimalsResponse), nil
	}

}

func (c *ERC20ChaincodeStubInvoker) TotalSupply(ctx cckit_router.Context, in *emptypb.Empty) (*TotalSupplyResponse, error) {

	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx.Stub(), ERC20Chaincode_TotalSupply, []interface{}{in}, &TotalSupplyResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*TotalSupplyResponse), nil
	}

}

func (c *ERC20ChaincodeStubInvoker) BalanceOf(ctx cckit_router.Context, in *BalanceOfRequest) (*BalanceOfResponse, error) {

	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx.Stub(), ERC20Chaincode_BalanceOf, []interface{}{in}, &BalanceOfResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*BalanceOfResponse), nil
	}

}

func (c *ERC20ChaincodeStubInvoker) Transfer(ctx cckit_router.Context, in *TransferRequest) (*TransferResponse, error) {

	return nil, cckit_gateway.ErrInvokeMethodNotAllowed

}

func (c *ERC20ChaincodeStubInvoker) Allowance(ctx cckit_router.Context, in *AllowanceRequest) (*AllowanceResponse, error) {

	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx.Stub(), ERC20Chaincode_Allowance, []interface{}{in}, &AllowanceResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*AllowanceResponse), nil
	}

}

func (c *ERC20ChaincodeStubInvoker) Approve(ctx cckit_router.Context, in *ApproveRequest) (*ApproveResponse, error) {

	return nil, cckit_gateway.ErrInvokeMethodNotAllowed

}

func (c *ERC20ChaincodeStubInvoker) TransferFrom(ctx cckit_router.Context, in *TransferFromRequest) (*TransferResponse, error) {

	return nil, cckit_gateway.ErrInvokeMethodNotAllowed

}