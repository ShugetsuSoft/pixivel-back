package main

import (
	"flag"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/config"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	"github.com/ShugetsuSoft/pixivel-back/modules"
	"github.com/ShugetsuSoft/pixivel-back/modules/spider"
	"log"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "Config File Path")
	flag.Parse()

	conf, err := modules.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	telemetry.RegisterSpider()
	errchan := telemetry.RunLoki(conf.General.Loki, "spider")
	go func() {
		for err := range errchan {
			log.Println(err)
		}
	}()
	go telemetry.RunPrometheus(conf.General.Prometheus)

	mq, err := modules.NewMQ(conf)
	if err != nil {
		log.Fatal(err)
	}

	spiderI, err := spider.NewSpider(mq, config.CrawlTaskQueue, config.CrawlResponsesExchange, conf.Spider.PixivToken, conf.Spider.CrawlingThreads)
	if err != nil {
		log.Fatal(err)
	}

	spiderI.Crawl()
}
