package taskx

import (
	"time"

	"github.com/IguoChan/go-project/pkg/mqx"
)

type BGTasker interface {
	Run()
	Tick() <-chan time.Time
	Name() string
}

type MQTasker interface {
	Subscribe()
	UnSubscribe()
	MessageIn() <-chan *mqx.Msg
	Handle(msg *mqx.Msg)
	Name() string
}
