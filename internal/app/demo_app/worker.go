package demo_app

import (
	"context"

	"github.com/IguoChan/go-project/internal/app/demo_app/task/mq_task"

	"github.com/IguoChan/go-project/internal/app/demo_app/task/bg_task"
	"github.com/IguoChan/go-project/pkg/appx"
	"github.com/IguoChan/go-project/pkg/taskx"
)

type demoWorker struct {
	ctx    context.Context
	cancel context.CancelFunc
	bgTask *taskx.BGTask
	mqTask *taskx.MQTask
}

func NewDemoWorker() appx.Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &demoWorker{
		ctx:    ctx,
		cancel: cancel,
		bgTask: taskx.NewBGTask(ctx),
		mqTask: taskx.NewMQTask(ctx),
	}
}

func (w *demoWorker) Start() error {
	// register background tasks
	w.bgTask.Register(bg_task.NewCheckStatusTask())
	//w.bgTask.Register(bg_task.NewPrint())

	// register mq tasks
	w.mqTask.Register(mq_task.NewEventTask())

	// run background task
	w.bgTask.Run(w.ctx)

	// run mq task
	//w.mqTask.Run(w.ctx)

	return nil
}

func (w *demoWorker) Stop() error {
	w.cancel() // stop the task
	return nil
}
