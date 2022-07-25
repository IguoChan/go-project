package demo_app

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/IguoChan/go-project/pkg/grpcx"

	"github.com/IguoChan/go-project/api/genproto/demo_app/server_streampb"
	"google.golang.org/grpc"
)

type ServerStream struct {
	grpcx.EmptyGateway
	server_streampb.UnimplementedStreamServerServer
}

func NewServerStream() ServerStream {
	return ServerStream{}
}

func (s ServerStream) RegisterPBServer(server *grpc.Server) {
	server_streampb.RegisterStreamServerServer(server, s)
}

func (s ServerStream) Route(ctx context.Context, request *server_streampb.SimpleRequest) (*server_streampb.SimpleResponse, error) {
	hn, _ := os.Hostname()
	res := &server_streampb.SimpleResponse{
		Code:  200,
		Value: "DE " + hn + " " + request.Data + " " + time.Now().String(),
	}
	return res, nil
}

func (s ServerStream) ListValue(request *server_streampb.SimpleRequest, server server_streampb.StreamServer_ListValueServer) error {
	for n := 0; n < 5; n++ {
		// 向流中发送消息， 默认每次send送消息最大长度为`math.MaxInt32`bytes
		err := server.Send(&server_streampb.StreamResponse{
			StreamValue: request.Data + strconv.Itoa(n),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
