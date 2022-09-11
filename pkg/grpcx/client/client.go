package client

import (
	"context"
	"time"

	"github.com/IguoChan/go-project/pkg/grpcx/resolver"
	"google.golang.org/grpc"
)

// Config client端必须的一些参数
type Config struct {
	target string
}

type Client struct {
	cc  *grpc.ClientConn
	cfg *Config
	opt *options
}

func NewClient(cfg *Config, opts ...Option) (*Client, error) {
	// options
	dsOpts := defaultOptions()
	for _, apply := range opts {
		apply(dsOpts)
	}
	d := dsOpts.discovery
	if d == nil {
		d = resolver.NewEmptyResolver(cfg.target)
	}

	// register resolver
	d.Register()

	// dial
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Second)
	cc, err := grpc.DialContext(ctx, d.Target(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return &Client{
		cc:  cc,
		cfg: cfg,
		opt: dsOpts,
	}, nil
}

func (c *Client) Conn() *grpc.ClientConn {
	return c.cc
}

func (c *Client) Close() {
	_ = c.cc.Close()
}
