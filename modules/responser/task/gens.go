package task

import (
	"context"
	"fmt"
	"time"

	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
)

type stringmap map[string]string

func (gen *TaskGenerator) IllustDetailTask(ctx context.Context, id uint64) error {
	tid, existed, err := gen.tracer.NewTaskGroup(stringmap{"iid": utils.Itoa(id)}, 1)
	if err != nil {
		return err
	}
	if !existed {
		task := models.CrawlTask{
			Group:      tid,
			Type:       models.CrawlIllustDetail,
			Params:     map[string]string{"id": utils.Itoa(id)},
			RetryCount: gen.retrys,
		}
		err = gen.SendTask(task)
		if err != nil {
			gen.tracer.RemoveTaskGroup(tid)
			return err
		}
	}
	err = gen.tracerlistener.WaitFor(ctx, tid, time.Minute)
	if err != nil {
		gen.tracer.RemoveTaskGroup(tid)
		return err
	}
	return nil
}

func (gen *TaskGenerator) UgoiraDetailTask(ctx context.Context, id uint64) error {
	tid, existed, err := gen.tracer.NewTaskGroup(stringmap{"ugid": utils.Itoa(id)}, 1)
	if err != nil {
		return err
	}
	fmt.Println(id)
	if !existed {
		task := models.CrawlTask{
			Group:      tid,
			Type:       models.CrawlUgoiraDetail,
			Params:     map[string]string{"id": utils.Itoa(id)},
			RetryCount: gen.retrys,
		}
		err = gen.SendTask(task)
		if err != nil {
			gen.tracer.RemoveTaskGroup(tid)
			return err
		}
	}
	err = gen.tracerlistener.WaitFor(ctx, tid, time.Minute)
	if err != nil {
		gen.tracer.RemoveTaskGroup(tid)
		return err
	}
	return nil
}

func (gen *TaskGenerator) UserDetailTask(ctx context.Context, id uint64) error {
	tid, existed, err := gen.tracer.NewTaskGroup(stringmap{"uid": utils.Itoa(id)}, 1)
	if err != nil {
		return err
	}
	fmt.Println(id)
	if !existed {
		task := models.CrawlTask{
			Group:      tid,
			Type:       models.CrawlUserDetail,
			Params:     map[string]string{"id": utils.Itoa(id)},
			RetryCount: gen.retrys,
		}
		err = gen.SendTask(task)
		if err != nil {
			gen.tracer.RemoveTaskGroup(tid)
			return err
		}
	}
	err = gen.tracerlistener.WaitFor(ctx, tid, time.Minute)
	if err != nil {
		gen.tracer.RemoveTaskGroup(tid)
		return err
	}
	return nil
}

func (gen *TaskGenerator) UserIllustsTask(ctx context.Context, id uint64) error {
	tid, existed, err := gen.tracer.NewTaskGroup(stringmap{"uiid": utils.Itoa(id)}, 1)
	if err != nil {
		return err
	}
	fmt.Println(id)
	if !existed {
		task := models.CrawlTask{
			Group:      tid,
			Type:       models.CrawlUserIllusts,
			Params:     map[string]string{"id": utils.Itoa(id)},
			RetryCount: gen.retrys,
		}
		err = gen.SendTask(task)
		if err != nil {
			gen.tracer.RemoveTaskGroup(tid)
			return err
		}
	}
	err = gen.tracerlistener.WaitFor(ctx, tid, time.Minute)
	if err != nil {
		gen.tracer.RemoveTaskGroup(tid)
		return err
	}
	return nil
}

func (gen *TaskGenerator) RankInitTask(ctx context.Context, mode string, date string, content string) error {
	tid, existed, err := gen.tracer.NewTaskGroup(stringmap{"rtmd": mode, "rtdt": date}, 1)
	if err != nil {
		return err
	}
	if !existed {
		task := models.CrawlTask{
			Group:      tid,
			Type:       models.CrawlRankIllusts,
			Params:     map[string]string{"mode": mode, "page": "1", "date": date, "content": content},
			RetryCount: gen.retrys,
		}
		err = gen.SendTask(task)
		if err != nil {
			gen.tracer.RemoveTaskGroup(tid)
			return err
		}
	}
	err = gen.tracerlistener.WaitFor(ctx, tid, 5*time.Minute)
	if err != nil {
		gen.tracer.RemoveTaskGroup(tid)
		return err
	}
	return nil
}
