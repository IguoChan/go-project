package demo_app

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/IguoChan/go-project/pkg/util"

	"github.com/IguoChan/go-project/api/genproto/demo_app/simplepb"
	runtime2 "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type SimpleServer struct {
	simplepb.UnimplementedSimpleServer
	Ext string
}

func NewSimpleServer() SimpleServer {
	return SimpleServer{}
}

func (s SimpleServer) RegisterPBServer(server *grpc.Server) {
	simplepb.RegisterSimpleServer(server, s)
}

func (s SimpleServer) RegisterPBGateway(ctx context.Context, mux *runtime2.ServeMux, rpcAddr string, opts []grpc.DialOption) {
	if err := simplepb.RegisterSimpleHandlerFromEndpoint(ctx, mux, rpcAddr, opts); err != nil {
		logrus.Errorf("register gateway failed: %+v", err)
	}
}

func (s SimpleServer) Route(ctx context.Context, request *simplepb.SimpleRequest) (*simplepb.SimpleResponse, error) {
	hn, _ := os.Hostname()
	res := &simplepb.SimpleResponse{
		Code:  200,
		Value: strconv.FormatInt(util.GoroutineId(), 10) + " DE " + hn + " " + request.Data + " " + time.Now().String() + " " + s.Ext,
	}
	return res, nil
}
