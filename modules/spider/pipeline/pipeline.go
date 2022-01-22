package pipeline

import (
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	"github.com/ShugetsuSoft/pixivel-back/modules/spider/scheduler"
	"github.com/ShugetsuSoft/pixivel-back/modules/spider/storage"
	"log"
)

type Pipeline struct {
	exchangename string
	storer       *storage.BetterInMemoryStorage
}

func NewPipeline(exchangename string, storer *storage.BetterInMemoryStorage) *Pipeline {
	return &Pipeline{
		exchangename: exchangename,
		storer:       storer,
	}
}

func (pipe *Pipeline) Send(tid string, typeid uint, data interface{}, priority uint8, taskq *scheduler.TaskQueue) {
	res := models.CrawlResponse{
		Group:    tid,
		Type:     typeid,
		Response: data,
	}
	binres, err := utils.MsgPack(res)
	if err == nil {
		err = taskq.Publish(pipe.exchangename, models.MQMessage{
			Data:     binres,
			Priority: priority,
		})
		if err != nil {
			telemetry.Log(telemetry.Label{"pos": "SpiderPipeline"}, err.Error())
			log.Fatal(err)
		}
	} else {
		telemetry.Log(telemetry.Label{"pos": "SpiderPipeline"}, err.Error())
		log.Fatal(err)
	}
}
