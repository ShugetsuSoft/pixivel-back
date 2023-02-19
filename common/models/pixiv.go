package models

import "time"

type IllustRawUrls struct {
	Mini     string `json:"mini"`
	Thumb    string `json:"thumb"`
	Small    string `json:"small"`
	Regular  string `json:"regular"`
	Original string `json:"original"`
}
type IllustRawTranslation struct {
	En string `json:"en"`
}
type IllustRawTag struct {
	Tag         string               `json:"tag"`
	Translation IllustRawTranslation `json:"translation"`
}

type IllustRawTagPre struct {
	Tags []IllustRawTag `json:"tags"`
}

type IllustRaw struct {
	ID          uint64          `json:"id,string"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	IllustType  uint            `json:"illustType"`
	CreateDate  time.Time       `json:"createDate"`
	UploadDate  time.Time       `json:"uploadDate"`
	Restrict    uint            `json:"restrict"`
	XRestrict   uint            `json:"xRestrict"`
	Sl          uint            `json:"sl"`
	Urls        IllustRawUrls   `json:"urls"`
	Tags        IllustRawTagPre `json:"tags"`
	Alt         string          `json:"alt"`
	UserID      uint            `json:"userId,string"`
	AIType      uint            `json:"aiType"`

	Width     uint `json:"width"`
	Height    uint `json:"height"`
	PageCount uint `json:"pageCount"`

	BookmarkCount uint `json:"bookmarkCount"`
	LikeCount     uint `json:"likeCount"`
	CommentCount  uint `json:"commentCount"`
	ViewCount     uint `json:"viewCount"`
}

type IllustRawResponse struct {
	Error   bool      `json:"error"`
	Message string    `json:"message"`
	Body    IllustRaw `json:"body"`
}

type UserRawBackground struct {
	URL string `json:"url"`
}

type UserRaw struct {
	UserID     uint64            `json:"userId,string"`
	Name       string            `json:"name"`
	Image      string            `json:"image"`
	ImageBig   string            `json:"imageBig"`
	Comment    string            `json:"comment"`
	Background UserRawBackground `json:"background"`
}

type UserRawResponse struct {
	Error   bool    `json:"error"`
	Message string  `json:"message"`
	Body    UserRaw `json:"body"`
}

type UserIllustsRaw struct {
	Illusts interface{} `json:"illusts"`
	Manga   interface{} `json:"manga"`
}

type UserIllustsRawResponse struct {
	Error   bool           `json:"error"`
	Message string         `json:"message"`
	Body    UserIllustsRaw `json:"body"`
}

type RankIllustsRawResponse struct {
	Contents []RankIllustRaw `json:"contents"`
	Mode     string          `json:"mode"`
	Page     int             `json:"page"`
	Date     string          `json:"date"`
	Next     interface{}     `json:"next"`
	Content  string          `json:"content"`
}

type RankIllustRaw struct {
	IllustID              uint64 `json:"illust_id"`
	UserID                uint   `json:"user_id"`
	Rank                  uint   `json:"rank"`
	YesRank               uint   `json:"yes_rank"`
	RatingCount           uint   `json:"rating_count"`
	IllustUploadTimestamp uint   `json:"illust_upload_timestamp"`
}

type UgoiraRawResponse struct {
	Error   bool      `json:"error"`
	Message string    `json:"message"`
	Body    UgoiraRaw `json:"body"`
}

type UgoiraFrameRaw struct {
	File  string `json:"file"`
	Delay int    `json:"delay"`
}

type UgoiraRaw struct {
	Src         string           `json:"src"`
	Originalsrc string           `json:"originalSrc"`
	MimeType    string           `json:"mime_type"`
	Frames      []UgoiraFrameRaw `json:"frames"`
}

type ErrorRawResponse struct {
	Message string `json:"message"`
}
