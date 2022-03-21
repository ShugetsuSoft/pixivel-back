package drivers

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	Context context.Context
	Client  *mongo.Client
	DB      *mongo.Database
}

func NewMongoDatabase(context context.Context, uri string, dbname string) (*MongoDatabase, error) {
	clientOptions := options.Client().ApplyURI(uri).SetMaxPoolSize(20)
	client, err := mongo.Connect(context, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context, nil)
	if err != nil {
		return nil, err
	}

	return &MongoDatabase{
		Context: context,
		Client:  client,
		DB:      client.Database(dbname),
	}, nil
}

func (db *MongoDatabase) Session() (mongo.Session, error) {
	session, err := db.Client.StartSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (db *MongoDatabase) Collection(colname string) *mongo.Collection {
	return db.DB.Collection(colname)
}

func (db *MongoDatabase) Close() error {
	return db.Client.Disconnect(db.Context)
}
