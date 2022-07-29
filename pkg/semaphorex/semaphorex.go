package semaphorex

import (
	"context"
	"errors"
)

type SemaphoreType string

const (
	SemaphoreProcess SemaphoreType = "process"
	SemaphoreRedis   SemaphoreType = "redis"
)

var (
	ErrGetSem = errors.New("get semaphore failed")
)

type Semaphore interface {
	Acquire(ctx context.Context) error
	TryAcquire() bool
	Release()
}

func NewSemaphore(t SemaphoreType, n int64, opts ...Option) Semaphore {
	switch t {
	case SemaphoreProcess:
		return NewProcessSem(n)
	case SemaphoreRedis:
		defaultOpts := defaultOptions()
		for _, apply := range opts {
			apply(defaultOpts)
		}
		return NewRedisSem(n, defaultOpts.name, defaultOpts.rc, defaultOpts)
	default:
		return NewProcessSem(n)
	}
}
