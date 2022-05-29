package demo_app

import (
	"context"
	"os"
	"time"

	"github.com/IguoChan/go-project/api/genproto/demo_app/simplepb"
	runtime2 "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type DemoServer struct {
	simplepb.UnimplementedSimpleServer
}

func NewDemoServer() DemoServer {
	return DemoServer{}
}

func (s DemoServer) RegisterPBServer(server *grpc.Server) {
	simplepb.RegisterSimpleServer(server, s)
}

func (s DemoServer) RegisterPBGateway(ctx context.Context, mux *runtime2.ServeMux, rpcAddr string, opts []grpc.DialOption) {
	if err := simplepb.RegisterSimpleHandlerFromEndpoint(ctx, mux, rpcAddr, opts); err != nil {
		logrus.Errorf("register gateway failed: %+v", err)
	}
}

func (s DemoServer) Route(ctx context.Context, request *simplepb.SimpleRequest) (*simplepb.SimpleResponse, error) {
	hn, _ := os.Hostname()
	res := &simplepb.SimpleResponse{
		Code:  200,
		Value: "DE " + hn + " " + request.Data + " " + time.Now().String(),
	}
	return res, nil
}
