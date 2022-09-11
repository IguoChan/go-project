package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"

	"github.com/sirupsen/logrus"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
)

// Config server端必须的一些参数
type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Server struct {
	s   *grpc.Server
	cfg *Config
	opt *options
	*zap.SugaredLogger
}

func NewServer(cfg *Config, rs []PBServerRegister, opts ...Option) (*Server, error) {
	// options
	dsOpts := defaultOptions()
	for _, apply := range opts {
		apply(dsOpts)
	}

	// new server
	server, err := newGrpc(dsOpts)
	if err != nil {
		return nil, err
	}

	// register proto service
	for _, r := range rs {
		r.RegisterPBServer(server)
	}

	// 注册服务，譬如etcd
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = dsOpts.register.Registry(ctx, cfg.Addr())
	if err != nil {
		return nil, err
	}

	return &Server{
		s:             server,
		cfg:           cfg,
		opt:           dsOpts,
		SugaredLogger: dsOpts.logger.Sugar(),
	}, nil
}

func (s *Server) Serve() error {
	defer s.s.Stop()

	// 监听本地端口
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		logrus.Errorf("Net lister err: %v.", err)
		return err
	}
	s.Infof("grpc start listening...:%d", s.cfg.Port)

	// 用服务器Serve() 方法以及端口信息区实现阻塞等待，直到进程被杀死或者Stop() 被调用
	if err := s.s.Serve(listen); err != nil {
		s.Errorf("grpc serve err: %v.", err)
		return err
	}

	return nil
}

func newGrpc(dsOpts *options) (*grpc.Server, error) {

	// logger
	if dsOpts.logger != nil {
		grpc_zap.ReplaceGrpcLoggerV2(dsOpts.logger)
	}

	// ServerOption
	srvOpts := make([]grpc.ServerOption, 0, 2)
	// 一元请求
	srvOpts = append(srvOpts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(dsOpts.usis...)))
	// 流式请求
	srvOpts = append(srvOpts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(dsOpts.ssis...)))

	// 新建grpc服务实例
	server := grpc.NewServer(srvOpts...)

	return server, nil
}

func (s *Server) Stop() {
	_ = s.opt.register.UnRegistry
	s.s.GracefulStop()
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
