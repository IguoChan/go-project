package appx

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IguoChan/go-project/pkg/grpcx"
	"github.com/sirupsen/logrus"
)

const (
	ErrServiceDependency = iota + 1
	ErrGrpcGateway
)

type app struct {
	serviceName string
	ctx         context.Context
	cancel      context.CancelFunc

	grpcServer grpcx.GrpcServer

	workers []Worker
	wmu     sync.Mutex
}

func New(serviceName string) *app {
	ctx, cancel := context.WithCancel(context.Background())
	return &app{
		serviceName: serviceName,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// AddWorker 添加work
func (a *app) AddWorker(w Worker) {
	a.wmu.Lock()
	defer a.wmu.Unlock()
	a.workers = append(a.workers, w)
}

func (a *app) SetGrpcServer(opt *grpcx.ServerOptions, rs ...grpcx.PBServerRegister) error {
	if len(rs) == 0 {
		a.cancel()
		return errors.New("no register service")
	}

	server, err := grpcx.NewServer(opt, rs[0], rs[1:]...)
	if err != nil {
		a.cancel()
		return err
	}

	a.grpcServer = server

	return nil
}

// SetGrpcGateway Gateway包含了启动GrpcServer，无需再调用 SetGrpcServer
func (a *app) SetGrpcGateway(opt *grpcx.ServerOptions, rs ...grpcx.PBGatewayRegister) error {
	if len(rs) == 0 {
		a.cancel()
		return errors.New("no register service")
	}

	server, err := grpcx.NewGateway(opt, rs[0], rs[1:]...)
	if err != nil {
		a.cancel()
		return err
	}

	a.grpcServer = server

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
	go func() {
		err := a.grpcServer.Serve()
		if err != nil {
			a.cancel()
		}
	}()

	// wait
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	for {
		select {
		case <-a.ctx.Done():
			logrus.Info("server exit by ctx cancel, bye!")

			// stop the server
			a.grpcServer.Stop()

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
