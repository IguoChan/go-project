package bg_task

import (
	"time"

	"github.com/sirupsen/logrus"
)

type CheckStatusTask struct {
}

func NewCheckStatusTask() *CheckStatusTask {
	return &CheckStatusTask{}
}

func (c *CheckStatusTask) Run() {
	logrus.Info(c.Name())
}

func (c *CheckStatusTask) Tick() <-chan time.Time {
	return time.NewTicker(5 * time.Second).C
}

func (c *CheckStatusTask) Name() string {
	return "CheckStatusTask"
}
