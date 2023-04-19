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

	// 查询RedDock.Global文档中的“Total number of users”字段的当前值
	collection := client.Database("RedDock").Collection("Global")
	filter := bson.M{"_id": "RedDock.Global"}
	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 如果文档不存在，则创建一个新文档并将“Total number of users”字段设置为1000001
			_, err = collection.InsertOne(context.Background(), bson.M{"_id": "RedDock.Global", "Total number of users": 100000})
			if err != nil {
				log.Fatal(err)
			}
			return nil
		}
		log.Fatal(err)
	}

	// 根据需要增加“Total number of users”字段的值
	currentUsers := result["Total number of users"].(int32)
	if currentUsers == 0 {
		// 如果当前值为0，则将其设置为1000001
		update := bson.M{"$set": bson.M{"Total number of users": 100000}}
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// 否则，增加其值
		update := bson.M{"$inc": bson.M{"Total number of users": 1}}
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
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
	options := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
	var result bson.M
	err = collection.FindOneAndUpdate(context.Background(), filter, options).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	// 返回更新后的“Total number of users”字段值
	return result["Total number of users"].(int32), nil
}
