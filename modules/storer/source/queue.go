package source

import (
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
)

type DataQueue struct {
	queue <-chan models.MQMessage
	mq    models.MessageQueue
}

func NewDataQueue(mq models.MessageQueue, channame string) *DataQueue {
	taskchan, err := mq.Consume(channame)
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "StorerQueue"}, err.Error())
		log.Fatal(err)
	}
	return &DataQueue{
		mq:    mq,
		queue: taskchan,
	}
}

func (dataq *DataQueue) GetData() (*models.CrawlResponse, uint64, uint8, error) {
	data, ok := <-dataq.queue
	if ok {
		var task models.CrawlResponse
		err := utils.MsgUnpack(data.Data, &task)
		if err != nil {
			return nil, 0, 0, err
		}
		return &task, data.Tag, data.Priority, nil
	}
	return nil, 0, 0, nil
}

func (dataq *DataQueue) Publish(name string, message models.MQMessage) error {
	return dataq.mq.Publish(name, message)
}

func (dataq *DataQueue) Ack(ackid uint64) error {
	return dataq.mq.Ack(ackid)
}

func (dataq *DataQueue) Reject(ackid uint64) error {
	return dataq.mq.Reject(ackid)
}
