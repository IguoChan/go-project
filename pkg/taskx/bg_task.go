package taskx

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type BGTask struct {
	ctx context.Context
	rs  []BGTasker
	mu  *sync.Mutex
}

func NewBGTask(ctx context.Context) *BGTask {
	return &BGTask{
		ctx: ctx,
		mu:  &sync.Mutex{},
		rs:  make([]BGTasker, 0),
	}
}

func (t *BGTask) Register(r BGTasker) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.rs = append(t.rs, r)
}

func (t *BGTask) Run(ctx context.Context) {
	for i := range t.rs {
		go t.run(ctx, t.rs[i])
	}
}

func (t *BGTask) run(ctx context.Context, r BGTasker) {
	for {
		select {
		case <-ctx.Done():
			logrus.Infof("stop bg work: %+v", r.Name())
			return
		case <-r.Tick():
			r.Run()
		}
	}
}
