package rpcx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IguoChan/go-project/pkg/etcdx"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

type Client struct {
	*grpc.ClientConn
	opt *ClientOptions
	grpclog.LoggerV2
}

type ClientOptions struct {
	EtcdOpt     *etcdx.Options
	ServiceName string

	// logger
	LogrusLogger *logrus.Logger
}

func NewClient(opt *ClientOptions) (*Client, error) {
	if opt == nil {
		return nil, errors.New("options is nil")
	}

	d, err := etcdx.NewDiscovery(opt.EtcdOpt)
	if err != nil {
		return nil, err
	}

	resolver.Register(d)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	dialUrl := fmt.Sprintf("%s://%s", d.Scheme(), opt.ServiceName)
	conn, err := grpc.DialContext(ctx, dialUrl, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	logger := grpclog.NewLoggerV2(logrus.StandardLogger().Out, logrus.StandardLogger().Out, logrus.StandardLogger().Out)
	if opt.LogrusLogger != nil {
		logger = grpclog.NewLoggerV2(opt.LogrusLogger.Out, opt.LogrusLogger.Out, opt.LogrusLogger.Out)
	}
	grpclog.SetLoggerV2(logger)

	return &Client{
		ClientConn: conn,
		opt:        opt,
		LoggerV2:   logger,
	}, nil
}

func (c *Client) Conn() *grpc.ClientConn {
	return c.ClientConn
}
