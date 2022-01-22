package main

import (
	"context"
	"flag"
	"github.com/ShugetsuSoft/pixivel-back/common/database/operations"
	"github.com/ShugetsuSoft/pixivel-back/common/database/tasktracer"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/config"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	"github.com/ShugetsuSoft/pixivel-back/modules"
	"github.com/ShugetsuSoft/pixivel-back/modules/storer"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var configPath string
	var runType int
	flag.StringVar(&configPath, "config", "config.yaml", "Config File Path")
	flag.IntVar(&runType, "run", 0, "Store Type")
	flag.Parse()

	conf, err := modules.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	telemetry.RegisterStorer()
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

	messageRedis := modules.NewRD(conf.Redis.MessageRedis.URI)
	tracer := tasktracer.NewTaskTracer(messageRedis, config.TaskTracerChannel)
	ft := messageRedis.NewCuckooFilter(config.DatabaseName)

	switch runType {
	case 0:
		ope := operations.NewDatabaseOperations(ctx, db, ft, es, ndb)
		stor := storer.NewStorer(mq, config.CrawlTaskQueue, conf.General.SpiderRetry, ope, tracer)
		stor.StoreDB(config.CrawlResponsesMongodb)
	case 1:
		ope := operations.NewDatabaseOperations(ctx, db, ft, es, ndb)
		stor := storer.NewStorer(mq, config.CrawlTaskQueue, conf.General.SpiderRetry, ope, tracer)
		stor.StoreES(config.CrawlResponsesElastic)
	}
}
