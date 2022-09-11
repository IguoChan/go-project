package demo_app

import (
	"sync"

	"github.com/IguoChan/go-project/api/genproto/demo_app/server_streampb"

	"github.com/IguoChan/go-project/api/genproto/demo_app/simplepb"

	"github.com/IguoChan/go-project/pkg/rpcx"
)

type RpcClient struct {
	c *rpcx.RpcClient
}

var (
	rpcClient *RpcClient
	once      sync.Once
)

func NewRpcClient(opts map[string]*rpcx.ClientOptions) *RpcClient {
	once.Do(func() {
		rpcClient = &RpcClient{
			c: rpcx.NewRpcClient(opts),
		}
	})
	return rpcClient
}

func (c RpcClient) Simple() (simplepb.SimpleClient, error) {
	rc, err := c.c.Get(ServiceName)
	if err != nil {
		return nil, err
	}
	return simplepb.NewSimpleClient(rc.Conn()), nil
}

func (c RpcClient) SS() (server_streampb.StreamServerClient, error) {
	rc, err := c.c.Get(ServiceName)
	if err != nil {
		return nil, err
	}
	return server_streampb.NewStreamServerClient(rc.Conn()), nil
}
