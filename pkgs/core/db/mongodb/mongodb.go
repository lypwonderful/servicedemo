package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoDB struct {
	Client      *mongo.Client
	Ctx         context.Context
	ServiceAddr string
}

type MongoDBer interface {
	NewDBClient() *MongoDB
}

type MongoFactory struct{}

func (MongoFactory) Create(sAddr string) MongoDBer {
	return &MongoDB{
		ServiceAddr: sAddr,
	}
}
func (mongoDB *MongoDB) NewDBClient() *MongoDB {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDB.ServiceAddr))
	if err != nil {
		panic(err)
	}
	return &MongoDB{
		Client: client,
		Ctx:    ctx,
	}
}

func (mongoDB *MongoDB) PingTest() {
	err := mongoDB.Client.Ping(mongoDB.Ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
}
