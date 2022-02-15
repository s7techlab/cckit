package gateway

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type (
	RegisterHandlerFromEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

	ServiceDef struct {
		name                        string
		swagger                     []byte
		Desc                        *grpc.ServiceDesc
		Service                     interface{}
		HandlerFromEndpointRegister RegisterHandlerFromEndpoint
	}

	Service interface {
		Name() string
		Swagger() []byte
		GRPCDesc() *grpc.ServiceDesc
		Impl() interface{}
		GRPCGatewayRegister() RegisterHandlerFromEndpoint
	}
)

func NewServiceDef(name string, swagger []byte, desc *grpc.ServiceDesc, service interface{}, registerHandler RegisterHandlerFromEndpoint) ServiceDef {
	return ServiceDef{
		name:                        name,
		swagger:                     swagger,
		Desc:                        desc,
		Service:                     service,
		HandlerFromEndpointRegister: registerHandler,
	}
}

func (s ServiceDef) Name() string {
	return s.name
}

func (s ServiceDef) Swagger() []byte {
	return s.swagger
}

func (s ServiceDef) GRPCDesc() *grpc.ServiceDesc {
	return s.Desc
}

func (s ServiceDef) Impl() interface{} {
	return s.Service
}

func (s ServiceDef) GRPCGatewayRegister() RegisterHandlerFromEndpoint {
	return s.HandlerFromEndpointRegister
}
