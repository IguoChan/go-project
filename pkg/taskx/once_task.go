package taskx

import (
	"context"

	"github.com/sirupsen/logrus"
)

// OnceTask 一次性任务，可以选择失败后的重试次数或者一直执行
type OnceTask struct {
	rs []OnceTasker
}

func NewOnceTask(rs ...OnceTasker) *OnceTask {
	return &OnceTask{
		rs: rs,
	}
}

func (t *OnceTask) Run(ctx context.Context) {
	for i := range t.rs {
		go t.run(ctx, t.rs[i])
	}
}

func (t *OnceTask) run(ctx context.Context, r OnceTasker) {
	for {
		select {
		case <-ctx.Done():
			logrus.Infof("stop once work: %+v", r.Name())
			return
		case <-r.Tick():
			err := r.Run()
			if err != nil && r.Continue() {
				r.Reset()
			} else {
				return
			}
		}
	}
}
