package responser

import (
	"github.com/ShugetsuSoft/pixivel-back/common/database/drivers"
	"github.com/ShugetsuSoft/pixivel-back/common/database/operations"
	"github.com/ShugetsuSoft/pixivel-back/common/database/tasktracer"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type Responser struct {
	app  *gin.Engine
	addr string
}

func NewResponser(addr string, dbops *operations.DatabaseOperations, mq models.MessageQueue, taskchaname string, retrys uint, tracer *tasktracer.TaskTracer, redis *drivers.RedisPool, debug bool) *Responser {
	app := gin.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://pixivel.moe", "https://beta.pixivel.moe"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router := NewRouter(dbops, mq, taskchaname, retrys, tracer, redis, debug)
	router.mount(app)

	return &Responser{
		app,
		addr,
	}
}

func (res *Responser) Run() error {
	return res.app.Run(res.addr)
}
