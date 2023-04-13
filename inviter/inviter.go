package inviter

import (
	"context"
	"log"

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
		defaultInviter := bson.M{
			"inviter":    "#19890six0four",
			"userscount": 0,
		}

		_, err := inviterCollection.InsertOne(ctx, defaultInviter)
		if err != nil {
			log.Fatalf("Failed to insert default inviter: %v", err)
		}

		log.Println("Inserted default inviter: #19890six0four")
	}
}
