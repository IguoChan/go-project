package client

import (
	"github.com/IguoChan/go-project/pkg/etcdx"
	"github.com/IguoChan/go-project/pkg/grpcx/resolver"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type options struct {
	logger    *zap.Logger
	ucis      []grpc.UnaryClientInterceptor  // 简单rpc拦截器
	scis      []grpc.StreamClientInterceptor // 流式rpc拦截器
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

func SetEndpointsDiscovery(endpoints []string) Option {
	return func(opts *options) {
		d := resolver.NewEndpointsResolver(endpoints)
		opts.discovery = d
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
		opts.ucis = append(opts.ucis, grpc_zap.UnaryClientInterceptor(logger))
		opts.scis = append(opts.scis, grpc_zap.StreamClientInterceptor(logger))
	}
}

func WithPrometheus(enableClientHandlingTimeHistogram bool) Option {
	return func(opts *options) {
		if enableClientHandlingTimeHistogram {
			grpc_prometheus.EnableClientHandlingTimeHistogram()
		}
		opts.ucis = append(opts.ucis, grpc_prometheus.UnaryClientInterceptor)
		opts.scis = append(opts.scis, grpc_prometheus.StreamClientInterceptor)
	}
}
