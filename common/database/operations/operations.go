package operations

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ShugetsuSoft/pixivel-back/common/database/drivers"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseCollections struct {
	Illust *mongo.Collection
	User   *mongo.Collection
	Rank   *mongo.Collection
	Ugoira *mongo.Collection
}

type DatabaseOperations struct {
	Flt  models.Filter
	Cols *DatabaseCollections
	Sc   *SearchOperations
}

type SearchOperations struct {
	es  *drivers.ElasticSearch
	ndb *drivers.NearDB
}

func NewDatabaseOperations(ctx context.Context, db *drivers.MongoDatabase, filter models.Filter, es *drivers.ElasticSearch, ndb *drivers.NearDB) *DatabaseOperations {
	var err error
	if es != nil {
		err = es.CreateIndex(ctx, config.IllustSearchIndexName, models.IllustSearchMapping)
		if err != nil && err != models.ErrorIndexExist {
			log.Fatal(err)
		}
		err = es.CreateIndex(ctx, config.UserSearchIndexName, models.UserSearchMapping)
		if err != nil && err != models.ErrorIndexExist {
			log.Fatal(err)
		}
	}

	illustCol := db.Collection(config.IllustTableName)
	userCol := db.Collection(config.UserTableName)
	rankCol := db.Collection(config.RankTableName)
	ugoiraCol := db.Collection(config.UgoiraTableName)
	rankCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "date", Value: -1}, {Key: "mode", Value: -1}, {Key: "content", Value: -1}},
		Options: options.Index().SetUnique(true),
	})
	return &DatabaseOperations{
		Flt: filter,
		Cols: &DatabaseCollections{
			Illust: illustCol,
			User:   userCol,
			Rank:   rankCol,
			Ugoira: ugoiraCol,
		},
		Sc: &SearchOperations{
			es:  es,
			ndb: ndb,
		},
	}
}
