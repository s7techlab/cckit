package gateway

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type (
	RegisterHandlerFromEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

	ServiceDef struct {
		Desc                        *grpc.ServiceDesc
		Service                     interface{}
		HandlerFromEndpointRegister RegisterHandlerFromEndpoint
	}
)
