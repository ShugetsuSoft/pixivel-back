package operations

import (
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ops *DatabaseOperations) InsertRank(mode string, date string, content string) (bool, error) {
	rank := &models.Rank{
		Date:    date,
		Mode:    mode,
		Content: content,
		Illusts: []models.RankIllust{},
	}
	_, err := ops.Cols.Rank.InsertOne(ops.Ctx, rank)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ops *DatabaseOperations) AddRankIllusts(mode string, date string, content string, illusts []models.RankIllust) error {
	filter := bson.M{"date": date, "mode": mode, "content": content}
	update := bson.M{
		"$push": bson.M{
			"illusts": bson.M{
				"$each": illusts,
				"$sort": bson.M{"rank": 1},
			},
		},
	}
	_, err := ops.Cols.Rank.UpdateOne(ops.Ctx, filter, update)
	return err
}

func (ops *DatabaseOperations) QueryRankIllusts(mode string, date string, content string, page int, limit int) ([]models.RankAggregateResult, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"date", date},
			{"mode", mode},
			{"content", content},
		}}},
		{{"$unwind", bson.D{
			{"path", "$illusts"},
		}}},
		{{"$skip", page * limit}},
		{{"$limit", limit}},
		{{"$lookup", bson.D{
			{"from", "Illust"},
			{"localField", "illusts.id"},
			{"foreignField", "_id"},
			{"as", "illusts"},
		}}},
		{{"$project", bson.D{
			{"content", bson.D{{"$arrayElemAt", bson.A{"$illusts", 0}}}},
			{"_id", 0},
		}}},
	}
	cursor, err := ops.Cols.Rank.Aggregate(ops.Ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var results []models.RankAggregateResult
	if err = cursor.All(ops.Ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (ops *DatabaseOperations) GetSampleIllusts(quality int, limit int, resultbanned bool) ([]models.Illust, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"popularity", bson.D{{"$gt", quality}}},
			{"type", 0},
			{"banned", resultbanned},
		}}},
		{{"$sample", bson.D{
			{"size", limit},
		}}},
	}
	cursor, err := ops.Cols.Illust.Aggregate(ops.Ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var results []models.Illust
	if err = cursor.All(ops.Ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (ops *DatabaseOperations) GetSampleUsers(limit int, resultbanned bool) ([]models.User, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"banned", resultbanned},
		}}},
		{{"$sample", bson.D{
			{"size", limit},
		}}},
	}
	cursor, err := ops.Cols.User.Aggregate(ops.Ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var results []models.User
	if err = cursor.All(ops.Ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
