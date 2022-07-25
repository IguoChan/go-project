package grpcx

import (
	"context"

	runtime2 "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type PBServerRegister interface {
	RegisterPBServer(server *grpc.Server)
}

type PBGatewayRegister interface {
	PBServerRegister
	RegisterPBGateway(ctx context.Context, mux *runtime2.ServeMux, rpcAddr string, opts []grpc.DialOption)
}

type EmptyGateway struct{}

func (EmptyGateway) RegisterPBGateway(ctx context.Context, mux *runtime2.ServeMux, rpcAddr string, opts []grpc.DialOption) {
}
