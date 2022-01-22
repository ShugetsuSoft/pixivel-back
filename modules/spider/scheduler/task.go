package scheduler

import (
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
)

type TaskQueue struct {
	queue    <-chan models.MQMessage
	mq       models.MessageQueue
	channame string
}

func NewTaskQueue(mq models.MessageQueue, channame string) *TaskQueue {
	taskchan, err := mq.Consume(channame)
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "SpiderSchedulerError"}, err.Error())
		log.Fatal(err)
	}
	return &TaskQueue{
		mq:       mq,
		queue:    taskchan,
		channame: channame,
	}
}

func (taskq *TaskQueue) GetTask() (*models.CrawlTask, uint64, uint8, error) {
	data, ok := <-taskq.queue
	if ok {
		var task models.CrawlTask
		err := utils.MsgUnpack(data.Data, &task)
		if err != nil {
			return nil, 0, 0, err
		}
		return &task, data.Tag, data.Priority, nil
	}
	return nil, 0, 0, nil
}

func (taskq *TaskQueue) Resend(task *models.CrawlTask, priority uint8) error {
	bin, err := utils.MsgPack(task)
	if err != nil {
		return err
	}
	return taskq.mq.Publish(taskq.channame, models.MQMessage{
		Data:     bin,
		Priority: priority,
	})
}

func (taskq *TaskQueue) Publish(name string, message models.MQMessage) error {
	return taskq.mq.PublishToExchange(name, message)
}

func (taskq *TaskQueue) Ack(ackid uint64) error {
	return taskq.mq.Ack(ackid)
}

func (taskq *TaskQueue) Reject(ackid uint64) error {
	return taskq.mq.Reject(ackid)
}
