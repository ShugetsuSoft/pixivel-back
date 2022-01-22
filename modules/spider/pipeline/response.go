package pipeline

import (
	"encoding/json"
	"errors"
	"github.com/ShugetsuSoft/pixivel-back/common/convert"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	"github.com/ShugetsuSoft/pixivel-back/modules/spider/scheduler"
	"github.com/gocolly/colly"
)

func (pipe *Pipeline) Response(r *colly.Response, taskq *scheduler.TaskQueue) {
	t := r.Ctx.GetAny("Task")
	ack := r.Ctx.GetAny("Ack")
	priority := r.Ctx.GetAny("Priority")
	switch task := t.(type) {
	case *models.CrawlTask:
		switch ackid := ack.(type) {
		case uint64:
			switch priority := priority.(type) {
			case uint8:
				resp, err := pipe.Deal(task.Type, r.Body, task)
				if resp != nil || err == nil {
					pipe.Send(task.Group, task.Type, resp, priority, taskq)
					taskq.Ack(ackid)
				} else {
					taskq.Reject(ackid)
					if task.RetryCount > 0 {
						task.RetryCount -= 1
						taskq.Resend(task, priority)
					} else {
						telemetry.SpiderErrorTaskCount.Inc()
						pipe.Send(task.Group, models.CrawlError, err.Error(), priority, taskq)
					}
				}
				switch urihash := r.Ctx.GetAny("Uri").(type) {
				case uint64:
					pipe.storer.ClearVisited(urihash)
				}
			}
		}
	}
}

func (pipe *Pipeline) Error(r *colly.Response, taskq *scheduler.TaskQueue) {
	telemetry.Log(telemetry.Label{"pos": "SipderError"}, utils.StringOut(r.Body))
	t := r.Ctx.GetAny("Task")
	ack := r.Ctx.GetAny("Ack")
	switch task := t.(type) {
	case *models.CrawlTask:
		switch ackid := ack.(type) {
		case uint64:
			taskq.Reject(ackid)
			priority := r.Ctx.GetAny("Priority")
			switch priority := priority.(type) {
			case uint8:
				if task.RetryCount > 0 {
					task.RetryCount -= 1
					switch urihash := r.Ctx.GetAny("Uri").(type) {
					case uint64:
						pipe.storer.ClearVisited(urihash)
					}
					taskq.Resend(task, priority)
				} else {
					telemetry.SpiderErrorTaskCount.Inc()
					pipe.Send(task.Group, models.CrawlError, nil, priority, taskq)
				}
			}
		}
	}
}

func (pipe *Pipeline) Deal(tasktype uint, rawb []byte, task *models.CrawlTask) (interface{}, error) {
	switch tasktype {
	case models.CrawlIllustDetail:
		var raw models.IllustRawResponse
		err := json.Unmarshal(rawb, &raw)
		if err != nil {
			return nil, err
		}
		if raw.Error {
			return nil, errors.New(raw.Message)
		}
		return convert.IllustRaw2Illust(&raw.Body), nil
	case models.CrawlUserDetail:
		var raw models.UserRawResponse
		err := json.Unmarshal(rawb, &raw)
		if err != nil {
			return nil, err
		}
		if raw.Error {
			return nil, errors.New(raw.Message)
		}
		return convert.UserRaw2User(&raw.Body), err
	case models.CrawlUserIllusts:
		var raw models.UserIllustsRawResponse
		err := json.Unmarshal(rawb, &raw)
		if err != nil {
			return nil, err
		}
		if raw.Error {
			return nil, errors.New(raw.Message)
		}
		return convert.UserIllusts2UserIllustsResponse(utils.Atoi(task.Params["id"]), &raw.Body), nil
	case models.CrawlRankIllusts:
		var raw models.RankIllustsRawResponse
		err := json.Unmarshal(rawb, &raw)
		if err != nil {
			return nil, err
		}
		return convert.RankIllusts2RankIllustsResponse(&raw), nil
	case models.CrawlUgoiraDetail:
		var raw models.UgoiraRawResponse
		err := json.Unmarshal(rawb, &raw)
		if err != nil {
			return nil, err
		}
		if raw.Error {
			return nil, errors.New(raw.Message)
		}
		return convert.UgoiraRaw2Ugoira(&raw.Body, utils.Atoi(task.Params["id"])), nil
	}
	return nil, nil
}
