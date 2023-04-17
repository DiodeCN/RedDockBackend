package inviter

import (
	"context"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeInviter(ctx context.Context, inviterCollection *mongo.Collection) {
	// 检查是否有任何邀请人文档
	count, err := inviterCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to count inviter documents: %v", err)
	}

	// 如果没有邀请人文档，创建一个新的邀请人
	if count == 0 {
		// 生成一个随机的八位字符串
		rand.Seed(time.Now().UnixNano())
		inviterID := "#" + randString(8)

		defaultInviter := bson.M{
			"inviter":    inviterID,
			"userscount": 0,
		}

		_, err := inviterCollection.InsertOne(ctx, defaultInviter)
		if err != nil {
			log.Fatalf("Failed to insert default inviter: %v", err)
		}

		log.Println("Inserted default inviter:", inviterID)
	}
}

func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
