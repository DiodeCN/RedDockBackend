package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/DiodeCN/RedDockBackend/Components/inviter"
	"github.com/DiodeCN/RedDockBackend/Components/login"
	"github.com/DiodeCN/RedDockBackend/Components/register"
	"github.com/DiodeCN/RedDockBackend/Components/tweet"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 添加CORS中间件，允许来自所有域的请求
	r.Use(cors.Default())
	ctx := context.Background()

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	twitterDatabase := client.Database("RedDock")
	tweetsCollection := twitterDatabase.Collection("Tweets")
	usersCollection := twitterDatabase.Collection("Users")
	inviterCollection := client.Database("RedDock").Collection("Inviter")

	inviter.InitializeInviter(ctx, inviterCollection)

	r.GET("/api/tweets", func(c *gin.Context) {
		tweets, err := tweet.GetAllTweets(ctx, tweetsCollection)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, tweets)
	})

	r.POST("/api/login", login.HandleLogin(usersCollection))
	r.POST("/api/send_VC", register.SendVerificationCodeHandler(usersCollection))
	r.POST("/api/register", register.RegisterHandler(usersCollection, inviterCollection))

	r.GET("/api/avatar/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join(cwd, "PictureStorage", "Avatar", filename+".png")

		c.FileAttachment(filePath, filename+".png")
	})

	if err := r.Run(":10628"); err != nil {
		log.Fatal(err)
	}

}
