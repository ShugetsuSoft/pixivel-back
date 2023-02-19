package convert

import "github.com/ShugetsuSoft/pixivel-back/common/models"

func Illust2IllustResponse(illust *models.Illust, user *models.User) *models.IllustResponse {
	var userres *models.UserResponse
	if user != nil {
		userres = User2UserResponse(user)
	} else {
		userres = nil
	}
	illustTags := Tags2TagResponses(illust.Tags)

	return &models.IllustResponse{
		ID:          illust.ID,
		Title:       illust.Title,
		AltTitle:    illust.AltTitle,
		Description: illust.Description,
		Type:        illust.Type,
		CreateDate:  illust.CreateDate.Format("2006-01-02T15:04:05"),
		UploadDate:  illust.UploadDate.Format("2006-01-02T15:04:05"),
		Sanity:      illust.Sanity,
		Width:       illust.Width,
		Height:      illust.Height,
		PageCount:   illust.PageCount,
		Tags:        illustTags,
		Statistic: models.IllustStatisticResponse{
			Bookmarks: illust.Statistic.Bookmarks,
			Likes:     illust.Statistic.Likes,
			Comments:  illust.Statistic.Comments,
			Views:     illust.Statistic.Views,
		},
		User:   userres,
		Image:  illust.Image.Format("2006-01-02T15:04:05"),
		AIType: illust.AIType,
	}
}

func User2UserResponse(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Bio:  user.Bio,
		Image: models.UserImageResponse{
			Url:        user.Image.Url,
			BigUrl:     user.Image.BigUrl,
			Background: user.Image.Background,
		},
	}
}

func Illusts2IllustsResponse(illusts []models.Illust, hasnext bool) *models.IllustsResponse {
	illustsresponse := make([]models.IllustResponse, len(illusts))
	for i := range illusts {
		illustsresponse[i] = *Illust2IllustResponse(&illusts[i], nil)
	}
	return &models.IllustsResponse{
		Illusts: illustsresponse,
		HasNext: hasnext,
	}
}

func Illusts2IllustsSearchResponse(illusts []models.Illust, hasnext bool, scores []float64, highlights []*string) *models.IllustsSearchResponse {
	illustsresponse := make([]models.IllustResponse, len(illusts))
	for i := range illusts {
		illustsresponse[i] = *Illust2IllustResponse(&illusts[i], nil)
	}
	return &models.IllustsSearchResponse{
		Illusts:   illustsresponse,
		Scores:    scores,
		HighLight: highlights,
		HasNext:   hasnext,
	}
}

func Users2UsersSearchResponse(users []models.User, hasnext bool, scores []float64, highlights []*string) *models.UsersSearchResponse {
	usersesponse := make([]models.UserResponse, len(users))
	for i := range users {
		usersesponse[i] = *User2UserResponse(&users[i])
	}
	return &models.UsersSearchResponse{
		Users:     usersesponse,
		Scores:    scores,
		HighLight: highlights,
		HasNext:   hasnext,
	}
}

func Tags2TagResponses(tags []models.IllustTag) []models.IllustTagResponse {
	var illustTags []models.IllustTagResponse

	illustTags = make([]models.IllustTagResponse, len(tags))
	for i := 0; i < len(tags); i++ {
		illustTags[i] = models.IllustTagResponse{
			Name:        tags[i].Name,
			Translation: tags[i].Translation,
		}
	}
	return illustTags
}

func RankAggregateResult2IllustsResponses(rank []models.RankAggregateResult, hasnext bool) *models.IllustsResponse {
	illusts := make([]models.Illust, len(rank))
	for i := range rank {
		illusts[i] = rank[i].Content
	}
	return Illusts2IllustsResponse(illusts, hasnext)
}

func Ugoira2UgoiraResponse(ugoira *models.Ugoira) *models.UgoiraResponse {
	frames := make([]models.UgoiraFrameResponse, len(ugoira.Frames))
	for key := range ugoira.Frames {
		frames[key].Delay = ugoira.Frames[key].Delay
		frames[key].File = ugoira.Frames[key].File
	}
	return &models.UgoiraResponse{
		Image:    ugoira.Image.Format("2006-01-02T15:04:05"),
		MimeType: ugoira.MimeType,
		Frames:   frames,
		ID:       ugoira.ID,
	}
}

func Users2UsersResponse(users []models.User, hasnext bool) *models.UsersResponse {
	usersresponse := make([]models.UserResponse, len(users))
	for i := range users {
		usersresponse[i] = *User2UserResponse(&users[i])
	}
	return &models.UsersResponse{
		Users:   usersresponse,
		HasNext: hasnext,
	}
}
