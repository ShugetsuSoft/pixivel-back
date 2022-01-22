package reader

import (
	"github.com/ShugetsuSoft/pixivel-back/common/convert"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/config"
)

func (r *Reader) SearchIllustsResponse(keyword string, page int, limit int, sortpopularity bool, sortdate bool) (*models.IllustsSearchResponse, error) {
	illusts, hits, scores, highlights, err := r.dbops.SearchIllust(keyword, page, limit, sortpopularity, sortdate, false)
	if err != nil {
		return nil, err
	}

	return convert.Illusts2IllustsSearchResponse(illusts, hits > int64(limit*(page+1)), scores, highlights), err
}

func (r *Reader) SearchIllustsSuggestResponse(keyword string) (*models.SearchSuggestResponse, error) {
	suggests, err := r.dbops.SearchIllustSuggest(keyword)
	if err != nil {
		return nil, err
	}

	return &models.SearchSuggestResponse{
		SuggestWords: suggests,
	}, nil
}

func (r *Reader) SearchUsersResponse(keyword string, page int, limit int) (*models.UsersSearchResponse, error) {
	users, hits, scores, highlights, err := r.dbops.SearchUser(keyword, page, limit, false)
	if err != nil {
		return nil, err
	}

	return convert.Users2UsersSearchResponse(users, hits > int64(limit*(page+1)), scores, highlights), err
}

func (r *Reader) SearchUsersSuggestResponse(keyword string) (*models.SearchSuggestResponse, error) {
	suggests, err := r.dbops.SearchUserSuggest(keyword)
	if err != nil {
		return nil, err
	}

	return &models.SearchSuggestResponse{
		SuggestWords: suggests,
	}, nil
}

func (r *Reader) SearchTagsSuggestResponse(keyword string) (*models.SearchSuggestTagsResponse, error) {
	suggests, err := r.dbops.SearchTagSuggest(keyword)
	if err != nil {
		return nil, err
	}

	return &models.SearchSuggestTagsResponse{
		SuggestTags: convert.Tags2TagResponses(suggests),
	}, nil
}

func (r *Reader) SearchIllustsByTagsResponse(musttags []string, shouldtags []string, perfectmatch bool, page int, limit int, sortpopularity bool, sortdate bool) (*models.IllustsResponse, error) {
	if perfectmatch {
		illusts, err := r.dbops.QueryIllustsByTags(musttags, shouldtags, int64(page), int64(limit), sortpopularity, sortdate, false)
		if err != nil {
			return nil, err
		}

		return convert.Illusts2IllustsResponse(illusts, len(illusts) >= limit), nil
	}
	return nil, nil
}

func (r *Reader) RecommendIllustsByIllustId(illustId uint64, k int) ([]models.Illust, error) {
	illusts, err := r.dbops.RecommendIllustsByIllustId(illustId, k, config.RecommendDrift, false)
	if err != nil {
		return nil, err
	}
	return illusts, nil
}
