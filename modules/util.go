package modules

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"gopkg.in/yaml.v3"

	"github.com/ShugetsuSoft/pixivel-back/common/database/drivers"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/config"
)

func ReadConfig(path string) (*models.Config, error) {
	var config models.Config
	confile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer confile.Close()
	content, err := ioutil.ReadAll(confile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, &config)
	return &config, err
}

func NewDB(conf *models.Config, ctx context.Context) (*drivers.MongoDatabase, error) {
	db, err := drivers.NewMongoDatabase(ctx, conf.Mongodb.URI, config.DatabaseName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewMQ(conf *models.Config) (*drivers.RabbitMQ, error) {
	mq, err := drivers.NewRabbitMQ(conf.Rabbitmq.URI)
	if err != nil {
		return nil, err
	}
	mq.QueueDeclare(config.CrawlTaskQueue)
	mq.QueueDeclare(config.CrawlResponsesMongodb)
	mq.QueueDeclare(config.CrawlResponsesElastic)
	mq.ExchangeDeclare(config.CrawlResponsesExchange, "fanout")
	mq.QueueBindExchange(config.CrawlResponsesMongodb, config.CrawlResponsesMongodb, config.CrawlResponsesExchange)
	mq.QueueBindExchange(config.CrawlResponsesElastic, config.CrawlResponsesElastic, config.CrawlResponsesExchange)
	return mq, nil
}

func NewES(conf *models.Config, ctx context.Context) (*drivers.ElasticSearch, error) {
	es, err := drivers.NewElasticSearchClient(ctx, conf.Elasticsearch.URI, conf.Elasticsearch.User, conf.Elasticsearch.Pass)
	return es, err
}

func NewRD(redisuri string) *drivers.RedisPool {
	redis := drivers.NewRedisPool(redisuri)
	return redis
}

func NewNB(conf *models.Config, ctx context.Context) (*drivers.NearDB, error) {
	ndb, err := drivers.NewNearDB(conf.Neardb.URI)
	return ndb, err
}
