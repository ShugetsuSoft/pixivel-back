package operations

import (
	"github.com/ShugetsuSoft/pixivel-back/common/convert"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/config"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (ops *DatabaseOperations) InsertUser(user *models.User) error {
	var err error
	is, err := ops.Flt.Exists(config.UserTableName, utils.Itoa(user.ID))
	if err != nil {
		return err
	}
	user.UpdateTime = time.Now()

	if is {
		goto REPLACE
	} else {
		user.IllustsUpdateTime = time.Unix(0, 0)
		user.IllustsCount = 0
		_, err = ops.Cols.User.InsertOne(ops.Ctx, user)

		if mongo.IsDuplicateKeyError(err) {
			goto REPLACE
		}

		if err != nil {
			return err
		}

		_, err = ops.Flt.Add(config.UserTableName, utils.Itoa(user.ID))
		if err != nil {
			return err
		}
	}

	return nil

REPLACE:
	result, err := ops.Cols.User.ReplaceOne(ops.Ctx, bson.M{"_id": user.ID}, user)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		_, err = ops.Cols.User.InsertOne(ops.Ctx, user)
		if err != nil {
			return err
		}
		_, err = ops.Flt.Add(config.UserTableName, utils.Itoa(user.ID))
	}

	return err
}

func (ops *DatabaseOperations) UpdateUserIllustsTime(userId uint64) error {
	var err error
	is, err := ops.Flt.Exists(config.UserTableName, utils.Itoa(userId))
	if err != nil {
		return err
	}
	if is {
		_, err = ops.Cols.User.UpdateOne(ops.Ctx, bson.M{"_id": userId}, bson.M{"$set": bson.M{
			"illusts_update_time": time.Now(),
			"update_time":         time.Now(),
		}})
		return err
	}
	return nil
}

func (ops *DatabaseOperations) InsertUserSearch(user *models.User) error {
	usersearch := convert.User2UserSearch(user)
	err := ops.Sc.es.InsertDocument(config.UserSearchIndexName, utils.Itoa(user.ID), usersearch)
	if err != nil {
		return err
	}
	return nil
}

func (ops *DatabaseOperations) QueryUser(userId uint64, resultbanned bool) (*models.User, error) {
	is, err := ops.Flt.Exists(config.UserTableName, utils.Itoa(userId))

	if err != nil {
		return nil, err
	}

	if is {
		result := models.User{
			Image: models.UserImage{},
		}
		err := ops.Cols.User.FindOne(ops.Ctx, bson.M{"_id": userId}).Decode(&result)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			} else {
				return nil, err
			}
		}

		if result.Banned && !resultbanned {
			return nil, models.ErrorItemBanned
		}

		return &result, err
	}
	return nil, nil
}

func (ops *DatabaseOperations) QueryUsers(userIds []uint64, resultbanned bool) ([]models.User, error) {
	query := bson.M{"_id": bson.M{"$in": userIds}}
	cursor, err := ops.Cols.User.Find(ops.Ctx, query)
	defer cursor.Close(ops.Ctx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	users := make([]models.User, 0, len(userIds))

	for cursor.Next(ops.Ctx) {
		result := models.User{
			Image: models.UserImage{},
		}
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}

		if result.Banned && !resultbanned {
			continue
		}

		users = append(users, result)
	}

	return users, err
}

func (ops *DatabaseOperations) SearchUserSuggest(keyword string) ([]string, error) {
	source := elastic.NewSearchSource().
		Suggester(ops.Sc.es.Suggest("name-completion-suggest").Field("name.suggest").Text(keyword).Fuzziness(2).Analyzer("kuromoji")).
		FetchSource(false).TrackScores(true)
	query := ops.Sc.es.Search(config.UserSearchIndexName).Source(nil).
		SearchSource(source)

	results, err := ops.Sc.es.DoSearch(query)
	if err != nil {
		return nil, err
	}
	suggests := results.Suggest["name-completion-suggest"][0].Options
	res := make([]string, len(suggests))
	for i, suggest := range suggests {
		res[i] = suggest.Text
	}
	return res, nil
}

func (ops *DatabaseOperations) SearchUser(keyword string, page int, limit int, resultbanned bool) ([]models.User, int64, []float64, []*string, error) {
	query := ops.Sc.es.Search(config.UserSearchIndexName).
		Query(ops.Sc.es.BoolQuery().
			Should(ops.Sc.es.Query("name", keyword).Boost(3)).
			Should(ops.Sc.es.Query("bio", keyword).Boost(1)),
		).
		Size(limit).From(page * limit).
		Highlight(elastic.NewHighlight().Field("name")).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("_id"))

	results, err := ops.Sc.es.DoSearch(query)
	if err != nil {
		return nil, 0, nil, nil, err
	}

	if results.Hits.TotalHits.Value > 0 {
		scores := make([]float64, 0, len(results.Hits.Hits))
		highlights := make([]*string, 0, len(results.Hits.Hits))
		userids := make([]uint64, len(results.Hits.Hits))
		for i, hit := range results.Hits.Hits {
			userids[i] = utils.Atoi(hit.Id)
			if hit.Score != nil {
				scores = append(scores, *hit.Score)
			} else {
				scores = append(scores, -1)
			}
			if highl, ok := hit.Highlight["name"]; ok {
				highlights = append(highlights, &highl[0])
			} else {
				highlights = append(highlights, nil)
			}
		}

		users, err := ops.QueryUsers(userids, resultbanned)
		if err != nil {
			return nil, 0, nil, nil, err
		}

		usersmap := make(map[uint64]models.User)
		for _, user := range users {
			usersmap[user.ID] = user
		}

		result := make([]models.User, 0, len(results.Hits.Hits))
		for _, userid := range userids {
			if _, exist := usersmap[userid]; exist {
				result = append(result, usersmap[userid])
			}
		}

		return result, results.Hits.TotalHits.Value, scores, highlights, nil
	} else {
		return nil, 0, nil, nil, models.ErrorNoResult
	}
	return nil, 0, nil, nil, err
}

func (ops *DatabaseOperations) DeleteUser(userId uint64) error {
	is, err := ops.Flt.Exists(config.UserTableName, utils.Itoa(userId))

	if err != nil {
		return err
	}

	if is {
		_, err := ops.Cols.User.DeleteOne(ops.Ctx, bson.M{"_id": userId})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil
			}
			return err
		}

		return err
	}

	return nil
}

func (ops *DatabaseOperations) ClearUserIllusts(userId uint64) error {
	illusts, err := ops.QueryIllustByUser(userId, true)
	if err != nil {
		return err
	}

	for i := 0; i < len(illusts); i++ {
		err = ops.DeleteIllust((illusts)[i].ID)
		if err != nil {
			return err
		}
	}

	is, err := ops.Flt.Exists(config.UserTableName, utils.Itoa(userId))

	if err != nil {
		return err
	}

	if is {
		_, err = ops.Cols.User.UpdateOne(ops.Ctx, bson.M{"_id": userId}, bson.M{"$set": bson.M{
			"illusts_update_time": time.Unix(0, 0),
			"illusts_count":       0,
			"update_time":         time.Now(),
		}})
		return err
	}
	return nil
}

func (ops *DatabaseOperations) SetIllustsCount(userId uint64, count uint) error {
	var err error
	is, err := ops.Flt.Exists(config.UserTableName, utils.Itoa(userId))
	if err != nil {
		return err
	}
	if is {
		_, err = ops.Cols.User.UpdateOne(ops.Ctx, bson.M{"_id": userId}, bson.M{"$set": bson.M{
			"illusts_update_time": time.Now(),
			"illusts_count":       count,
			"update_time":         time.Now(),
		}})
		return err
	}
	return nil
}

func (ops *DatabaseOperations) IsUserExist(userId uint64) (bool, error) {
	is, err := ops.Flt.Exists(config.UserTableName, utils.Itoa(userId))
	if err != nil {
		return false, err
	}
	return is, nil
}
