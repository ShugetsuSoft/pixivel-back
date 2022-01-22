package task

import (
	"context"
	"fmt"
	"github.com/ShugetsuSoft/pixivel-back/common/database/tasktracer"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
)

type TaskGenerator struct {
	mq models.MessageQueue
	taskchaname string
	retrys uint
	tracer *tasktracer.TaskTracer
	tracerlistener *tasktracer.TaskListener
}

func NewTaskGenerator(mq models.MessageQueue, taskchaname string, retrys uint, tracer *tasktracer.TaskTracer) *TaskGenerator {
	return &TaskGenerator{
		mq,
		taskchaname,
		retrys,
		tracer,
		tracer.NewListener(context.Background()),
	}
}

func (gen *TaskGenerator) SendTask(task models.CrawlTask) error {
	fmt.Println(task)
	bin, err := utils.MsgPack(task)
	if err != nil {
		return err
	}
	return gen.mq.Publish(gen.taskchaname, models.MQMessage{Data: bin, Priority: 5})
}