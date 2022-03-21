package reader

import (
	"time"

	"github.com/ShugetsuSoft/pixivel-back/common/convert"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
)

func (r *Reader) IllustResponse(illustId uint64, forcefetch bool) (*models.IllustResponse, error) {
	retry := 2
START:
	illust, err := r.dbops.QueryIllust(illustId, false)
	if err != nil {
		return nil, err
	}

	if illust == nil || forcefetch && time.Now().After(illust.UpdateTime.Add(time.Hour*24*2)) {
		if retry == 0 {
			return nil, models.ErrorRetrivingFinishedTask
		}
		err = r.gen.IllustDetailTask(illustId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START
	}

	userId := uint64(illust.User)
	user, err := r.dbops.QueryUser(userId, false)
	if err != nil {
		return nil, err
	}

	if user == nil {
		if retry == 0 {
			return nil, models.ErrorRetrivingFinishedTask
		}
		err = r.gen.UserDetailTask(userId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START
	}

	return convert.Illust2IllustResponse(illust, user), nil
}

func (r *Reader) UgoiraResponse(ugoiraId uint64, forcefetch bool) (*models.UgoiraResponse, error) {
	retry := 2
START:
	ugoira, err := r.dbops.QueryUgoira(ugoiraId)
	if err != nil {
		return nil, err
	}

	if ugoira == nil || forcefetch && time.Now().After(ugoira.UpdateTime.Add(time.Hour*24*2)) {
		if retry == 0 {
			return nil, models.ErrorRetrivingFinishedTask
		}
		err = r.gen.UgoiraDetailTask(ugoiraId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START
	}

	return convert.Ugoira2UgoiraResponse(ugoira), nil
}

func (r *Reader) IllustsResponse(illustIds []uint64) (*models.IllustsResponse, error) {
	var err error
	var illust *models.Illust
	illusts := make([]models.Illust, 0, len(illustIds))
	for _, i := range illustIds {
		illust, err = r.dbops.QueryIllust(illustIds[i], false)
		if err != nil {
			return nil, err
		}
		if illust == nil {
			continue
		}
		illusts = append(illusts, *illust)
	}
	return convert.Illusts2IllustsResponse(illusts, false), nil
}

func (r *Reader) UserDetailResponse(userId uint64) (*models.UserResponse, error) {
	retry := 1
START:
	user, err := r.dbops.QueryUser(userId, false)
	if err != nil {
		return nil, err
	}

	if user == nil {
		if retry == 0 {
			return nil, models.ErrorRetrivingFinishedTask
		}
		err = r.gen.UserDetailTask(userId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START
	}

	return convert.User2UserResponse(user), nil
}

func (r *Reader) UserIllustsResponse(userId uint64, page int64, limit int64) (*models.IllustsResponse, error) {
	retry := 2
START:
	user, err := r.dbops.QueryUser(userId, false)
	if err != nil {
		return nil, err
	}

	if user == nil {
		if retry == 0 {
			return nil, models.ErrorRetrivingFinishedTask
		}
		err = r.gen.UserDetailTask(userId)
		if err != nil {
			return nil, err
		}
		err = r.gen.UserIllustsTask(userId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START
	}

	if time.Now().After(user.IllustsUpdateTime.Add(time.Hour*24*2)) || user.IllustsCount == 0 {
		if retry == 0 {
			return nil, models.ErrorRetrivingFinishedTask
		}
		err = r.gen.UserIllustsTask(userId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START
	}

	illusts, err := r.dbops.QueryIllustByUserWithPage(userId, page, limit, false)
	if err != nil {
		return nil, err
	}

	return convert.Illusts2IllustsResponse(illusts, user.IllustsCount > uint(limit*(page+1))), err
}

func (r *Reader) RankIllustsResponse(mode string, date string, page int, content string, limit int) (*models.IllustsResponse, error) {
	results, err := r.dbops.QueryRankIllusts(mode, date, content, page, limit)
	if err != nil {
		return nil, err
	}
	return convert.RankAggregateResult2IllustsResponses(results, page < 9 && len(results) != 0), nil
}

func (r *Reader) SampleIllustsResponse(quality int, limit int) (*models.IllustsResponse, error) {
	results, err := r.dbops.GetSampleIllusts(quality, limit, false)
	if err != nil {
		return nil, err
	}
	return convert.Illusts2IllustsResponse(results, false), nil
}

func (r *Reader) SampleUsersResponse(limit int) (*models.UsersResponse, error) {
	results, err := r.dbops.GetSampleUsers(limit, false)
	if err != nil {
		return nil, err
	}
	return convert.Users2UsersResponse(results, false), nil
}
