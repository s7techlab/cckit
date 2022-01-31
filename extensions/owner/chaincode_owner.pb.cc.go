// Code generated by protoc-gen-cc-gateway. DO NOT EDIT.
// source: owner/chaincode_owner.proto

/*
Package owner contains
  *   chaincode methods names {service_name}Chaincode_{method_name}
  *   chaincode interface definition {service_name}Chaincode
  *   chaincode gateway definition {service_name}}Gateway
  *   chaincode service to cckit router registration func
*/
package owner

import (
	context "context"
	_ "embed"

	cckit_gateway "github.com/s7techlab/cckit/gateway"
	cckit_router "github.com/s7techlab/cckit/router"
	cckit_defparam "github.com/s7techlab/cckit/router/param/defparam"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ChaincodeOwnerServiceChaincode  method names
const (
	ChaincodeOwnerServiceChaincode_GetOwnerByTxCreator = "GetOwnerByTxCreator"

	ChaincodeOwnerServiceChaincode_ListOwners = "ListOwners"

	ChaincodeOwnerServiceChaincode_GetOwner = "GetOwner"

	ChaincodeOwnerServiceChaincode_CreateOwner = "CreateOwner"

	ChaincodeOwnerServiceChaincode_CreateOwnerTxCreator = "CreateOwnerTxCreator"

	ChaincodeOwnerServiceChaincode_UpdateOwner = "UpdateOwner"

	ChaincodeOwnerServiceChaincode_DeleteOwner = "DeleteOwner"
)

// ChaincodeOwnerServiceChaincodeResolver interface for service resolver
type ChaincodeOwnerServiceChaincodeResolver interface {
	ChaincodeOwnerServiceChaincode(ctx cckit_router.Context) (ChaincodeOwnerServiceChaincode, error)
}

// ChaincodeOwnerServiceChaincode chaincode methods interface
type ChaincodeOwnerServiceChaincode interface {
	GetOwnerByTxCreator(cckit_router.Context, *emptypb.Empty) (*ChaincodeOwner, error)

	ListOwners(cckit_router.Context, *emptypb.Empty) (*ChaincodeOwners, error)

	GetOwner(cckit_router.Context, *OwnerId) (*ChaincodeOwner, error)

	CreateOwner(cckit_router.Context, *CreateOwnerRequest) (*ChaincodeOwner, error)

	CreateOwnerTxCreator(cckit_router.Context, *emptypb.Empty) (*ChaincodeOwner, error)

	UpdateOwner(cckit_router.Context, *UpdateOwnerRequest) (*ChaincodeOwner, error)

	DeleteOwner(cckit_router.Context, *OwnerId) (*ChaincodeOwner, error)
}

// RegisterChaincodeOwnerServiceChaincode registers service methods as chaincode router handlers
func RegisterChaincodeOwnerServiceChaincode(r *cckit_router.Group, cc ChaincodeOwnerServiceChaincode) error {

	r.Query(ChaincodeOwnerServiceChaincode_GetOwnerByTxCreator,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.GetOwnerByTxCreator(ctx, ctx.Param().(*emptypb.Empty))
		},
		cckit_defparam.Proto(&emptypb.Empty{}))

	r.Query(ChaincodeOwnerServiceChaincode_ListOwners,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.ListOwners(ctx, ctx.Param().(*emptypb.Empty))
		},
		cckit_defparam.Proto(&emptypb.Empty{}))

	r.Query(ChaincodeOwnerServiceChaincode_GetOwner,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.GetOwner(ctx, ctx.Param().(*OwnerId))
		},
		cckit_defparam.Proto(&OwnerId{}))

	r.Invoke(ChaincodeOwnerServiceChaincode_CreateOwner,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.CreateOwner(ctx, ctx.Param().(*CreateOwnerRequest))
		},
		cckit_defparam.Proto(&CreateOwnerRequest{}))

	r.Invoke(ChaincodeOwnerServiceChaincode_CreateOwnerTxCreator,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.CreateOwnerTxCreator(ctx, ctx.Param().(*emptypb.Empty))
		},
		cckit_defparam.Proto(&emptypb.Empty{}))

	r.Invoke(ChaincodeOwnerServiceChaincode_UpdateOwner,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.UpdateOwner(ctx, ctx.Param().(*UpdateOwnerRequest))
		},
		cckit_defparam.Proto(&UpdateOwnerRequest{}))

	r.Invoke(ChaincodeOwnerServiceChaincode_DeleteOwner,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.DeleteOwner(ctx, ctx.Param().(*OwnerId))
		},
		cckit_defparam.Proto(&OwnerId{}))

	return nil
}

//go:embed chaincode_owner.swagger.json
var ChaincodeOwnerServiceSwagger []byte

// NewChaincodeOwnerServiceGateway creates gateway to access chaincode method via chaincode service
func NewChaincodeOwnerServiceGateway(ccService cckit_gateway.ChaincodeServiceServer, channel, chaincode string, opts ...cckit_gateway.Opt) *ChaincodeOwnerServiceGateway {
	return &ChaincodeOwnerServiceGateway{Gateway: cckit_gateway.NewChaincode(ccService, channel, chaincode, opts...)}
}

// gateway implementation
// gateway can be used as kind of SDK, GRPC or REST server ( via grpc-gateway or clay )
type ChaincodeOwnerServiceGateway struct {
	Gateway cckit_gateway.Chaincode
}

// ServiceDef returns service definition
func (c *ChaincodeOwnerServiceGateway) ServiceDef() cckit_gateway.ServiceDef {
	return cckit_gateway.ServiceDef{
		Desc:                        &_ChaincodeOwnerService_serviceDesc,
		Service:                     c,
		HandlerFromEndpointRegister: RegisterChaincodeOwnerServiceHandlerFromEndpoint,
	}
}

// ApiDef deprecated, use ServiceDef
func (c *ChaincodeOwnerServiceGateway) ApiDef() cckit_gateway.ServiceDef {
	return c.ServiceDef()
}

// Events returns events subscription
func (c *ChaincodeOwnerServiceGateway) Events(ctx context.Context) (cckit_gateway.ChaincodeEventSub, error) {
	return c.Gateway.Events(ctx)
}

func (c *ChaincodeOwnerServiceGateway) GetOwnerByTxCreator(ctx context.Context, in *emptypb.Empty) (*ChaincodeOwner, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Query(ctx, ChaincodeOwnerServiceChaincode_GetOwnerByTxCreator, []interface{}{in}, &ChaincodeOwner{}); err != nil {
		return nil, err
	} else {
		return res.(*ChaincodeOwner), nil
	}
}

func (c *ChaincodeOwnerServiceGateway) ListOwners(ctx context.Context, in *emptypb.Empty) (*ChaincodeOwners, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Query(ctx, ChaincodeOwnerServiceChaincode_ListOwners, []interface{}{in}, &ChaincodeOwners{}); err != nil {
		return nil, err
	} else {
		return res.(*ChaincodeOwners), nil
	}
}

func (c *ChaincodeOwnerServiceGateway) GetOwner(ctx context.Context, in *OwnerId) (*ChaincodeOwner, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Query(ctx, ChaincodeOwnerServiceChaincode_GetOwner, []interface{}{in}, &ChaincodeOwner{}); err != nil {
		return nil, err
	} else {
		return res.(*ChaincodeOwner), nil
	}
}

func (c *ChaincodeOwnerServiceGateway) CreateOwner(ctx context.Context, in *CreateOwnerRequest) (*ChaincodeOwner, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Invoke(ctx, ChaincodeOwnerServiceChaincode_CreateOwner, []interface{}{in}, &ChaincodeOwner{}); err != nil {
		return nil, err
	} else {
		return res.(*ChaincodeOwner), nil
	}
}

func (c *ChaincodeOwnerServiceGateway) CreateOwnerTxCreator(ctx context.Context, in *emptypb.Empty) (*ChaincodeOwner, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Invoke(ctx, ChaincodeOwnerServiceChaincode_CreateOwnerTxCreator, []interface{}{in}, &ChaincodeOwner{}); err != nil {
		return nil, err
	} else {
		return res.(*ChaincodeOwner), nil
	}
}

func (c *ChaincodeOwnerServiceGateway) UpdateOwner(ctx context.Context, in *UpdateOwnerRequest) (*ChaincodeOwner, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Invoke(ctx, ChaincodeOwnerServiceChaincode_UpdateOwner, []interface{}{in}, &ChaincodeOwner{}); err != nil {
		return nil, err
	} else {
		return res.(*ChaincodeOwner), nil
	}
}

func (c *ChaincodeOwnerServiceGateway) DeleteOwner(ctx context.Context, in *OwnerId) (*ChaincodeOwner, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Gateway.Invoke(ctx, ChaincodeOwnerServiceChaincode_DeleteOwner, []interface{}{in}, &ChaincodeOwner{}); err != nil {
		return nil, err
	} else {
		return res.(*ChaincodeOwner), nil
	}
}
