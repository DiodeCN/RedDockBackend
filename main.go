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

	"github.com/DiodeCN/RedDockBackend/SimpleModule/iwantatoken"
	"github.com/DiodeCN/RedDockBackend/SimpleModule/requestlogger"
	"github.com/DiodeCN/RedDockBackend/SimpleModule/whereismyavatar"

	"github.com/DiodeCN/RedDockBackend/RefactoredModule/getuserinfo"
	initprint "github.com/DiodeCN/RedDockBackend/RefactoredModule/initprint"

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

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowHeaders = append(config.AllowHeaders, "")

	r := gin.Default()
	r.Use(cors.New(config))
	r.Use(requestlogger.RequestLogger())

	rt := gin.Default()
	rt.Use(cors.New(config))
	rt.Use(requestlogger.RequestLogger())



	rt.Use(iwantatoken.TokenMiddleware()) // 使用TokenMiddleware

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

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initprint.PrintArt()

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
	r.POST("/api/tokencheck", iwantatoken.TokenHandler(usersCollection))

	rt.POST("/api/posttweet", tweet.PostTweetHandler(tweetsCollection))
	rt.GET("/api/avatar/:filename", whereismyavatar.AvatarHandler(usersCollection, cwd))
	
	rt.GET("/api/userinfo/:userid", getuserinfo.GetUserInfoHandler(usersCollection))

	// Start rt router on a separate goroutine
	go func() {
		if err := rt.Run(":10629"); err != nil {
			log.Fatal(err)
		}
	}()

	// Start r router on the main goroutine
	if err := r.Run(":10628"); err != nil {
		log.Fatal(err)
	}
}
