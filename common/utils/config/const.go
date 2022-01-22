package config

const (
	DatabaseName    = "Pixivel"
	IllustTableName = "Illust"
	UserTableName   = "User"
	RankTableName   = "Rank"
	UgoiraTableName = "Ugoira"

	IllustSearchIndexName = "illust"
	UserSearchIndexName   = "user"

	TaskTracerChannel      = "TaskTracer"
	CrawlTaskQueue         = "CrawlTasks"
	CrawlResponsesMongodb  = "CrawlResponsesMongodb"
	CrawlResponsesElastic  = "CrawlResponsesElastic"
	CrawlResponsesExchange = "CrawlResponses"

	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"

	RecommendDrift = 0.00001
)
