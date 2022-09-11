package rpcx

import (
	"fmt"

	memo "github.com/IguoChan/go-project/pkg/memoiz"
)

type RpcClient struct {
	memo *memo.Memo
}

func NewRpcClient(opts map[string]*ClientOptions) *RpcClient {
	return &RpcClient{
		memo: memo.New(connFunc(opts)),
	}
}

func connFunc(opts map[string]*ClientOptions) memo.Func {
	return func(key string) (any, error) {
		opt, ok := opts[key]
		if !ok {
			return nil, fmt.Errorf("no key [%+v] in rpc client options: %+v", key, opts)
		}
		c, err := NewClient(opt)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}

func (c *RpcClient) Get(serviceName string) (*Client, error) {
	client, err := c.memo.Get(serviceName)
	if err != nil {
		c.memo.Delete(serviceName)
		return nil, err
	}
	return client.(*Client), nil
}
