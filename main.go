package main

import (
	"context"
	"log"
	"time"

	"github.com/DiodeCN/RedDockBackend/tweet"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := gin.Default()

	// 添加CORS中间件，允许来自所有域的请求
	r.Use(cors.Default())

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	twitterDatabase := client.Database("RedDock")
	tweetsCollection := twitterDatabase.Collection("Tweets")

	if err != nil {
		log.Fatal(err)
	}

	r.GET("/api/tweets", func(c *gin.Context) {
		tweets, err := tweet.GetAllTweets(ctx, tweetsCollection)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, tweets)

	})

	r.Run() // 默认监听8080端口
}
