package main

import (
	"os"

	"github.com/IguoChan/go-project/internal/app/demo_app"
	"github.com/IguoChan/go-project/pkg/appx"
)

func main() {
	os.Exit(Run())
}

func Run() int {
	// config
	conf := demo_app.Conf()

	// new app
	app := appx.New(conf.ServiceName)

	// add worker
	app.AddWorker(demo_app.NewDemoWorker())

	// set server
	if err := app.SetGrpcGateway(conf.GrpcServer, demo_app.NewSimpleServer(), demo_app.NewServerStream()); err != nil {
		return appx.ErrGrpcGateway
	}

	return app.Run()
}
