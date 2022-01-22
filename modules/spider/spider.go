package spider

import (
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	storage2 "github.com/ShugetsuSoft/pixivel-back/modules/spider/storage"
	"log"
	"net/http"
	"time"

	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/modules/spider/pipeline"
	"github.com/ShugetsuSoft/pixivel-back/modules/spider/scheduler"
	"github.com/gocolly/colly"
)

type Spider struct {
	col     *colly.Collector
	sche    *scheduler.Scheduler
	pipe    *pipeline.Pipeline
	mq      models.MessageQueue
	inqueue string
}

func NewSpider(mq models.MessageQueue, inqueue string, outexchange string, loginss string, threads int) (*Spider, error) {
	col := colly.NewCollector(
		colly.Async(true),
	)
	col.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: threads,
		RandomDelay: 600 * time.Millisecond,
	})

	storage := &storage2.BetterInMemoryStorage{}
	if err := col.SetStorage(storage); err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:  "PHPSESSID",
		Value: loginss,
	}
	return &Spider{
		mq:      mq,
		inqueue: inqueue,
		col:     col,
		sche:    scheduler.NewScheduler(cookie, storage),
		pipe:    pipeline.NewPipeline(outexchange, storage),
	}, nil
}

func (s *Spider) Crawl() {
BEGIN:
	taskq := scheduler.NewTaskQueue(s.mq, s.inqueue)
	s.pipe.Hook(s.col, taskq)
	err := s.sche.Schedule(s.col, taskq)
	if err != nil {
		if err == models.ErrorChannelClosed {
			goto BEGIN
		}
		telemetry.Log(telemetry.Label{"pos": "Spider"}, err.Error())
		log.Fatal(err)
	}
}
