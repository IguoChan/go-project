package config

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/IguoChan/go-project/pkg/configx"

	"github.com/spf13/viper"
)

var (
	comm = &Common{}
	once = &sync.Once{}
)

func DefaultCommon() *Common {
	once.Do(func() {
		cfg := &Config{}
		err := loadFs(&ConfigFS, cfg)
		if err != nil {
			once = &sync.Once{}
		}
		comm = cfg.Common
	})
	return comm
}

func loadFs(fs *embed.FS, config any) error {
	v := viper.GetViper()
	v.SetConfigType("yml")

	bs, err := fs.ReadFile("config.yml")
	if err != nil {
		return err
	}

	err = v.ReadConfig(bytes.NewBuffer(bs))
	if err != nil {
		return err
	}

	v.AutomaticEnv()
	err = v.Unmarshal(config)
	if err != nil {
		return err
	}

	return nil
}

func (c *Common) Merge() {
	if c.ServiceName == "" {
		return
	}

	as, ok := c.Apps[c.ServiceName]
	if !ok {
		return
	}
	c.GrpcServer = as
	c.GrpcServer.EtcdOpt = c.Etcd
}

func NewConfig(fs *embed.FS, config any) error {
	err := loadFs(fs, config)
	if err != nil {
		return err
	}

	if c, ok := config.(configx.Config); ok {
		c.Merge()
	}

	return nil
}

func Print(config any) {
	b, _ := json.MarshalIndent(config, "", "  ")
	fmt.Println(string(b))
}
