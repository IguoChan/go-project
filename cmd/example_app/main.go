package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/IguoChan/go-project/internal/app/example_app"
	"github.com/IguoChan/go-project/pkg/appx"
	"github.com/IguoChan/go-project/pkg/etcdx"
	"github.com/IguoChan/go-project/pkg/grpcx"
)

func main() {
	os.Exit(Run())
}

func Run() int {
	app := appx.New()

	// add server
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	host := ""
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && strings.Contains(ipnet.IP.String(), "192.168.0") {
				fmt.Println(ipnet.IP.String())
				host = ipnet.IP.String()
			}
		}
	}
	opt := &grpcx.ServerOptions{
		EtcdOpt: &etcdx.Options{
			Addrs:  []string{"192.168.0.102:2379"},
			TTL:    30,
			Scheme: "grpc",
		},
		Host:         host,
		Port:         9413,
		GWPort:       9414,
		ServiceName:  "cyg_service1",
		LogrusLogger: nil,
	}
	if err := app.AddGrpcGateway(opt, example_app.NewExampleServer()); err != nil {
		panic(err)
	}

	return app.Run()
}
