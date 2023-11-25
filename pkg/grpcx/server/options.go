package server

import (
	"github.com/IguoChan/go-project/pkg/etcdx"
	"github.com/IguoChan/go-project/pkg/grpcx/resolver"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// ------------------------------------------- server options -------------------------------------------
type options struct {
	logger   *zap.Logger
	usis     []grpc.UnaryServerInterceptor  // 简单rpc拦截器
	ssis     []grpc.StreamServerInterceptor // 流式rpc拦截器
	register resolver.Register
}

type Option func(opts *options)

func defaultOptions() *options {
	logger, _ := zap.NewDevelopment()
	return &options{
		logger:   logger,
		register: resolver.NewEmbedResolver(resolver.SchemePassThrough, ""),
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
		opts.usis = append(opts.usis, grpc_zap.UnaryServerInterceptor(logger))
		opts.ssis = append(opts.ssis, grpc_zap.StreamServerInterceptor(logger))
	}
}

func WithRecovery() Option {
	return func(opts *options) {
		opts.usis = append(opts.usis, grpc_recovery.UnaryServerInterceptor())
		opts.ssis = append(opts.ssis, grpc_recovery.StreamServerInterceptor())
	}
}

func SetEtcdRegister(serviceName string, etcdOpts *etcdx.Options, resolverOpts ...resolver.Option) Option {
	return func(opts *options) {
		r, _ := resolver.NewEtcdResolver(serviceName, etcdOpts, resolverOpts...)
		opts.register = r
	}
}

func WithEtcdRegister(r *resolver.EtcdResolver) Option {
	return func(opts *options) {
		opts.register = r
	}
}

func WithPrometheus(enableHandlingTimeHistogram bool) Option {
	return func(opts *options) {
		if enableHandlingTimeHistogram {
			grpc_prometheus.EnableHandlingTimeHistogram()
		}
		opts.usis = append(opts.usis, grpc_prometheus.UnaryServerInterceptor)
		opts.ssis = append(opts.ssis, grpc_prometheus.StreamServerInterceptor)
	}
}
