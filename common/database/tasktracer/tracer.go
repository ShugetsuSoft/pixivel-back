package tasktracer

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ShugetsuSoft/pixivel-back/common/database/drivers"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
)

type TaskTracer struct {
	redis      *drivers.RedisPool
	tracerchan string
}

type TaskListener struct {
	cond   *sync.Cond
	finish string
	cancle context.CancelFunc
	redis  *drivers.RedisClient
}

func NewTaskTracer(redis *drivers.RedisPool, tracerchan string) *TaskTracer {
	return &TaskTracer{
		redis,
		tracerchan,
	}
}

func (tc *TaskTracer) NewTaskGroup(params map[string]string, tasknum uint) (string, bool, error) {
	tid := utils.HashMap(params)
	cli := tc.redis.NewRedisClient()
	defer cli.Close()
	isexist, err := cli.KeyExist(tid)
	if err != nil {
		return "", false, err
	}
	if isexist {
		return tid, true, nil
	}
	err = cli.SetValueExpire(tid, utils.Itoa(tasknum), 60)
	if err != nil {
		return "", false, err
	}
	return tid, false, nil
}

func (tc *TaskTracer) NewTaskByNum(tid string, tasknum uint) error {
	if tid == "" {
		return nil
	}
	cli := tc.redis.NewRedisClient()
	defer cli.Close()
	_, err := cli.IncreaseBy(tid, tasknum)
	return err
}

func (tc *TaskTracer) NewTask(tid string) error {
	if tid == "" {
		return nil
	}
	cli := tc.redis.NewRedisClient()
	defer cli.Close()
	_, err := cli.Increase(tid)
	return err
}

func (tc *TaskTracer) RemoveTaskGroup(tid string) error {
	if tid == "" {
		return nil
	}
	cli := tc.redis.NewRedisClient()
	defer cli.Close()
	return cli.DeleteValue(tid)
}

func (tc *TaskTracer) FinishTask(tid string) bool {
	if tid == "" {
		return true
	}
	cli := tc.redis.NewRedisClient()
	defer cli.Close()
	num, err := cli.Decrease(tid)
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "tasktracer"}, err.Error())
		log.Fatal(err)
	}
	if num <= 0 {
		err = tc.RemoveTaskGroup(tid)
		if err != nil {
			telemetry.Log(telemetry.Label{"pos": "tasktracer"}, err.Error())
			log.Fatal(err)
		}
		err = cli.Publish(tc.tracerchan, tid)
		if err != nil {
			telemetry.Log(telemetry.Label{"pos": "tasktracer"}, err.Error())
			log.Fatal(err)
		}
		return true
	}
	return false
}

func (tc *TaskTracer) FailTask(tid string, info string) bool {
	if tid == "" {
		return true
	}
	err := tc.RemoveTaskGroup(tid)
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "tasktracer"}, err.Error())
		log.Fatal(err)
	}
	cli := tc.redis.NewRedisClient()
	defer cli.Close()
	err = cli.Publish(tc.tracerchan, tid+"-ERR-"+info)
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "tasktracer"}, err.Error())
		log.Fatal(err)
	}
	return true
}

func (tc *TaskTracer) NewListener(ctx context.Context) *TaskListener {
	ctx, cancel := context.WithCancel(ctx)
	cond := sync.NewCond(&sync.Mutex{})
	cli := tc.redis.NewRedisClient()
	tracerchan, err := cli.Subscribe(ctx, tc.tracerchan)
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "tasktracer"}, err.Error())
		log.Fatal(err)
	}
	tasklistener := &TaskListener{
		cond,
		"",
		cancel,
		cli,
	}
	go func() {
		for {
			tid, ok := <-tracerchan
			if !ok {
				return
			}
			tasklistener.cond.L.Lock()
			tasklistener.finish = tid
			tasklistener.cond.L.Unlock()
			tasklistener.cond.Broadcast()
		}
	}()
	return tasklistener
}

func (tl *TaskListener) CloseChan() {
	tl.cancle()
	tl.redis.Close()
}

func (tl *TaskListener) WaitFor(ctx context.Context, tid string) error {
	timeout := time.After(1 * time.Minute)
	signal := make(chan error, 1)
	go func() {
		tl.cond.L.Lock()
		for {
			if tl.finish == tid {
				tl.cond.L.Unlock()
				signal <- nil
				return
			} else if strings.Contains(tl.finish, tid) {
				tl.cond.L.Unlock()
				signal <- errors.New(tl.finish)
				return
			}
			tl.cond.Wait()
		}
	}()
	select {
	case err := <-signal:
		return err
	case <-timeout:
		return models.ErrorTimeOut
	case <-ctx.Done():
		return models.ErrorTimeOut
	}
}
