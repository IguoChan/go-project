package client

import (
	"context"
	"fmt"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

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
	//balancer.NewVersionBuilder(&balancer.Options{})
	// options
	dsOpts := defaultOptions()
	for _, apply := range opts {
		apply(dsOpts)
	}
	d := dsOpts.discovery
	if d == nil {
		d = resolver.NewEmbedResolver(resolver.SchemePassThrough, cfg.target)
	}

	//policy := `{
	//	"loadBalancingConfig": [ { "round_robin": {} } ],
	//	"methodConfig": [{
	//	  "retryPolicy": {
	//		  "MaxAttempts": 4,
	//		  "InitialBackoff": ".01s",
	//		  "MaxBackoff": ".01s",
	//		  "BackoffMultiplier": 1.0,
	//		  "RetryableStatusCodes": [ "UNAVAILABLE" ]
	//	  }
	//	}]}`

	// dial
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Second)
	cc, err := grpc.DialContext(ctx, d.Target(),
		grpc.WithResolvers(d),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "weight")),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(dsOpts.ucis...)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(dsOpts.scis...)),
	)
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
