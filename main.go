package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/DiodeCN/RedDockBackend/Components/inviter"
	"github.com/DiodeCN/RedDockBackend/Components/login"
	"github.com/DiodeCN/RedDockBackend/Components/register"
	"github.com/DiodeCN/RedDockBackend/Components/tweet"

	"github.com/DiodeCN/RedDockBackend/SimpleModule/isTokenOK"
	"github.com/DiodeCN/RedDockBackend/SimpleModule/requestlogger"
	"github.com/DiodeCN/RedDockBackend/SimpleModule/whereismyavatar"

	art "github.com/DiodeCN/RedDockBackend/RefactoredModule/ArtOfCMD"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	// 添加请求日志中间件
	r.Use(requestlogger.RequestLogger())

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

	// log.Println(iwantatoken.EncryptFile("ddd"))
	twitterDatabase := client.Database("RedDock")
	tweetsCollection := twitterDatabase.Collection("Tweets")
	usersCollection := twitterDatabase.Collection("Users")
	inviterCollection := client.Database("RedDock").Collection("Inviter")

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	art.PrintArt()

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
	r.POST("/api/token", isTokenOK.TokenHandler(usersCollection))

	r.GET("/api/avatar/:filename", whereismyavatar.AvatarHandler(usersCollection, cwd))

	if err := r.Run(":10628"); err != nil {
		log.Fatal(err)
	}

}
