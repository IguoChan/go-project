package once_task

import (
	"time"

	"github.com/sirupsen/logrus"
)

type OnceTask struct {
	retry, cnt int
	timer      *time.Timer
}

func NewOnceTask(retry int) *OnceTask {
	return &OnceTask{
		retry: retry,
		timer: time.NewTimer(5 * time.Second),
	}
}

func (o *OnceTask) Run() error {
	o.cnt++
	logrus.Infof("do once work %d times", o.cnt)
	return nil
}

func (o *OnceTask) Reset() {
	o.timer.Reset(6 * time.Second)
}

func (o *OnceTask) Continue() bool {
	return o.retry > o.cnt
}

func (o *OnceTask) Tick() <-chan time.Time {
	return o.timer.C
}

func (o *OnceTask) Name() string {
	return "ONCE WORK"
}
