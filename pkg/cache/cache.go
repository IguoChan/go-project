package cache

import (
	"context"

	"github.com/IguoChan/go-project/pkg/cache/bigcachex"
	"github.com/IguoChan/go-project/pkg/cache/redisx"
)

type Type int

const (
	BigCache Type = iota
	RedisCache
)

type Cacher interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type Options struct {
	Type
	BigOpt *bigcachex.Options
	ReOpt  *redisx.Options
}

func NewCacher(opt *Options) (Cacher, error) {
	switch opt.Type {
	case BigCache:
		return bigcachex.NewClient(opt.BigOpt)
	default:
		return redisx.NewClient(opt.ReOpt)
	}
}
