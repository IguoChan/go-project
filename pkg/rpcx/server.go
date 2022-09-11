package rpcx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/IguoChan/go-project/pkg/etcdx"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	runtime2 "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

const (
	pprofUrlPrefix = "/debug/pprof"
)

type Server struct {
	gs  *grpc.Server
	opt *ServerOptions
	grpclog.LoggerV2
	er *etcdx.Register
}

type Gateway struct {
	*Server
	mux *http.ServeMux //ServeMux
	hs  *http.Server
}

type ServerOptions struct {
	EtcdOpt     *etcdx.Options `json:"-"`
	Host        string         `json:"host" mapstructure:"host"`
	Port        int            `json:"port" mapstructure:"port"`
	Gateway     *GWOptions     `json:"gateway" mapstructure:"gateway"`
	ServiceName string         `json:"name" mapstructure:"name"`

	// logger
	LogrusLogger *logrus.Logger `json:"-"`
}

type GWOptions struct {
	Port int `json:"port" mapstructure:"port"`
}

// NewServer r&rs for at least one demo_app register
func NewServer(opt *ServerOptions, r PBServerRegister, rs ...PBServerRegister) (*Server, error) {
	if opt == nil || opt.Host == "" || opt.Port == 0 {
		return nil, errors.New("options is nil, or addr is nil")
	}

	// new server
	server, err := newGrpc()
	if err != nil {
		return nil, err
	}

	// 注册etcd服务
	er, err := etcdx.NewRegister(opt.EtcdOpt)
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = er.Registry(ctx, opt.ServiceName, opt.Addr())
	if err != nil {
		return nil, err
	}

	// 注册pb服务
	r.RegisterPBServer(server)
	for _, ra := range rs {
		ra.RegisterPBServer(server)
	}

	logger := grpclog.NewLoggerV2(logrus.StandardLogger().Out, logrus.StandardLogger().Out, logrus.StandardLogger().Out)
	if opt.LogrusLogger != nil {
		logger = grpclog.NewLoggerV2(opt.LogrusLogger.Out, opt.LogrusLogger.Out, opt.LogrusLogger.Out)
	}
	grpclog.SetLoggerV2(logger)

	return &Server{
		gs:       server,
		opt:      opt,
		LoggerV2: logger,
		er:       er,
	}, nil
}

func (s *Server) Serve() error {
	defer s.gs.Stop()

	// 监听本地端口
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.opt.Port))
	if err != nil {
		logrus.Errorf("Net lister err: %v.", err)
		return err
	}
	s.LoggerV2.Infof("grpc start listening...:%d", s.opt.Port)

	// 用服务器Serve() 方法以及端口信息区实现阻塞等待，直到进程被杀死或者Stop() 被调用
	if err := s.gs.Serve(listen); err != nil {
		s.LoggerV2.Errorf("grpc serve err: %v.", err)
		return err
	}

	return nil
}

func (s *Server) Stop() {
	_ = s.er.Revoke()
	s.gs.GracefulStop()
}

func newGrpc() (*grpc.Server, error) {
	opts := make([]grpc.ServerOption, 0, 2)

	// 一元请求
	usi := make([]grpc.UnaryServerInterceptor, 0, 2)
	usi = append(usi, grpc_recovery.UnaryServerInterceptor(), unaryServerInterceptor())
	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(usi...)))

	// 流式请求
	ssi := make([]grpc.StreamServerInterceptor, 0, 2)
	ssi = append(ssi, grpc_recovery.StreamServerInterceptor(), streamServerInterceptor())
	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(ssi...)))

	opts = append(opts)

	// 新建grpc服务实例
	server := grpc.NewServer(opts...)

	return server, nil
}

func NewGateway(opt *ServerOptions, gr PBGatewayRegister, grs ...PBGatewayRegister) (*Gateway, error) {
	if opt == nil || opt.Gateway == nil || opt.Gateway.Port == 0 {
		return nil, errors.New("options is nil, or gateway port is 0")
	}

	rs := make([]PBServerRegister, 0, len(grs))
	for _, gr := range grs {
		rs = append(rs, gr.(PBServerRegister))
	}
	server, err := NewServer(opt, gr, rs...)
	if err != nil {
		return nil, err
	}

	rMux := runtime2.NewServeMux(runtime2.WithIncomingHeaderMatcher(customMatcher))
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	addr := opt.Addr()
	gr.RegisterPBGateway(context.Background(), rMux, addr, opts)
	for _, gra := range grs {
		gra.RegisterPBGateway(context.Background(), rMux, addr, opts)
	}

	mux := http.NewServeMux()
	mux.Handle("/", rMux)
	mux.Handle(pprofUrlPrefix+"/", http.HandlerFunc(pprof.Index))
	mux.Handle(pprofUrlPrefix+"/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle(pprofUrlPrefix+"/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle(pprofUrlPrefix+"/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle(pprofUrlPrefix+"/trace", http.HandlerFunc(pprof.Trace))
	mux.Handle(pprofUrlPrefix+"/heap", pprof.Handler("heap"))
	mux.Handle(pprofUrlPrefix+"/goroutine", pprof.Handler("goroutine"))
	mux.Handle(pprofUrlPrefix+"/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle(pprofUrlPrefix+"/block", pprof.Handler("block"))

	return &Gateway{
		Server: server,
		mux:    mux,
	}, nil
}

// Serve include grpc.Serve
func (g *Gateway) Serve() error {
	errCh := make(chan error)
	go func() {
		err := g.Server.Serve()
		if err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(time.Second): // wait 1 second, just wait for error handling
	}

	g.LoggerV2.Infof("grpc gateway start listening...:%d", g.opt.Gateway.Port)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", g.opt.Gateway.Port),
		Handler: g.mux,
	}
	g.hs = server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		g.LoggerV2.Errorf("grpc gateway serve err: %v.", err)
		g.Server.Stop()
		return err
	}

	return nil
}

func (g *Gateway) Stop() {
	g.Server.Stop()
	_ = g.hs.Shutdown(context.Background())
	_ = g.hs.Close()
}

func (o *ServerOptions) Addr() string {
	return fmt.Sprintf("%s:%d", o.Host, o.Port)
}

func customMatcher(key string) (string, bool) {
	switch key {
	case "TraceId":
		return key, true
	case "Token":
		return key, true
	default:
		return runtime2.DefaultHeaderMatcher(key)
	}
}
