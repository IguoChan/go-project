package semaphorex

import (
	"github.com/IguoChan/go-project/pkg/cache/redisx"
	"github.com/go-redis/redis/v8"
)

type setOptions struct {
	name string
	rc   *redisx.Client
}

type Option func(*setOptions)

func defaultOptions() *setOptions {
	return &setOptions{
		name: "default_semaphore",
		rc: &redisx.Client{
			UniversalClient: redis.NewUniversalClient(&redis.UniversalOptions{
				Addrs: []string{":6379"},
			}),
		},
	}
}

func SetName(name string) Option {
	return func(options *setOptions) {
		if name == "" {
			return
		}
		options.name = name
	}
}

func SetRedisClient(rc *redisx.Client) Option {
	return func(options *setOptions) {
		if rc == nil {
			return
		}
		options.rc = rc
	}
}
