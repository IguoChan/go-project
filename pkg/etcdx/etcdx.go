package etcdx

import (
	"errors"
	v3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"time"
)

const (
	defaultScheme     = "scheme"
	defaultGrpcScheme = "grpc"
	defaultTTL        = 5
)

type ClientV3 struct {
	*v3.Client
	*Options
	*zap.Logger
}

type Options struct {
	Addrs       []string      // etcd主机地址
	DialTimeout time.Duration // 连接失败超时时间
	Username    string
	Password    string
	Scheme      string

	// logger
	*zap.Logger

	// for service register
	TTL int64 `json:"ttl"` // Lease TTL时间，单位：s；每次KeepAlive 续租频率为TTL/3
}

func NewClient(opt *Options) (*ClientV3, error) {
	if opt == nil {
		return nil, errors.New("options is nil")
	}

	c, err := v3.New(v3.Config{
		Endpoints:   opt.Addrs,
		DialTimeout: opt.DialTimeout,
		Username:    opt.Username,
		Password:    opt.Password,
		Logger:      opt.Logger,
	})
	if err != nil {
		return nil, err
	}

	return &ClientV3{
		Client:  c,
		Options: opt,
		Logger:  c.GetLogger(),
	}, nil
}
