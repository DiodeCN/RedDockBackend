package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Comment struct {
	ID        string `json:"id" bson:"id"`
	TweetID   string `json:"tweet_id" bson:"tweet_id"`
	Author    string `json:"author" bson:"author"`
	Content   string `json:"content" bson:"content"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
}

func StartCommentsApp(router *gin.Engine, twitterDatabase *mongo.Database) {
	commentsCollection := twitterDatabase.Collection("Comments")

	router.POST("/api/tweets/:tweet_id/comments", func(c *gin.Context) {
		tweetID := c.Param("tweet_id")
		var newComment Comment
		if err := c.BindJSON(&newComment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newComment.TweetID = tweetID
		newComment.Timestamp = time.Now().Unix()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := commentsCollection.InsertOne(ctx, newComment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, newComment)
	})

	router.GET("/api/tweets/:tweet_id/comments", func(c *gin.Context) {
		tweetID := c.Param("tweet_id")
		filter := bson.M{"tweet_id": tweetID}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := commentsCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var comments []Comment
		err = cursor.All(ctx, &comments)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, comments)
	})
}
