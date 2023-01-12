package scheduler

import (
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/config"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	"github.com/ShugetsuSoft/pixivel-back/modules/spider/apis"
	"github.com/ShugetsuSoft/pixivel-back/modules/spider/storage"
	"github.com/gocolly/colly"
	"log"
	"net/http"
)

type Scheduler struct {
	storer       *storage.BetterInMemoryStorage
	cookie       *http.Cookie
	exchangename string
}

func NewScheduler(cookie *http.Cookie, storer *storage.BetterInMemoryStorage, exchangename string) *Scheduler {
	return &Scheduler{
		storer:       storer,
		cookie:       cookie,
		exchangename: exchangename,
	}
}

func (sch *Scheduler) Schedule(c *colly.Collector, taskq *TaskQueue) error {
	for {
		newTask, tag, priority, err := taskq.GetTask()
		if err != nil {
			return err
		}
		if newTask == nil {
			return models.ErrorChannelClosed
		}
		uri, needlogin := ConstructRequest(newTask)
		uhash := storage.GetUrlHash(uri)
		c.DisableCookies()
		ctx := colly.NewContext()
		ctx.Put("Task", newTask)
		ctx.Put("Priority", priority)
		ctx.Put("Ack", tag)
		ctx.Put("Uri", uhash)
		header := http.Header{
			"User-Agent": []string{config.UserAgent},
		}
		if needlogin {
			header.Add("Cookie", sch.cookie.String())
		}
		if uri != "" {
			telemetry.SpiderTaskCount.Inc()
			err = c.Request("GET", uri, nil, ctx, header)
			if err != nil {
				taskq.Reject(tag)
				telemetry.Log(telemetry.Label{"pos": "SpiderScheduler"}, err.Error())
				if err != colly.ErrAlreadyVisited {
					if newTask.RetryCount > 0 {
						newTask.RetryCount -= 1
						sch.storer.ClearVisited(uhash)
						taskq.Resend(newTask, priority)
					}
				} else {
					res := models.CrawlResponse{
						Group: newTask.Group,
						Type:  models.CrawlError,
						Response: &models.CrawlErrorResponse{
							TaskType: newTask.Type,
							Params:   newTask.Params,
							Message:  "Error Visited",
						},
					}
					binRes, err := utils.MsgPack(res)
					if err == nil {
						err = taskq.Publish(sch.exchangename, models.MQMessage{
							Data:     binRes,
							Priority: priority,
						})
						if err != nil {
							telemetry.Log(telemetry.Label{"pos": "SpiderScheduler"}, err.Error())
							log.Fatal(err)
						}
					} else {
						telemetry.Log(telemetry.Label{"pos": "SpiderScheduler"}, err.Error())
						log.Fatal(err)
					}
				}
			}
		}
	}
}

func ConstructRequest(task *models.CrawlTask) (string, bool) {
	switch task.Type {
	case models.CrawlIllustDetail:
		isLogin := false
		if v, ok := task.Params["login"]; ok && v == "1" {
			isLogin = true
		}
		return apis.IllustDetailG(task.Params["id"]), isLogin
	case models.CrawlUserDetail:
		return apis.UserDetailG(task.Params["id"]), false
	case models.CrawlUserIllusts:
		return apis.UserIllustsG(task.Params["id"]), true
	case models.CrawlRankIllusts:
		return apis.RankIllustsG(task.Params["mode"], task.Params["page"], task.Params["date"], task.Params["content"]), false
	case models.CrawlUgoiraDetail:
		return apis.UgoiraDetailG(task.Params["id"]), true
	}
	return "", false
}
