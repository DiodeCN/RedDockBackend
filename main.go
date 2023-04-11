package main

import (
	"context"
	"log"
	"net/http"

	"github.com/DiodeCN/RedDockBackend/register"
	"github.com/DiodeCN/RedDockBackend/tweet"
	"github.com/DiodeCN/RedDockBackend/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
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

	r.GET("/api/tweets", func(c *gin.Context) {
		tweets, err := tweet.GetAllTweets(ctx, tweetsCollection)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, tweets)
	})

	r.POST("/login", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		user, authenticated := user.AuthenticateUser(c.Request.Context(), usersCollection, email, password)
		if !authenticated {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// 登录成功，将用户信息返回到前端（注意：不要返回密码）
		sanitizedUser := user.Sanitize()
		c.JSON(http.StatusOK, sanitizedUser)
	})

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}

	register.JustSendASMS()
}
