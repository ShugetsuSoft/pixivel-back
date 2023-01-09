package models

const (
	CrawlIllustDetail uint = iota
	CrawlUserDetail
	CrawlUserIllusts
	CrawlRankIllusts
	CrawlUgoiraDetail

	CrawlError
)

type IDList []uint64

type IDWithPos struct {
	ID  uint64
	Pos uint
}

type UserIllustsResponse struct {
	UserID  uint64
	Illusts IDList
}

type RankIllustsResponseMessage struct {
	Mode    string
	Date    string
	Content string
	Page    int
	Next    bool
	Illusts []IDWithPos
}

type MQMessage struct {
	Data     []byte
	Tag      uint64
	Priority uint8
}

type CrawlTask struct {
	Group      string
	Type       uint
	Params     map[string]string
	RetryCount uint
}

type CrawlResponse struct {
	Group    string
	Type     uint
	Response interface{}
}

type CrawlErrorResponse struct {
	TaskType uint
	Params   map[string]string
	Message  string
}
