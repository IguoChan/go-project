package mq_task

import "github.com/IguoChan/go-project/pkg/mqx"

type EventTask struct {
}

func NewEventTask() *EventTask {
	return &EventTask{}
}

func (e *EventTask) Subscribe() {
	//TODO implement me
	panic("implement me")
}

func (e *EventTask) UnSubscribe() {
	//TODO implement me
	panic("implement me")
}

func (e *EventTask) MessageIn() <-chan *mqx.Msg {
	//TODO implement me
	panic("implement me")
}

func (e *EventTask) Handle(msg *mqx.Msg) {
	//TODO implement me
	panic("implement me")
}

func (e *EventTask) Name() string {
	//TODO implement me
	panic("implement me")
}
