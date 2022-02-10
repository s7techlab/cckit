package gateway

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type (
	RegisterHandlerFromEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

	ServiceDef struct {
		SwaggerJSON                 []byte
		Desc                        *grpc.ServiceDesc
		Service                     interface{}
		HandlerFromEndpointRegister RegisterHandlerFromEndpoint
	}

	Service interface {
		Swagger() []byte
		GRPCDesc() *grpc.ServiceDesc
		Impl() interface{}
		GRPCGatewayRegister() RegisterHandlerFromEndpoint
	}
)

func (s ServiceDef) Swagger() []byte {
	return s.SwaggerJSON
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
