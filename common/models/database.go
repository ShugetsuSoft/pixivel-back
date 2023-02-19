package models

import "time"

type UserImage struct {
	Url        string `bson:"url"`
	BigUrl     string `bson:"bigUrl"`
	Background string `bson:"background"`
}

type User struct {
	ID         uint64    `bson:"_id"`
	UpdateTime time.Time `bson:"update_time"`

	Name              string    `bson:"name"`
	Bio               string    `bson:"bio,omitempty"`
	Image             UserImage `bson:"image"`
	IllustsUpdateTime time.Time `bson:"illusts_update_time"`
	IllustsCount      uint      `bson:"illusts_count"`

	Banned bool `bson:"banned"`
}

type IllustTag struct {
	Name        string `bson:"name" json:"name"`
	Translation string `bson:"translation,omitempty" json:"translation,omitempty"`
}

type IllustStatistic struct {
	Bookmarks uint `bson:"bookmarks" json:"bookmarks"`
	Likes     uint `bson:"likes" json:"likes"`
	Comments  uint `bson:"comments" json:"comments"`
	Views     uint `bson:"views" json:"views"`
}

type Illust struct {
	ID         uint64    `bson:"_id"`
	UpdateTime time.Time `bson:"update_time"`

	Title       string          `bson:"title"`
	AltTitle    string          `bson:"altTitle"`
	Description string          `bson:"description,omitempty"`
	Type        uint            `bson:"type"`
	CreateDate  time.Time       `bson:"createDate"`
	UploadDate  time.Time       `bson:"uploadDate"`
	Sanity      uint            `bson:"sanity"`
	Width       uint            `bson:"width"`
	Height      uint            `bson:"height"`
	PageCount   uint            `bson:"pageCount"`
	Tags        []IllustTag     `bson:"tags"`
	Statistic   IllustStatistic `bson:"statistic"`
	User        uint            `bson:"user"`
	Image       time.Time       `bson:"image"`
	AIType      uint            `bson:"aiType,omitempty"`

	Popularity uint `bson:"popularity"`
	Banned     bool `bson:"banned"`
}

type RankIllust struct {
	ID   uint64 `bson:"id"`
	Rank uint   `bson:"rank"`
}

type RankAggregateResult struct {
	Content Illust `bson:"content"`
}

type Rank struct {
	Date    string       `bson:"date"`
	Mode    string       `bson:"mode"`
	Content string       `bson:"content"`
	Illusts []RankIllust `bson:"illusts"`
}

type Ugoira struct {
	ID         uint64    `bson:"_id"`
	UpdateTime time.Time `bson:"update_time"`

	Image    time.Time     `bson:"image"`
	MimeType string        `bson:"mimeType"`
	Frames   []UgoiraFrame `bson:"frames"`
}

type UgoiraFrame struct {
	File  string `bson:"file"`
	Delay int    `bson:"delay"`
}

type IllustSearch struct {
	Title       string      `json:"title"`
	AltTitle    string      `json:"alt_title,omitempty"`
	Description string      `json:"description,omitempty"`
	Type        uint        `json:"type"`
	CreateDate  time.Time   `json:"create_date"`
	Sanity      uint        `json:"sanity"`
	Popularity  uint        `json:"popularity"`
	User        uint        `json:"user"`
	Tags        []IllustTag `json:"tags"`
}

type UserSearch struct {
	Name string `json:"name"`
	Bio  string `json:"bio,omitempty"`
}

const (
	IllustSearchMapping = `{"mappings":{"properties":{"alt_title":{"type":"text","analyzer":"kuromoji","search_analyzer":"kuromoji","fields":{"keyword":{"type":"keyword","ignore_above":256},"suggest":{"type": "completion","analyzer":"kuromoji"}}},"create_date":{"type":"date"},"description":{"type":"text","analyzer":"kuromoji","search_analyzer":"kuromoji"},"sanity":{"type":"short"},"popularity":{"type":"long"},"title":{"type":"text","analyzer":"kuromoji","search_analyzer":"kuromoji","fields":{"keyword":{"type":"keyword","ignore_above":256},"suggest":{"type": "completion","analyzer":"kuromoji"}}},"type":{"type":"short"},"user":{"type":"long"},"tags":{"type":"nested", "properties":{"name":{"type":"keyword","fields":{"suggest":{"type": "completion","analyzer":"kuromoji"},"fuzzy":{"type":"text","analyzer":"kuromoji","search_analyzer":"kuromoji"}}},"translation":{"type":"keyword","fields":{"suggest":{"type": "completion","analyzer":"smartcn"},"fuzzy":{"type":"text","analyzer":"smartcn","search_analyzer":"smartcn"}}}}}}}}`
	UserSearchMapping   = `{"mappings":{"properties":{"bio":{"type":"text","analyzer":"kuromoji","search_analyzer":"kuromoji","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"name":{"type":"text","analyzer":"kuromoji","search_analyzer":"kuromoji","fields":{"keyword":{"type":"keyword","ignore_above":256},"suggest":{"type": "completion","analyzer":"kuromoji"}}}}}}`
)
