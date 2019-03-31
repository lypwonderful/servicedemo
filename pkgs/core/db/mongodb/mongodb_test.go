package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

type UserInfo struct {
	Name   string  `bson:"name"`
	Value1 float64 `bson:"value"`
}

func TestMongoDB_PingTest(t *testing.T) {
	dbObj := &MongoFactory{}
	dbObj.Create("mongodb://127.0.0.1:27017").NewDBClient().PingTest()
}

func TestMongoDB_Find(t *testing.T) {
	dbObj := &MongoFactory{}
	dbCli := dbObj.Create("mongodb://root:example@localhost:27017").NewDBClient()

	uInfo := bson.M{"name": "pi", "value": 3.14159}
	inseartInfo, err := dbCli.Client.Database("userTest").Collection("userinfo").InsertOne(dbCli.Ctx, uInfo)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("inseart info:", inseartInfo)
}

func TestMongoDB_Find_Index(t *testing.T) {
	dbObj := &MongoFactory{}
	dbCli := dbObj.Create("mongodb://root:example@localhost:27017").NewDBClient()

	uInfo1 := bson.M{"$set": bson.M{"value": 3.24159}}
	uInfo12 := bson.M{"name": "p3i"}
	//uInfo := bson.M{"name": "p2i", "value": 3.14159}
	c := dbCli.Client.Database("userTest").Collection("userinfounique2")
	o := options.Index()
	o.SetUnique(true)
	var trueFlag = true
	c.Indexes().CreateOne(dbCli.Ctx, mongo.IndexModel{Keys: bson.M{"name": 1}, Options: &options.IndexOptions{
		Unique: &trueFlag,
	}})

	inseartInfo, err := c.UpdateOne(dbCli.Ctx, uInfo12, uInfo1, &options.UpdateOptions{
		Upsert: &trueFlag,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("inseart info:", inseartInfo)
}

func TestFind(t *testing.T) {
	dbObj := &MongoFactory{}
	dbCli := dbObj.Create("mongodb://root:example@localhost:27017").NewDBClient()
	c := dbCli.Client.Database("userTest").Collection("userinfounique2")
	//retUserInfo := &UserInfo{}
	retUserInfo, err := c.Find(dbCli.Ctx, bson.D{})
	if err != nil {
		t.Fatal(err)
	}
	var out interface{}
	for retUserInfo.Next(dbCli.Ctx) {
		retUserInfo.Decode(&out)
		t.Log(out)
	}

}
