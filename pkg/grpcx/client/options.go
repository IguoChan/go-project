package client

import (
	"github.com/IguoChan/go-project/pkg/etcdx"
	"github.com/IguoChan/go-project/pkg/grpcx/resolver"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type options struct {
	logger    *zap.Logger
	usis      []grpc.UnaryClientInterceptor  // 简单rpc拦截器
	ssis      []grpc.StreamClientInterceptor // 流式rpc拦截器
	discovery resolver.Discovery
}

type Option func(opts *options)

func defaultOptions() *options {
	return &options{}
}

func SetEtcdDiscovery(serviceName string, etcdOpts *etcdx.Options, resolverOpts ...resolver.Option) Option {
	return func(opts *options) {
		d, _ := resolver.NewEtcdResolver(serviceName, etcdOpts, resolverOpts...)
		opts.discovery = d
	}
}

func WithEtcdDiscovery(r *resolver.EtcdResolver) Option {
	return func(opts *options) {
		opts.discovery = r
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
		opts.usis = append(opts.usis, grpc_zap.UnaryClientInterceptor(logger))
		opts.ssis = append(opts.ssis, grpc_zap.StreamClientInterceptor(logger))
	}
}

func SetEndpointsDiscovery(endpoints []string) Option {
	return func(opts *options) {
		d := resolver.NewEndpointsResolver(endpoints)
		opts.discovery = d
	}
}
