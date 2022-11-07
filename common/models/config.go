package models

type Modes int

const (
	NormalMode Modes = iota
	ArchiveMode
)

type Config struct {
	Mongodb       MongoDBConfig       `yaml:"mongodb"`
	Rabbitmq      RabbitMQConfig      `yaml:"rabbitmq"`
	Elasticsearch ElasticSearchConfig `yaml:"elasticsearch"`
	Redis         RedisConfig         `yaml:"redis"`
	Neardb        NearDBConfig        `yaml:"neardb"`
	Responser     ResponserConfig     `yaml:"responser"`
	General       GeneralConfig       `yaml:"general"`
	Spider        SpiderConfig        `yaml:"spider"`
}

type ResponserConfig struct {
	Listen string `yaml:"listen"`
	Debug  bool   `yaml:"debug"`
	Mode   Modes  `yaml:"mode"`
}

type SpiderConfig struct {
	PixivToken      string `yaml:"pixiv-token"`
	CrawlingThreads int    `yaml:"crawling-threads"`
}

type GeneralConfig struct {
	SpiderRetry uint   `yaml:"spider-retry"`
	Prometheus  string `yaml:"prometheus-listen"`
	Loki        string `yaml:"loki-uri"`
}

type MongoDBConfig struct {
	URI string `yaml:"uri"`
}

type RabbitMQConfig struct {
	URI string `yaml:"uri"`
}

type ElasticSearchConfig struct {
	URI  string `yaml:"uri"`
	User string `yaml:"user"`
	Pass string `yaml:"password"`
}

type RedisConfig struct {
	CacheRedis   CacheRedisConfig   `yaml:"cache"`
	MessageRedis MessageRedisConfig `yaml:"message"`
}

type CacheRedisConfig struct {
	URI string `yaml:"uri"`
}

type MessageRedisConfig struct {
	URI string `yaml:"uri"`
}

type NearDBConfig struct {
	URI string `yaml:"uri"`
}
