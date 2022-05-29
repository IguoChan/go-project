package taskx

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type MQTask struct {
	ctx context.Context
	rs  []MQTasker
	mu  *sync.Mutex
}

func NewMQTask(ctx context.Context) *MQTask {
	return &MQTask{
		ctx: ctx,
		mu:  &sync.Mutex{},
		rs:  make([]MQTasker, 0),
	}
}

func (t *MQTask) Register(r MQTasker) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.rs = append(t.rs, r)
}

func (t *MQTask) Run(ctx context.Context) {
	for i := range t.rs {
		go t.run(ctx, t.rs[i])
	}
}

func (t *MQTask) run(ctx context.Context, r MQTasker) {
	r.Subscribe()
	for {
		select {
		case <-ctx.Done():
			// 取消订阅
			r.UnSubscribe()
			logrus.Infof("stop mq task: %+v", r.Name())
			return
		case msg := <-r.MessageIn():
			r.Handle(msg)
		}
	}
}
