package bg_task

import "time"

type Print struct {
}

func NewPrint() *Print {
	return &Print{}
}

func (c *Print) Run() {
	//TODO implement me
	panic("implement me")
}

func (c *Print) Tick() <-chan time.Time {
	//TODO implement me
	panic("implement me")
}

func (c *Print) Name() string {
	//TODO implement me
	panic("implement me")
}
