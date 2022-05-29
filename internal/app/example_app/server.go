package example_app

import (
	"context"
	"os"
	"time"

	"github.com/IguoChan/go-project/api/genproto/example_app/examplepb"
	runtime2 "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ExampleServer struct {
	examplepb.UnimplementedExampleServer
}

func NewExampleServer() ExampleServer {
	return ExampleServer{}
}

func (e ExampleServer) RegisterPBServer(server *grpc.Server) {
	examplepb.RegisterExampleServer(server, e)
}

func (e ExampleServer) RegisterPBGateway(ctx context.Context, mux *runtime2.ServeMux, rpcAddr string, opts []grpc.DialOption) {
	if err := examplepb.RegisterExampleHandlerFromEndpoint(ctx, mux, rpcAddr, opts); err != nil {
		logrus.Errorf("register gateway failed: %+v", err)
	}
}

func (e ExampleServer) R(ctx context.Context, req *examplepb.ExReq) (*examplepb.ExResp, error) {
	hn, _ := os.Hostname()
	res := &examplepb.ExResp{
		Code:  200,
		Value: "EX " + hn + " " + req.S.Data + " " + time.Now().String(),
	}
	return res, nil
}
