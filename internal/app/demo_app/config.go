package demo_app

import (
	"github.com/IguoChan/go-project/config"
)

const (
	ServiceName string = "demo"
)

func Conf() *Config {
	conf := &Config{
		Common: config.DefaultCommon(),
	}
	conf.ServiceName = ServiceName

	err := config.NewConfig(&ConfigFS, conf)
	if err != nil {
		panic(err)
	}

	config.Print(conf)

	return conf
}

type Config struct {
	*config.Common
}
