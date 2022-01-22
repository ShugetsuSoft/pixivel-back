package pipeline

import (
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	"github.com/ShugetsuSoft/pixivel-back/modules/spider/scheduler"
	"github.com/gocolly/colly"
)

func (pipe *Pipeline) Hook(c *colly.Collector, taskq *scheduler.TaskQueue) {
	c.OnResponse(func(r *colly.Response) {
		telemetry.Log(telemetry.Label{"pos": "SpiderCrawlURL"}, r.Request.URL.String())
		pipe.Response(r, taskq)
	})
	c.OnError(func(r *colly.Response, err error) {
		telemetry.Log(telemetry.Label{"pos": "SpiderCrawlURL"}, r.Request.URL.String())
		if err != nil {
			telemetry.Log(telemetry.Label{"pos": "SpiderCrawl"}, err.Error())
			pipe.Error(r, taskq)
		}
	})
}
