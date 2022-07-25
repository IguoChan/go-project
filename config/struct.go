package config

import (
	"github.com/IguoChan/go-project/pkg/cache/redisx"
	"github.com/IguoChan/go-project/pkg/dbx/mysqlx"
	"github.com/IguoChan/go-project/pkg/etcdx"
	"github.com/IguoChan/go-project/pkg/grpcx"
)

type Config struct {
	*Common `json:"common" mapstructure:"common"`
}

type Common struct {
	ServiceName string                          `json:"service_name" mapstructure:"service_name"`
	Apps        map[string]*grpcx.ServerOptions `json:"apps"  mapstructure:"apps"`
	GrpcServer  *grpcx.ServerOptions            `json:"grpc_server" mapstructure:"grpc_server"`
	MySql       *mysqlx.Options                 `json:"mysql" mapstructure:"mysql"`
	Redis       *redisx.Options                 `json:"redis" mapstructure:"mysql"`
	Etcd        *etcdx.Options                  `json:"etcd" mapstructure:"etcd"`
}
