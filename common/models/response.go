package models

type UserImageResponse struct {
	Url        string `json:"url"`
	BigUrl     string `json:"bigUrl"`
	Background string `json:"background,omitempty"`
}

type UserResponse struct {
	ID    uint64            `json:"id"`
	Name  string            `json:"name"`
	Bio   string            `json:"bio"`
	Image UserImageResponse `json:"image"`
}

type IllustTagResponse struct {
	Name        string `json:"name"`
	Translation string `json:"translation,omitempty"`
}

type IllustStatisticResponse struct {
	Bookmarks uint `json:"bookmarks"`
	Likes     uint `json:"likes"`
	Comments  uint `json:"comments"`
	Views     uint `json:"views"`
}

type IllustResponse struct {
	ID          uint64                  `json:"id"`
	Title       string                  `json:"title"`
	AltTitle    string                  `json:"altTitle"`
	Description string                  `json:"description"`
	Type        uint                    `json:"type"`
	CreateDate  string                  `json:"createDate"`
	UploadDate  string                  `json:"uploadDate"`
	Sanity      uint                    `json:"sanity"`
	Width       uint                    `json:"width"`
	Height      uint                    `json:"height"`
	PageCount   uint                    `json:"pageCount"`
	Tags        []IllustTagResponse     `json:"tags"`
	Statistic   IllustStatisticResponse `json:"statistic"`
	User        *UserResponse           `json:"user,omitempty"`
	Image       string                  `json:"image"`
}

type IllustsResponse struct {
	Illusts []IllustResponse `json:"illusts"`
	HasNext bool             `json:"has_next"`
}

type UsersResponse struct {
	Users   []UserResponse `json:"users"`
	HasNext bool           `json:"has_next"`
}

type UsersSearchResponse struct {
	Users     []UserResponse `json:"users"`
	Scores    []float64      `json:"scores"`
	HighLight []*string      `json:"highlight,omitempty"`
	HasNext   bool           `json:"has_next"`
}

type IllustsSearchResponse struct {
	Illusts   []IllustResponse `json:"illusts"`
	Scores    []float64        `json:"scores,omitempty"`
	HighLight []*string        `json:"highlight,omitempty"`
	HasNext   bool             `json:"has_next"`
}

type SearchSuggestResponse struct {
	SuggestWords []string `json:"suggest_words"`
}

type SearchSuggestTagsResponse struct {
	SuggestTags []IllustTagResponse `json:"suggest_tags"`
}

type UgoiraResponse struct {
	ID       uint64                `json:"id"`
	Image    string                `json:"image"`
	MimeType string                `json:"mime_type"`
	Frames   []UgoiraFrameResponse `json:"frames"`
}

type UgoiraFrameResponse struct {
	File  string `json:"file"`
	Delay int    `json:"delay"`
}
