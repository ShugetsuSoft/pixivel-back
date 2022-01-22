package storer

import (
	"github.com/ShugetsuSoft/pixivel-back/common/database/operations"
	"github.com/ShugetsuSoft/pixivel-back/common/database/tasktracer"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	"github.com/ShugetsuSoft/pixivel-back/modules/storer/source"
	"log"
)

type Storer struct {
	mq     models.MessageQueue
	ops    *operations.DatabaseOperations
	tracer *tasktracer.TaskTracer
	task   *source.Task
	retrys uint
}

func NewStorer(mq models.MessageQueue, taskchanname string, retrys uint, ops *operations.DatabaseOperations, tracer *tasktracer.TaskTracer) *Storer {
	return &Storer{
		mq:     mq,
		ops:    ops,
		tracer: tracer,
		task:   source.NewTaskGenerator(mq, taskchanname),
		retrys: retrys,
	}
}

func (st *Storer) StoreDB(dbchanname string) {
BEGIN:
	dataq := source.NewDataQueue(st.mq, dbchanname)
	var err error
	err = st.handleDatabase(dataq)
	if err != nil {
		if err == models.ErrorChannelClosed {
			goto BEGIN
		}
		telemetry.Log(telemetry.Label{"pos": "Storer"}, err.Error())
		log.Fatal(err)
	}
}

func (st *Storer) StoreES(eschanname string) {
BEGIN:
	dataq := source.NewDataQueue(st.mq, eschanname)
	var err error
	err = st.handleElasticSearch(dataq)
	if err != nil {
		if err == models.ErrorChannelClosed {
			goto BEGIN
		}
		telemetry.Log(telemetry.Label{"pos": "Storer"}, err.Error())
		log.Fatal(err)
	}
}
