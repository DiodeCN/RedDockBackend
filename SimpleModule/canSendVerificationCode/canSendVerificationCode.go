package CanSendVerificationCode

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ipPhoneTimers = make(map[string]time.Time)

func CanSendVerificationCode(ip string, phoneNumber string) bool {
	key := ip + ":" + phoneNumber
	lastSendTime, exists := ipPhoneTimers[key]

	if exists && time.Since(lastSendTime) < 60*time.Second {
		return false
	}

	ipPhoneTimers[key] = time.Now()
	return true
}

func StoreVerificationCodeInDB(ctx context.Context, usersCollection *mongo.Collection, phoneNumber, verificationCode string, uid int) error {
	filter := bson.M{"phoneNumber": phoneNumber}

	// 获取用户文档
	user := bson.M{}
	err := usersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// 获取当前日期中的日
	now := time.Now()
	currentDay := now.Day()

	// 更新 sendvcday 和 daycount
	sendvcday := currentDay
	daycount := 1

	if user["sendvcday"] != nil && int(user["sendvcday"].(int32)) == currentDay {
		// sendvcday 与当前日期相同
		daycount = int(user["daycount"].(int32)) + 1
	}

	update := bson.M{
		"$set": bson.M{
			"verificationCode": verificationCode,
			"createdAt":        time.Now(),
			"sendvcday":        sendvcday,
			"daycount":         daycount,
			"_id":              uid,
			"phoneNumber":      phoneNumber,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err = usersCollection.UpdateOne(ctx, filter, update, opts)
	return err
}
