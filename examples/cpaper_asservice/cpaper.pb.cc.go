// Code generated by protoc-gen-cc-gateway. DO NOT EDIT.
// source: cpaper_asservice/cpaper.proto

/*
Package cpaper_asservice contains
  *   chaincode methods names {service_name}Chaincode_{method_name}
  *   chaincode interface definition {service_name}Chaincode
  *   chaincode gateway definition {service_name}}Gateway
  *   chaincode service to cckit router registration func
*/
package cpaper_asservice

import (
	context "context"
	_ "embed"

	cckit_croscc "github.com/s7techlab/cckit/extensions/crosscc"
	cckit_gateway "github.com/s7techlab/cckit/gateway"
	cckit_router "github.com/s7techlab/cckit/router"
	cckit_defparam "github.com/s7techlab/cckit/router/param/defparam"
	cckit_sdk "github.com/s7techlab/cckit/sdk"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CPaperServiceChaincode method names
const (

	// CPaperServiceChaincodeMethodPrefix allows to use multiple services with same method names in one chaincode
	CPaperServiceChaincodeMethodPrefix = ""

	CPaperServiceChaincode_List = CPaperServiceChaincodeMethodPrefix + "List"

	CPaperServiceChaincode_Get = CPaperServiceChaincodeMethodPrefix + "Get"

	CPaperServiceChaincode_GetByExternalId = CPaperServiceChaincodeMethodPrefix + "GetByExternalId"

	CPaperServiceChaincode_Issue = CPaperServiceChaincodeMethodPrefix + "Issue"

	CPaperServiceChaincode_Buy = CPaperServiceChaincodeMethodPrefix + "Buy"

	CPaperServiceChaincode_Redeem = CPaperServiceChaincodeMethodPrefix + "Redeem"

	CPaperServiceChaincode_Delete = CPaperServiceChaincodeMethodPrefix + "Delete"
)

// CPaperServiceChaincode chaincode methods interface
type CPaperServiceChaincode interface {
	List(cckit_router.Context, *emptypb.Empty) (*CommercialPaperList, error)

	Get(cckit_router.Context, *CommercialPaperId) (*CommercialPaper, error)

	GetByExternalId(cckit_router.Context, *ExternalId) (*CommercialPaper, error)

	Issue(cckit_router.Context, *IssueCommercialPaper) (*CommercialPaper, error)

	Buy(cckit_router.Context, *BuyCommercialPaper) (*CommercialPaper, error)

	Redeem(cckit_router.Context, *RedeemCommercialPaper) (*CommercialPaper, error)

	Delete(cckit_router.Context, *CommercialPaperId) (*CommercialPaper, error)
}

// RegisterCPaperServiceChaincode registers service methods as chaincode router handlers
func RegisterCPaperServiceChaincode(r *cckit_router.Group, cc CPaperServiceChaincode) error {

	r.Query(CPaperServiceChaincode_List,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.List(ctx, ctx.Param().(*emptypb.Empty))
		},
		cckit_defparam.Proto(&emptypb.Empty{}))

	r.Query(CPaperServiceChaincode_Get,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Get(ctx, ctx.Param().(*CommercialPaperId))
		},
		cckit_defparam.Proto(&CommercialPaperId{}))

	r.Query(CPaperServiceChaincode_GetByExternalId,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.GetByExternalId(ctx, ctx.Param().(*ExternalId))
		},
		cckit_defparam.Proto(&ExternalId{}))

	r.Invoke(CPaperServiceChaincode_Issue,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Issue(ctx, ctx.Param().(*IssueCommercialPaper))
		},
		cckit_defparam.Proto(&IssueCommercialPaper{}))

	r.Invoke(CPaperServiceChaincode_Buy,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Buy(ctx, ctx.Param().(*BuyCommercialPaper))
		},
		cckit_defparam.Proto(&BuyCommercialPaper{}))

	r.Invoke(CPaperServiceChaincode_Redeem,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Redeem(ctx, ctx.Param().(*RedeemCommercialPaper))
		},
		cckit_defparam.Proto(&RedeemCommercialPaper{}))

	r.Invoke(CPaperServiceChaincode_Delete,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Delete(ctx, ctx.Param().(*CommercialPaperId))
		},
		cckit_defparam.Proto(&CommercialPaperId{}))

	return nil
}

//go:embed cpaper.swagger.json
var CPaperServiceSwagger []byte

// NewCPaperServiceGateway creates gateway to access chaincode method via chaincode service
func NewCPaperServiceGateway(sdk cckit_sdk.SDK, channel, chaincode string, opts ...cckit_gateway.OptFunc) *CPaperServiceGateway {
	return &CPaperServiceGateway{
		Invoker: &cckit_gateway.ChaincodeInstanceServiceInvoker{
			ChaincodeInstance: cckit_gateway.NewChaincodeInstanceService(
				sdk,
				&cckit_gateway.ChaincodeLocator{Channel: channel, Chaincode: chaincode},
				opts...),
		},
	}
}

// gateway implementation
// gateway can be used as kind of SDK, GRPC or REST server ( via grpc-gateway or clay )
type CPaperServiceGateway struct {
	Invoker cckit_gateway.ChaincodeInstanceInvoker
}

// ServiceDef returns service definition
func (c *CPaperServiceGateway) ServiceDef() cckit_gateway.ServiceDef {
	return cckit_gateway.NewServiceDef(
		_CPaperService_serviceDesc.ServiceName,
		CPaperServiceSwagger,
		&_CPaperService_serviceDesc,
		c,
		RegisterCPaperServiceHandlerFromEndpoint,
	)
}

func (c *CPaperServiceGateway) List(ctx context.Context, in *emptypb.Empty) (*CommercialPaperList, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx, CPaperServiceChaincode_List, []interface{}{in}, &CommercialPaperList{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaperList), nil
	}
}

func (c *CPaperServiceGateway) Get(ctx context.Context, in *CommercialPaperId) (*CommercialPaper, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx, CPaperServiceChaincode_Get, []interface{}{in}, &CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaper), nil
	}
}

func (c *CPaperServiceGateway) GetByExternalId(ctx context.Context, in *ExternalId) (*CommercialPaper, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx, CPaperServiceChaincode_GetByExternalId, []interface{}{in}, &CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaper), nil
	}
}

func (c *CPaperServiceGateway) Issue(ctx context.Context, in *IssueCommercialPaper) (*CommercialPaper, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Invoke(ctx, CPaperServiceChaincode_Issue, []interface{}{in}, &CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaper), nil
	}
}

func (c *CPaperServiceGateway) Buy(ctx context.Context, in *BuyCommercialPaper) (*CommercialPaper, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Invoke(ctx, CPaperServiceChaincode_Buy, []interface{}{in}, &CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaper), nil
	}
}

func (c *CPaperServiceGateway) Redeem(ctx context.Context, in *RedeemCommercialPaper) (*CommercialPaper, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Invoke(ctx, CPaperServiceChaincode_Redeem, []interface{}{in}, &CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaper), nil
	}
}

func (c *CPaperServiceGateway) Delete(ctx context.Context, in *CommercialPaperId) (*CommercialPaper, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Invoke(ctx, CPaperServiceChaincode_Delete, []interface{}{in}, &CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaper), nil
	}
}

// CPaperServiceChaincodeResolver interface for service resolver
type (
	CPaperServiceChaincodeResolver interface {
		Resolve(ctx cckit_router.Context) (CPaperServiceChaincode, error)
	}

	CPaperServiceChaincodeLocalResolver struct {
		service CPaperServiceChaincode
	}

	CPaperServiceChaincodeLocatorResolver struct {
		locatorResolver cckit_gateway.ChaincodeLocatorResolver
		service         CPaperServiceChaincode
	}
)

func NewCPaperServiceChaincodeLocalResolver(service CPaperServiceChaincode) *CPaperServiceChaincodeLocalResolver {
	return &CPaperServiceChaincodeLocalResolver{
		service: service,
	}
}

func (r *CPaperServiceChaincodeLocalResolver) Resolve(ctx cckit_router.Context) (CPaperServiceChaincode, error) {
	if r.service == nil {
		return nil, cckit_croscc.ErrServiceNotForLocalChaincodeResolver
	}

	return r.service, nil
}

func NewCPaperServiceChaincodeResolver(locatorResolver cckit_gateway.ChaincodeLocatorResolver) *CPaperServiceChaincodeLocatorResolver {
	return &CPaperServiceChaincodeLocatorResolver{
		locatorResolver: locatorResolver,
	}
}

func (r *CPaperServiceChaincodeLocatorResolver) Resolve(ctx cckit_router.Context) (CPaperServiceChaincode, error) {
	if r.service != nil {
		return r.service, nil
	}

	locator, err := r.locatorResolver(ctx, _CPaperService_serviceDesc.ServiceName)
	if err != nil {
		return nil, err
	}

	r.service = NewCPaperServiceChaincodeStubInvoker(locator)
	return r.service, nil
}

type CPaperServiceChaincodeStubInvoker struct {
	Invoker cckit_gateway.ChaincodeStubInvoker
}

func NewCPaperServiceChaincodeStubInvoker(locator *cckit_gateway.ChaincodeLocator) *CPaperServiceChaincodeStubInvoker {
	return &CPaperServiceChaincodeStubInvoker{
		Invoker: &cckit_gateway.LocatorChaincodeStubInvoker{Locator: locator},
	}
}

func (c *CPaperServiceChaincodeStubInvoker) List(ctx cckit_router.Context, in *emptypb.Empty) (*CommercialPaperList, error) {

	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx.Stub(), CPaperServiceChaincode_List, []interface{}{in}, &CommercialPaperList{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaperList), nil
	}

}

func (c *CPaperServiceChaincodeStubInvoker) Get(ctx cckit_router.Context, in *CommercialPaperId) (*CommercialPaper, error) {

	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx.Stub(), CPaperServiceChaincode_Get, []interface{}{in}, &CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaper), nil
	}

}

func (c *CPaperServiceChaincodeStubInvoker) GetByExternalId(ctx cckit_router.Context, in *ExternalId) (*CommercialPaper, error) {

	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker.Query(ctx.Stub(), CPaperServiceChaincode_GetByExternalId, []interface{}{in}, &CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaper), nil
	}

}

func (c *CPaperServiceChaincodeStubInvoker) Issue(ctx cckit_router.Context, in *IssueCommercialPaper) (*CommercialPaper, error) {

	return nil, cckit_gateway.ErrInvokeMethodNotAllowed

}

func (c *CPaperServiceChaincodeStubInvoker) Buy(ctx cckit_router.Context, in *BuyCommercialPaper) (*CommercialPaper, error) {

	return nil, cckit_gateway.ErrInvokeMethodNotAllowed

}

func (c *CPaperServiceChaincodeStubInvoker) Redeem(ctx cckit_router.Context, in *RedeemCommercialPaper) (*CommercialPaper, error) {

	return nil, cckit_gateway.ErrInvokeMethodNotAllowed

}

func (c *CPaperServiceChaincodeStubInvoker) Delete(ctx cckit_router.Context, in *CommercialPaperId) (*CommercialPaper, error) {

	return nil, cckit_gateway.ErrInvokeMethodNotAllowed

}
