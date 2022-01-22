package main

import (
	"context"
	"flag"
	"github.com/ShugetsuSoft/pixivel-back/common/database/operations"
	"github.com/ShugetsuSoft/pixivel-back/common/database/tasktracer"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/config"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	"github.com/ShugetsuSoft/pixivel-back/modules"
	"github.com/ShugetsuSoft/pixivel-back/modules/responser"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "Config File Path")
	flag.Parse()

	conf, err := modules.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	telemetry.RegisterResponser()
	errchan := telemetry.RunLoki(conf.General.Loki, "responser")
	go func() {
		for err := range errchan {
			log.Println(err)
		}
	}()
	go telemetry.RunPrometheus(conf.General.Prometheus)

	db, err := modules.NewDB(conf, ctx)
	if err != nil {
		log.Fatal(err)
	}
	mq, err := modules.NewMQ(conf)
	if err != nil {
		log.Fatal(err)
	}
	ndb, err := modules.NewNB(conf, ctx)
	if err != nil {
		log.Fatal(err)
	}
	es, err := modules.NewES(conf, ctx)
	if err != nil {
		log.Fatal(err)
	}

	cacheRedis := modules.NewRD(conf.Redis.CacheRedis.URI)
	messageRedis := modules.NewRD(conf.Redis.MessageRedis.URI)

	ft := messageRedis.NewCuckooFilter(config.DatabaseName)

	ope := operations.NewDatabaseOperations(ctx, db, ft, es, ndb)

	tracer := tasktracer.NewTaskTracer(messageRedis, config.TaskTracerChannel)

	resp := responser.NewResponser(conf.Responser.Listen, ope, mq, config.CrawlTaskQueue, conf.General.SpiderRetry, tracer, cacheRedis, conf.Responser.Debug)

	err = resp.Run()
	if err != nil {
		log.Fatal(err)
	}
}
