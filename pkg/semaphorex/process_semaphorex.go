package semaphorex

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type ProcessSem struct {
	w *semaphore.Weighted
}

func NewProcessSem(n int64) *ProcessSem {
	return &ProcessSem{
		w: semaphore.NewWeighted(n),
	}
}

func (p *ProcessSem) Acquire(ctx context.Context) error {
	return p.w.Acquire(ctx, 1)
}

func (p *ProcessSem) TryAcquire() bool {
	return p.w.TryAcquire(1)
}

func (p *ProcessSem) Release() {
	p.w.Release(1)
}
