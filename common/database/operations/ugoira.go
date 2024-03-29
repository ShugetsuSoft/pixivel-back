package operations

import (
	"context"
	"time"

	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ops *DatabaseOperations) InsertUgoira(ctx context.Context, ugoira *models.Ugoira) error {
	var err error
	is, err := ops.Flt.Exists(config.UgoiraTableName, utils.Itoa(ugoira.ID))
	if err != nil {
		return err
	}
	ugoira.UpdateTime = time.Now()

	if is {
		goto REPLACE
	} else {
		_, err = ops.Cols.Ugoira.InsertOne(ctx, ugoira)

		if mongo.IsDuplicateKeyError(err) {
			_, err = ops.Flt.Add(config.UgoiraTableName, utils.Itoa(ugoira.ID))
			if err != nil {
				return err
			}
			goto REPLACE
		}

		if err != nil {
			return err
		}

		_, err = ops.Flt.Add(config.UgoiraTableName, utils.Itoa(ugoira.ID))
		if err != nil {
			return err
		}
	}

	return nil

REPLACE:
	result, err := ops.Cols.Ugoira.ReplaceOne(ctx, bson.M{"_id": ugoira.ID}, ugoira)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		_, err = ops.Cols.Ugoira.InsertOne(ctx, ugoira)
		if err != nil {
			return err
		}
		_, err = ops.Flt.Add(config.UgoiraTableName, utils.Itoa(ugoira.ID))
	}

	return err
}

func (ops *DatabaseOperations) QueryUgoira(ctx context.Context, ugoiraId uint64) (*models.Ugoira, error) {
	is, err := ops.Flt.Exists(config.UgoiraTableName, utils.Itoa(ugoiraId))

	if err != nil {
		return nil, err
	}

	if is {
		result := models.Ugoira{
			Frames: []models.UgoiraFrame{},
		}
		query := bson.M{"_id": ugoiraId}
		err := ops.Cols.Ugoira.FindOne(ctx, query).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			} else {
				return nil, err
			}
		}

		return &result, err
	}
	return nil, nil
}

func (ops *DatabaseOperations) IsUgoiraExist(ugoiraId uint64) (bool, error) {
	is, err := ops.Flt.Exists(config.UgoiraTableName, utils.Itoa(ugoiraId))
	if err != nil {
		return false, err
	}
	return is, nil
}

func (ops *DatabaseOperations) DeleteUgoira(ctx context.Context, ugoiraId uint64) error {
	is, err := ops.Flt.Exists(config.UgoiraTableName, utils.Itoa(ugoiraId))

	if err != nil {
		return err
	}

	if is {
		_, err := ops.Cols.Ugoira.DeleteOne(ctx, bson.M{"_id": ugoiraId})
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
