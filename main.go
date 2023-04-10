package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tweet struct {
	ID             string `json:"id" bson:"id"`
	Name           string `json:"name" bson:"name"`
	AvatarURL      string `json:"avatar_url" bson:"avatar_url"`
	HoursSincePost int    `json:"hours_since_post" bson:"hours_since_post"`
	Content        string `json:"content" bson:"content"`
	Likes          int    `json:"likes" bson:"likes"`
	Favorites      int    `json:"favorites" bson:"favorites"`
	Retweets       int    `json:"retweets" bson:"retweets"`
	Shares         int    `json:"shares" bson:"shares"`
	Views          int    `json:"views" bson:"views"`
	Comments       int    `json:"comments" bson:"comments"`
}

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
		reqCtx, reqCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer reqCancel()

		// Check if the collection is empty
		count, err := tweetsCollection.CountDocuments(reqCtx, bson.D{})
		if err != nil {
			log.Fatal(err)
		}

		// If the collection is empty, insert a new "helloworld" tweet
		if count == 0 {
			helloTweet := Tweet{
				ID:             "你好世界！",
				Name:           "他妈的",
				AvatarURL:      "helloworld",
				HoursSincePost: 0,
				Content:        "如果你看到这个东西，说明数据库已经被remade了。",
				Likes:          0,
				Favorites:      0,
				Retweets:       0,
				Shares:         0,
				Views:          0,
				Comments:       0,
			}

			_, err = tweetsCollection.InsertOne(reqCtx, helloTweet)
			if err != nil {
				log.Fatal(err)
			}
		}

		cur, err := tweetsCollection.Find(reqCtx, bson.D{})
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(reqCtx)

		var tweets []Tweet
		for cur.Next(reqCtx) {
			var tweet Tweet
			err := cur.Decode(&tweet)
			if err != nil {
				log.Fatal(err)
			}
			tweets = append(tweets, tweet)
		}
		if err := cur.Err(); err != nil {

			log.Fatal(err)
		}

		c.JSON(200, tweets)
	})

	r.Run() // 默认监听8080端口
}
