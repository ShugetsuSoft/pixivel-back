package source

import (
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
)

type Task struct {
	mq models.MessageQueue
	taskchaname string
}

func NewTaskGenerator(mq models.MessageQueue, taskchaname string) *Task {
	return &Task{
		mq,
		taskchaname,
	}
}

func (gen *Task) SendTask(task models.CrawlTask, priority uint8) error {
	bin, err := utils.MsgPack(task)
	if err != nil {
		return err
	}
	return gen.mq.Publish(gen.taskchaname, models.MQMessage{Data: bin, Priority: priority})
}
