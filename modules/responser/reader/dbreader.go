package reader

import (
	"context"
	"time"

	"github.com/ShugetsuSoft/pixivel-back/common/convert"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
)

func (r *Reader) IllustResponse(ctx context.Context, illustId uint64, forcefetch bool) (*models.IllustResponse, error) {
	retry := 1
START:
	illust, err := r.dbops.QueryIllust(ctx, illustId, false)
	if err != nil {
		return nil, err
	}

	if illust == nil || forcefetch && time.Now().After(illust.UpdateTime.Add(time.Hour*24*2)) {
		if r.mode == models.ArchiveMode {
			return nil, models.ErrorArchiveMode
		}
		if retry == 0 {
			return nil, models.ErrorRetrivingFinishedTask
		}
		err = r.gen.IllustDetailTask(ctx, illustId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START

	}

	userId := uint64(illust.User)
	user, err := r.dbops.QueryUser(ctx, userId, false)
	if err != nil {
		return nil, err
	}

	if user == nil {
		if r.mode == models.ArchiveMode {
			return nil, models.ErrorArchiveMode
		}
		if retry == 0 {
			return nil, models.ErrorRetrivingFinishedTask
		}
		err = r.gen.UserDetailTask(ctx, userId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START
	}

	return convert.Illust2IllustResponse(illust, user), nil
}

func (r *Reader) UgoiraResponse(ctx context.Context, ugoiraId uint64, forcefetch bool) (*models.UgoiraResponse, error) {
	retry := 1
START:
	ugoira, err := r.dbops.QueryUgoira(ctx, ugoiraId)
	if err != nil {
		return nil, err
	}

	if ugoira == nil || forcefetch && time.Now().After(ugoira.UpdateTime.Add(time.Hour*24*2)) {
		if r.mode == models.ArchiveMode {
			return nil, models.ErrorArchiveMode
		}
		if retry == 0 {
			return nil, models.ErrorRetrivingFinishedTask
		}
		err = r.gen.UgoiraDetailTask(ctx, ugoiraId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START
	}

	return convert.Ugoira2UgoiraResponse(ugoira), nil
}

func (r *Reader) IllustsResponse(ctx context.Context, illustIds []uint64) (*models.IllustsResponse, error) {
	var err error
	var illust *models.Illust
	illusts := make([]models.Illust, 0, len(illustIds))
	for _, i := range illustIds {
		illust, err = r.dbops.QueryIllust(ctx, illustIds[i], false)
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

func (r *Reader) UserDetailResponse(ctx context.Context, userId uint64) (*models.UserResponse, error) {
	retry := 1
START:
	user, err := r.dbops.QueryUser(ctx, userId, false)
	if err != nil {
		return nil, err
	}

	if user == nil {
		if r.mode == models.ArchiveMode {
			return nil, models.ErrorArchiveMode
		}
		if retry == 0 {
			return nil, models.ErrorRetrivingFinishedTask
		}
		err = r.gen.UserDetailTask(ctx, userId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START
	}

	return convert.User2UserResponse(user), nil
}

func (r *Reader) UserIllustsResponse(ctx context.Context, userId uint64, page int64, limit int64) (*models.IllustsResponse, error) {
	retry := 1
START:
	user, err := r.dbops.QueryUser(ctx, userId, false)
	if err != nil {
		return nil, err
	}

	if user == nil {
		if r.mode == models.ArchiveMode {
			return nil, models.ErrorArchiveMode
		}
		if retry == 0 {
			return nil, models.ErrorRetrivingFinishedTask
		}
		err = r.gen.UserDetailTask(ctx, userId)
		if err != nil {
			return nil, err
		}
		err = r.gen.UserIllustsTask(ctx, userId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START
	}

	if (time.Now().After(user.IllustsUpdateTime.Add(time.Hour*24*2)) || user.IllustsCount == 0) && r.mode != models.ArchiveMode {
		if retry == 0 {
			goto END
		}
		err = r.gen.UserIllustsTask(ctx, userId)
		if err != nil {
			return nil, err
		}
		retry--
		goto START
	}

END:
	illusts, err := r.dbops.QueryIllustByUserWithPage(ctx, userId, page, limit, false)
	if err != nil {
		return nil, err
	}

	return convert.Illusts2IllustsResponse(illusts, user.IllustsCount > uint(limit*(page+1))), err
}

func (r *Reader) RankIllustsResponse(ctx context.Context, mode string, date string, page int, content string, limit int) (*models.IllustsResponse, error) {
	results, err := r.dbops.QueryRankIllusts(ctx, mode, date, content, page, limit)
	if err != nil {
		return nil, err
	}
	return convert.RankAggregateResult2IllustsResponses(results, page < 9 && len(results) != 0), nil
}

func (r *Reader) SampleIllustsResponse(ctx context.Context, quality int, limit int) (*models.IllustsResponse, error) {
	results, err := r.dbops.GetSampleIllusts(ctx, quality, limit, false)
	if err != nil {
		return nil, err
	}
	return convert.Illusts2IllustsResponse(results, false), nil
}

func (r *Reader) SampleUsersResponse(ctx context.Context, limit int) (*models.UsersResponse, error) {
	results, err := r.dbops.GetSampleUsers(ctx, limit, false)
	if err != nil {
		return nil, err
	}
	return convert.Users2UsersResponse(results, false), nil
}
