package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/IguoChan/go-project/api/genproto/demo_app/server_streampb"

	"github.com/IguoChan/go-project/api/genproto/demo_app/simplepb"

	"github.com/IguoChan/go-project/internal/app/demo_app"

	"github.com/IguoChan/go-project/pkg/etcdx"
	"github.com/IguoChan/go-project/pkg/rpcx"
)

func main() {
	os.Exit(Run())
}

func Run() int {
	opts := make(map[string]*rpcx.ClientOptions)
	opts["demo"] = &rpcx.ClientOptions{
		EtcdOpt: &etcdx.Options{
			Addrs:       []string{"192.168.0.98:2379"},
			DialTimeout: 5 * time.Second,
			TTL:         30,
		},
		ServiceName: "demo",
	}

	dc := demo_app.NewRpcClient(opts)
	s, err := dc.Simple()
	if err != nil {
		panic(err)
	}
	req := simplepb.SimpleRequest{
		Data: "grpc 100",
	}
	for i := 0; i < 100; i++ {
		res, err := s.Route(context.Background(), &req)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}

	ss, err := dc.SS()
	if err != nil {
		panic(err)
	}

	req1 := &server_streampb.SimpleRequest{Data: "stream server grpc "}
	stream, err := ss.ListValue(context.Background(), req1)
	if err != nil {
		panic(err)
	}
	for {
		//Recv() 方法接收服务端消息，默认每次Recv()最大消息长度为`1024*1024*4`bytes(4M)
		res, err := stream.Recv()
		// 判断消息流是否已经结束
		if err == io.EOF {
			fmt.Println("end!", err)
			break
		}
		if err != nil {
			fmt.Errorf("ListStr get stream err: %v", err)
		}
		// 打印返回值
		log.Println(res.StreamValue)
	}

	return 0
}
