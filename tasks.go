package main

import (
	"context"
	"flag"
	"github.com/ShugetsuSoft/pixivel-back/common/database/operations"
	"github.com/ShugetsuSoft/pixivel-back/common/database/tasktracer"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/config"
	"github.com/ShugetsuSoft/pixivel-back/modules"
	"github.com/ShugetsuSoft/pixivel-back/modules/responser/task"
	"github.com/ShugetsuSoft/pixivel-back/tools/incrementor"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var configPath string
	var date string
	flag.StringVar(&configPath, "config", "config.yaml", "Config File Path")
	flag.StringVar(&date, "date", "", "Specify rank date")

	flag.Parse()

	conf, err := modules.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

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

	ft := messageRedis.NewBloomFilter(config.DatabaseName)
	ope := operations.NewDatabaseOperations(ctx, db, ft, es, ndb)

	taskgen := task.NewTaskGenerator(mq, config.CrawlTaskQueue, conf.General.SpiderRetry, tracer)
	var dateobj time.Time
	if date != "" {
		dateobj = time.Now().AddDate(0, 0, -2)
	} else {
		dateobj, err = time.Parse("20060102", date)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = incrementor.CrawlRank(taskgen, ope, dateobj)
	if err != nil {
		log.Fatal(err)
	}
}
