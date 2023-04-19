package globaldatamanipulation

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func IncrementUsers() error {
	// 连接到MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// 获取RedDock.Global文档
	collection := client.Database("RedDock").Collection("Global")
	filter := bson.M{"_id": "RedDock.Global"}
	update := bson.M{"$inc": bson.M{"Total number of users": 1}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func GetAndIncrementUsers() (int32, error) {
	// 连接到MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// 获取并更新RedDock.Global文档中的“Total number of users”字段
	collection := client.Database("RedDock").Collection("Global")
	filter := bson.M{"_id": "RedDock.Global"}
	update := bson.M{"$inc": bson.M{"Total number of users": 1000000}}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
	var result bson.M
	err = collection.FindOneAndUpdate(context.Background(), filter, update, options).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	// 返回更新后的“Total number of users”字段值
	return result["Total number of users"].(int32), nil
}
