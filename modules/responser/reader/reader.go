package reader

import (
	"github.com/ShugetsuSoft/pixivel-back/common/database/operations"
	"github.com/ShugetsuSoft/pixivel-back/common/database/tasktracer"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/modules/responser/task"
)

type Reader struct {
	dbops *operations.DatabaseOperations
	gen   *task.TaskGenerator
	//shops *operations.SearchOperations
}

func NewReader(dbops *operations.DatabaseOperations, mq models.MessageQueue, taskchaname string, retrys uint, tracer *tasktracer.TaskTracer) *Reader {
	gen := task.NewTaskGenerator(mq, taskchaname, retrys, tracer)
	return &Reader{
		dbops: dbops,
		gen:   gen,
	}
}
