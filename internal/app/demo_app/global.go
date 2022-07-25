package demo_app

import (
	"github.com/IguoChan/go-project/pkg/cache/redisx"
	"github.com/IguoChan/go-project/pkg/dbx/mysqlx"
	"github.com/IguoChan/go-project/pkg/httpx"
	"github.com/IguoChan/go-project/pkg/mqx"
)

type Global struct {
	Config  *Config
	DB      *mysqlx.Client
	RDB     *redisx.Client
	MQ      mqx.MQer
	HttpCli *httpx.Client
}

func NewGlobal(c *Config) (*Global, error) {
	g := &Global{Config: c}

	// mysql
	db, err := mysqlx.NewClient(c.MySql)
	if err != nil {
		return nil, err
	}
	g.DB = db

	// redis
	rdb, err := redisx.NewClient(c.Redis)
	if err != nil {
		return nil, err
	}
	g.RDB = rdb

	// mq

	return g, nil
}
