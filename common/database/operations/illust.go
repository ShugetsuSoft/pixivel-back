package operations

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ShugetsuSoft/pixivel-back/common/convert"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/config"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (ops *DatabaseOperations) InsertIllust(ctx context.Context, illust *models.Illust) error {
	var err error
	is, err := ops.Flt.Exists(config.IllustTableName, utils.Itoa(illust.ID))
	if err != nil {
		return err
	}
	illust.UpdateTime = time.Now()

	if is {
		goto REPLACE
	} else {
		_, err = ops.Cols.Illust.InsertOne(ctx, illust)

		if mongo.IsDuplicateKeyError(err) {
			_, err = ops.Flt.Add(config.IllustTableName, utils.Itoa(illust.ID))
			if err != nil {
				return err
			}
			goto REPLACE
		}

		if err != nil {
			return err
		}

		_, err = ops.Flt.Add(config.IllustTableName, utils.Itoa(illust.ID))
		if err != nil {
			return err
		}
	}

	return nil

REPLACE:
	result, err := ops.Cols.Illust.ReplaceOne(ctx, bson.M{"_id": illust.ID}, illust)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		_, err = ops.Cols.Illust.InsertOne(ctx, illust)
		if err != nil {
			return err
		}
		_, err = ops.Flt.Add(config.IllustTableName, utils.Itoa(illust.ID))
	}

	return err
}

func (ops *DatabaseOperations) AddIllusts(ctx context.Context, illusts []models.Illust) error {
	operations := make([]mongo.WriteModel, len(illusts))
	for i, illust := range illusts {
		ins := mongo.NewInsertOneModel()
		ins.SetDocument(illust)
		operations[i] = ins
	}
	bulkOption := &options.BulkWriteOptions{}
	bulkOption.SetOrdered(false)
	bulkOption.SetBypassDocumentValidation(false)
	_, err := ops.Cols.Illust.BulkWrite(ctx, operations, bulkOption)
	return err
}

func (ops *DatabaseOperations) InsertIllusts(ctx context.Context, illusts []models.Illust) error {
	nowIllusts := make([]interface{}, 0, len(illusts))
	updateoperations := make([]mongo.WriteModel, 0, len(illusts))

	for _, illust := range illusts {
		is, err := ops.Flt.Exists(config.IllustTableName, utils.Itoa(illust.ID))
		if err != nil {
			return err
		}
		if is {
			rep := mongo.NewReplaceOneModel()
			rep.SetFilter(bson.M{"_id": illust.ID})
			rep.SetReplacement(illust)
			rep.SetUpsert(true)
			updateoperations = append(updateoperations, rep)
		} else {
			_, err := ops.Flt.Add(config.IllustTableName, utils.Itoa(illust.ID))
			if err != nil {
				return err
			}
			nowIllusts = append(nowIllusts, illust)
		}
	}

	if len(nowIllusts) > 0 {
		_, err := ops.Cols.Illust.InsertMany(ctx, nowIllusts)
		if err != nil {
			return err
		}
	}

	if len(updateoperations) > 0 {
		bulkOption := &options.BulkWriteOptions{}
		bulkOption.SetOrdered(false)
		bulkOption.SetBypassDocumentValidation(false)
		_, err := ops.Cols.Illust.BulkWrite(ctx, updateoperations, bulkOption)
		return err
	}
	return nil
}

func (ops *DatabaseOperations) InsertIllustSearch(ctx context.Context, illust *models.Illust) error {
	illustsearch := convert.Illust2IllustSearch(illust)
	err := ops.Sc.es.InsertDocument(ctx, config.IllustSearchIndexName, utils.Itoa(illust.ID), illustsearch)
	if err != nil {
		return err
	}
	return nil
}

func (ops *DatabaseOperations) InsertIllustTagNearDB(ctx context.Context, illust *models.Illust) error {
	tagset := make([]string, len(illust.Tags))
	for i, tag := range illust.Tags {
		tagset[i] = tag.Name
	}

	return ops.Sc.ndb.Add(ctx, illust.ID, tagset, illust.Popularity)
}

func (ops *DatabaseOperations) RecommendIllustsByIllustId(ctx context.Context, illustId uint64, k int, drif float64, resultbanned bool) ([]models.Illust, error) {
	items, err := ops.Sc.ndb.QueryById(ctx, illustId, k, drif)
	if err != nil {
		return nil, err
	}
	if len(items) < 1 {
		illust, err := ops.QueryIllust(ctx, illustId, resultbanned)
		if err != nil {
			return nil, err
		}
		tags := make([]string, len(illust.Tags))
		for i, tag := range illust.Tags {
			tags[i] = tag.Name
		}
		err = ops.Sc.ndb.Add(ctx, illust.ID, tags, illust.Popularity)
		if err != nil {
			return nil, err
		}
		items, err = ops.Sc.ndb.QueryById(ctx, illustId, k, drif)
		if err != nil {
			return nil, err
		}
	}
	queryidlist := make([]uint64, 0, len(items))
	for _, item := range items {
		if item.Id == illustId {
			continue
		}
		queryidlist = append(queryidlist, item.Id)
	}

	illusts, err := ops.QueryIllusts(ctx, queryidlist, resultbanned)
	if err != nil {
		return nil, err
	}

	illustsmap := make(map[uint64]models.Illust)
	for _, illust := range illusts {
		illustsmap[illust.ID] = illust
	}

	res := make([]models.Illust, 0, len(items))
	for _, illustid := range queryidlist {
		if _, exist := illustsmap[illustid]; exist {
			res = append(res, illustsmap[illustid])
		}
	}

	return res, nil
}

func (ops *DatabaseOperations) QueryIllust(ctx context.Context, illustId uint64, resultbanned bool) (*models.Illust, error) {
	is, err := ops.Flt.Exists(config.IllustTableName, utils.Itoa(illustId))

	if err != nil {
		return nil, err
	}

	if is {
		result := models.Illust{
			Statistic: models.IllustStatistic{},
			Tags:      []models.IllustTag{},
		}
		query := bson.M{"_id": illustId}
		err := ops.Cols.Illust.FindOne(ctx, query).Decode(&result)
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

func (ops *DatabaseOperations) QueryIllusts(ctx context.Context, illustIds []uint64, resultbanned bool) ([]models.Illust, error) {
	query := bson.M{
		"_id":    bson.M{"$in": illustIds},
		"banned": false,
	}
	if resultbanned {
		query["banned"] = true
	}
	cursor, err := ops.Cols.Illust.Find(ctx, query)
	defer cursor.Close(ctx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	illusts := make([]models.Illust, 0, len(illustIds))

	for cursor.Next(ctx) {
		result := models.Illust{
			Statistic: models.IllustStatistic{},
			Tags:      []models.IllustTag{},
		}
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}

		illusts = append(illusts, result)
	}

	return illusts, err
}

func (ops *DatabaseOperations) SearchIllustSuggest(ctx context.Context, keyword string) ([]string, error) {
	source := elastic.NewSearchSource().
		Suggester(ops.Sc.es.Suggest("title-completion-suggest").Field("title.suggest").Text(keyword).Fuzziness(2).Analyzer("kuromoji").SkipDuplicates(true)).
		Suggester(ops.Sc.es.Suggest("alt_title-completion-suggest").Field("alt_title.suggest").Text(keyword).Fuzziness(3).Analyzer("kuromoji").SkipDuplicates(true)).
		Suggester(ops.Sc.es.Suggest("name-completion-suggest").Field("tags.name.suggest").Text(keyword).Fuzziness(2).Analyzer("kuromoji").SkipDuplicates(true)).
		Suggester(ops.Sc.es.Suggest("trans-completion-suggest").Field("tags.translation.suggest").Text(keyword).Fuzziness(2).Analyzer("smartcn").SkipDuplicates(true)).
		FetchSource(false)
	query := ops.Sc.es.Search(config.IllustSearchIndexName).Source(nil).
		SearchSource(source)

	results, err := ops.Sc.es.DoSearch(ctx, query)
	if err != nil {
		return nil, err
	}
	suggests := append(append(append(results.Suggest["title-completion-suggest"][0].Options, results.Suggest["alt_title-completion-suggest"][0].Options...), results.Suggest["name-completion-suggest"][0].Options...), results.Suggest["trans-completion-suggest"][0].Options...)
	res := make([]string, 0, len(suggests))
	uniqmap := make(map[string]bool, len(suggests))
	for _, suggest := range suggests {
		if uniqmap[suggest.Text] == false {
			uniqmap[suggest.Text] = true
			res = append(res, suggest.Text)
		}
	}
	return res, nil
}

func (ops *DatabaseOperations) SearchTagSuggest(ctx context.Context, keyword string) ([]models.IllustTag, error) {
	source := elastic.NewSearchSource().
		Suggester(ops.Sc.es.Suggest("name-completion-suggest").Field("tags.name.suggest").Text(keyword).Fuzziness(2).Analyzer("kuromoji").SkipDuplicates(true)).
		Suggester(ops.Sc.es.Suggest("trans-completion-suggest").Field("tags.translation.suggest").Text(keyword).Fuzziness(2).Analyzer("smartcn").SkipDuplicates(true)).
		FetchSource(true)
	query := ops.Sc.es.Search(config.IllustSearchIndexName).Source(nil).
		SearchSource(source)

	results, err := ops.Sc.es.DoSearch(ctx, query)
	if err != nil {
		return nil, err
	}
	suggests := append(results.Suggest["name-completion-suggest"][0].Options, results.Suggest["trans-completion-suggest"][0].Options...)
	//suggests := results.Suggest["name-completion-suggest"][0].Options
	res := make([]models.IllustTag, len(suggests))
	for i, suggest := range suggests {
		var tag models.IllustTag
		json.Unmarshal(suggest.Source, &tag)
		res[i] = tag
	}
	return res, nil
}

func (ops *DatabaseOperations) SearchIllust(ctx context.Context, keyword string, page int, limit int, sortpopularity bool, sortdate bool, resultbanned bool) ([]models.Illust, int64, []float64, []*string, error) {
	query := ops.Sc.es.Search(config.IllustSearchIndexName).
		Query(ops.Sc.es.BoolQuery().
			Should(ops.Sc.es.Query("title", keyword).Analyzer("kuromoji").Boost(4)).
			Should(elastic.NewMatchQuery("title.keyword", keyword).Boost(3)).
			Should(ops.Sc.es.Query("alt_title", keyword).Analyzer("kuromoji").Boost(2)).
			Should(elastic.NewNestedQuery("tags",
				ops.Sc.es.BoolQuery().Should(ops.Sc.es.Query("tags.name.fuzzy", keyword).Analyzer("kuromoji").Boost(2)).
					Should(ops.Sc.es.Query("tags.translation.fuzzy", keyword).Analyzer("kuromoji").Boost(1))),
			),
		).
		Size(limit).From(page * limit).
		Highlight(elastic.NewHighlight().Field("title")).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("_id")).TrackScores(true)
	if sortpopularity {
		query = query.Sort("popularity", false)
	}
	if sortdate {
		query = query.Sort("create_date", false)
	}
	query = query.Sort("_score", false)

	results, err := ops.Sc.es.DoSearch(ctx, query)
	if err != nil {
		return nil, 0, nil, nil, err
	}

	if results.Hits.TotalHits.Value > 0 {
		scores := make([]float64, 0, len(results.Hits.Hits))
		highlights := make([]*string, 0, len(results.Hits.Hits))
		illustids := make([]uint64, len(results.Hits.Hits))
		for i, hit := range results.Hits.Hits {
			illustids[i] = utils.Atoi(hit.Id)
			if hit.Score != nil {
				scores = append(scores, *hit.Score)
			} else {
				scores = append(scores, -1)
			}
			if highl, ok := hit.Highlight["title"]; ok {
				highlights = append(highlights, &highl[0])
			} else {
				highlights = append(highlights, nil)
			}
		}

		illusts, err := ops.QueryIllusts(ctx, illustids, resultbanned)
		if err != nil {
			return nil, 0, nil, nil, err
		}

		illustsmap := make(map[uint64]models.Illust)
		for _, illust := range illusts {
			illustsmap[illust.ID] = illust
		}

		result := make([]models.Illust, 0, len(results.Hits.Hits))
		for _, illustid := range illustids {
			if _, exist := illustsmap[illustid]; exist {
				result = append(result, illustsmap[illustid])
			}
		}

		return result, results.Hits.TotalHits.Value, scores, highlights, nil
	} else {
		return nil, 0, nil, nil, models.ErrorNoResult
	}
	return nil, 0, nil, nil, err
}

func (ops *DatabaseOperations) QueryIllustsByTags(ctx context.Context, tags []string, page int, limit int, sortpopularity bool, sortdate bool, resultbanned bool) ([]models.Illust, error) {
	/*
		{
		  "query": {
		    "nested": {
		      "path": "tags",
		      "query": {
		        "bool": {
		          "must": [
		            {
		              "match": {
		                "tags.name.fuzzy": "<TAG>"
		              }
		            },
		          ]
		        }
		      }
		    }
		  }
		}
	*/
	tagQuires := make([]elastic.Query, len(tags))
	for i, tag := range tags {
		tagQuires[i] = elastic.NewMatchQuery("tags.name.fuzzy", tag)
	}
	query := ops.Sc.es.Search(config.IllustSearchIndexName).Query(
		elastic.NewNestedQuery("tags", elastic.NewBoolQuery().Must(tagQuires...)),
	).Size(limit).From(page * limit).FetchSourceContext(elastic.NewFetchSourceContext(true).Include("_id"))

	if sortpopularity {
		query = query.Sort("popularity", false)
	}
	if sortdate {
		query = query.Sort("create_date", false)
	}

	results, err := ops.Sc.es.DoSearch(ctx, query)
	if err != nil {
		return nil, err
	}

	if results.Hits.TotalHits.Value > 0 {
		illustids := make([]uint64, len(results.Hits.Hits))
		for i, hit := range results.Hits.Hits {
			illustids[i] = utils.Atoi(hit.Id)
		}

		illusts, err := ops.QueryIllusts(ctx, illustids, resultbanned)
		if err != nil {
			return nil, err
		}

		illustsmap := make(map[uint64]models.Illust)
		for _, illust := range illusts {
			illustsmap[illust.ID] = illust
		}

		result := make([]models.Illust, 0, len(results.Hits.Hits))
		for _, illustid := range illustids {
			if _, exist := illustsmap[illustid]; exist {
				result = append(result, illustsmap[illustid])
			}
		}

		return result, nil
	}
	return nil, models.ErrorNoResult
}

func (ops *DatabaseOperations) QueryIllustByUser(ctx context.Context, userId uint64, resultbanned bool) ([]models.Illust, error) {
	var results []models.Illust
	opts := options.Find().SetSort(bson.M{"createDate": -1})
	query := bson.M{"user": userId, "banned": false}
	if resultbanned {
		query["banned"] = true
	}
	cursor, err := ops.Cols.Illust.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, err
}

func (ops *DatabaseOperations) QueryIllustByUserWithPage(ctx context.Context, userId uint64, page int64, limit int64, resultbanned bool) ([]models.Illust, error) {
	var results []models.Illust
	query := bson.M{"user": userId, "banned": false}
	if resultbanned {
		query["banned"] = true
	}
	opts := options.Find().SetSort(bson.M{"createDate": -1}).SetLimit(limit).SetSkip(page * limit)
	cursor, err := ops.Cols.Illust.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, err
}

func (ops *DatabaseOperations) DeleteIllust(ctx context.Context, illustId uint64) error {
	_, err := ops.Cols.Illust.DeleteOne(ctx, bson.M{"_id": illustId})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return err
}

func (ops *DatabaseOperations) IsIllustExist(illustId uint64) (bool, error) {
	is, err := ops.Flt.Exists(config.IllustTableName, utils.Itoa(illustId))
	if err != nil {
		return false, err
	}
	return is, nil
}
