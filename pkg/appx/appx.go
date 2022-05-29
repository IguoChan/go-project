package appx

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/IguoChan/go-project/pkg/grpcx"
)

type app struct {
	appName string
	ctx     context.Context
	cancel  context.CancelFunc

	grpcServers []grpcx.GrpcServer

	workers []Worker
	wmu     sync.Mutex
}

func New() *app {
	ctx, cancel := context.WithCancel(context.Background())
	return &app{
		ctx:    ctx,
		cancel: cancel,
	}
}

// AddWorker 添加work
func (a *app) AddWorker(w Worker) {
	a.wmu.Lock()
	defer a.wmu.Unlock()
	a.workers = append(a.workers, w)
}

func (a *app) AddGrpcServer(opt *grpcx.ServerOptions, rs []grpcx.PBServerRegister) error {
	if len(rs) == 0 {
		a.cancel()
		return errors.New("no register service")
	}

	server, err := grpcx.NewServer(opt, rs[0], rs[1:]...)
	if err != nil {
		a.cancel()
	}

	a.grpcServers = append(a.grpcServers, server)

	return nil
}

// AddGrpcGateway Gateway包含了启动GrpcServer，无需再调用 AddGrpcServer
func (a *app) AddGrpcGateway(opt *grpcx.ServerOptions, rs ...grpcx.PBGatewayRegister) error {
	if len(rs) == 0 {
		a.cancel()
		return errors.New("no register service")
	}

	server, err := grpcx.NewGateway(opt, rs[0], rs[1:]...)
	if err != nil {
		a.cancel()
	}

	a.grpcServers = append(a.grpcServers, server)

	return nil
}

func (a *app) Run() int {
	// task
	for _, w := range a.workers {
		worker := w
		go func() {
			err := worker.Start()
			if err != nil {
				logrus.Error("task start failed!")
				a.cancel()
			}
		}()
	}

	// server
	for _, server := range a.grpcServers {
		s := server
		go func() {
			err := s.Serve()
			if err != nil {
				a.cancel()
			}
		}()
	}

	// wait
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	for {
		select {
		case <-a.ctx.Done():
			logrus.Info("server exit by ctx cancel, bye!")

			// stop the server
			for _, server := range a.grpcServers {
				server.Stop()
			}

			// stop the worker
			for _, w := range a.workers {
				_ = w.Stop()
			}

			time.Sleep(3 * time.Second)
			return 0
		case <-interrupt:
			logrus.Info("server receive Ctrl+C!")
			a.cancel()
		}
	}

}
