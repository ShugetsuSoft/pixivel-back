package convert

import (
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"strings"
	"time"

	"github.com/ShugetsuSoft/pixivel-back/common/models"
)

func ParseImgTime(url string) time.Time {
	date_str := url[strings.Index(url, "/img/")+5 : strings.LastIndex(url, "/")]
	t, _ := time.Parse("2006/01/02/15/04/05", date_str)
	return t
}

func IllustRaw2Illust(raw *models.IllustRaw) *models.Illust {
	illustTags := make([]models.IllustTag, len(raw.Tags.Tags))
	for i := 0; i < len(raw.Tags.Tags); i++ {
		illustTags[i] = models.IllustTag{
			Name:        raw.Tags.Tags[i].Tag,
			Translation: raw.Tags.Tags[i].Translation.En,
		}
	}
	banned := false
	if raw.XRestrict == 1 {
		banned = true
	}
	sta := models.IllustStatistic{
		Bookmarks: raw.BookmarkCount,
		Likes:     raw.LikeCount,
		Comments:  raw.CommentCount,
		Views:     raw.ViewCount,
	}
	return &models.Illust{
		ID:          raw.ID,
		Title:       raw.Title,
		AltTitle:    raw.Alt,
		Description: raw.Description,
		Type:        raw.IllustType,
		CreateDate:  raw.CreateDate,
		UploadDate:  raw.UploadDate,
		Sanity:      raw.Sl,
		Width:       raw.Width,
		Height:      raw.Height,
		PageCount:   raw.PageCount,
		User:        raw.UserID,
		Tags:        illustTags,
		Popularity:  CalcIllustPop(sta),
		Statistic:   sta,
		Image:       ParseImgTime(raw.Urls.Original),
		Banned:      banned,
	}
}

func CalcIllustPop(sta models.IllustStatistic) uint {
	return (sta.Bookmarks*70 + sta.Likes*30) / 100
}

func Illust2IllustSearch(illustdb *models.Illust) *models.IllustSearch {
	return &models.IllustSearch{
		Title:       illustdb.Title,
		AltTitle:    illustdb.AltTitle,
		Description: illustdb.Description,
		Type:        illustdb.Type,
		CreateDate:  illustdb.CreateDate,
		Sanity:      illustdb.Sanity,
		Popularity:  illustdb.Popularity,
		User:        illustdb.User,
		Tags:        illustdb.Tags,
	}
}

func User2UserSearch(userdb *models.User) *models.UserSearch {
	return &models.UserSearch{
		Name: userdb.Name,
		Bio:  userdb.Bio,
	}
}

func UserRaw2User(raw *models.UserRaw) *models.User {
	return &models.User{
		ID:   raw.UserID,
		Name: raw.Name,
		Bio:  raw.Comment,
		Image: models.UserImage{
			Url:        raw.Image,
			BigUrl:     raw.ImageBig,
			Background: raw.Background.URL,
		},
	}
}

func UserIllusts2UserIllustsResponse(uid uint64, raw *models.UserIllustsRaw) *models.UserIllustsResponse {
	lenth := 0
	switch illusts := raw.Illusts.(type) {
	case map[string]interface{}:
		lenth += len(illusts)
	}
	switch manga := raw.Manga.(type) {
	case map[string]interface{}:
		lenth += len(manga)
	}
	lis := make(models.IDList, lenth)
	i := 0
	switch illusts := raw.Illusts.(type) {
	case map[string]interface{}:
		for key := range illusts {
			lis[i] = utils.Atoi(key)
			i++
		}
	}
	switch manga := raw.Manga.(type) {
	case map[string]interface{}:
		for key := range manga {
			lis[i] = utils.Atoi(key)
			i++
		}
	}
	return &models.UserIllustsResponse{
		UserID:  uid,
		Illusts: lis,
	}
}

func RankIllusts2RankIllustsResponse(raw *models.RankIllustsRawResponse) *models.RankIllustsResponseMessage {
	illusts := make([]models.IDWithPos, len(raw.Contents))
	for key := range raw.Contents {
		illusts[key].ID = raw.Contents[key].IllustID
		illusts[key].Pos = raw.Contents[key].Rank
	}
	hasnext := true
	switch next := raw.Next.(type) {
	case bool:
		if next == false {
			hasnext = false
		}
	}
	return &models.RankIllustsResponseMessage{
		Mode:    raw.Mode,
		Date:    raw.Date,
		Content: raw.Content,
		Page:    raw.Page,
		Next:    hasnext,
		Illusts: illusts,
	}
}

func UgoiraRaw2Ugoira(raw *models.UgoiraRaw, id uint64) *models.Ugoira {
	frames := make([]models.UgoiraFrame, len(raw.Frames))
	for key := range raw.Frames {
		frames[key].Delay = raw.Frames[key].Delay
		frames[key].File = raw.Frames[key].File
	}
	return &models.Ugoira{
		Image:    ParseImgTime(raw.Src),
		MimeType: raw.MimeType,
		Frames:   frames,
		ID:       id,
	}
}
