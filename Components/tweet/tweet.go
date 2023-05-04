package tweet

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tweet struct {
	UID       string    `json:"uid" bson:"uid"`
	User      string    `json:"user" bson:"user"`
	Name      string    `json:"name" bson:"name"`
	AvatarURL string    `json:"avatar_url" bson:"avatar_url"`
	Time      time.Time `json:"time" bson:"time"`
	Day       string    `json:"day" bson:"day"`
	Hour      int       `json:"hour" bson:"hour"`
	Minute    int       `json:"minute" bson:"minute"`
	Compid    string    `json:"compid" bson:"compid"`
	Dayid     string    `json:"dayid" bson:"dayid"`
	Content   string    `json:"content" bson:"content"`
	Likes     int       `json:"likes" bson:"likes"`
	Favorites int       `json:"favorites" bson:"favorites"`
	Retweets  int       `json:"retweets" bson:"retweets"`
	Shares    int       `json:"shares" bson:"shares"`
	Views     int       `json:"views" bson:"views"`
	Comments  int       `json:"comments" bson:"comments"`
	Sign      string    `json:"sign" bson:"sign"`
	Classification string `json:"classification" bson:"classification"`
	SenderUID string `json:"senderUID" bson:"senderUID"`
}

type PostTweet struct {
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}


func NewTweet(id, name, avatarURL, content string) *Tweet {
	return &Tweet{
		UID:       id,
		Name:      name,
		AvatarURL: avatarURL,
		Content:   content,
	}
}

func (t *Tweet) UpdateLikes(n int) {
	t.Likes += n
}

func (t *Tweet) UpdateFavorites(n int) {
	t.Favorites += n
}

func GetAllTweets(ctx context.Context, tweetsCollection *mongo.Collection) ([]Tweet, error) {
	reqCtx, reqCancel := context.WithCancel(ctx)
	defer reqCancel()

	// ... 执行其他操作

	// 在需要取消操作时调用 reqCancel()

	// Check if the collection is empty
	count, err := tweetsCollection.CountDocuments(reqCtx, bson.D{})
	if err != nil {
		return nil, err
	}

	// If the collection is empty, insert default tweets
	if count == 0 {
		defaultTweets := []Tweet{
			{
				UID:     "1",
				User:    "100004",
				Content: "This is a default tweet from Default User 1.",
			},
			{
				UID:     "2",
				User:    "100004",
				Content: "This is another default tweet from Default User 2.",
			},
		}

		for _, tweet := range defaultTweets {
			_, err = tweetsCollection.InsertOne(reqCtx, tweet)
			if err != nil {
				return nil, err
			}
		}
	}

	cur, err := tweetsCollection.Find(reqCtx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cur.Close(reqCtx); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	var tweets []Tweet
	for cur.Next(reqCtx) {
		var tweet Tweet
		err := cur.Decode(&tweet)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return tweets, nil
}

func PostTweetHandler(tweetsCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newTweet Tweet
		if err := c.ShouldBindJSON(&newTweet); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		newTweet.Timestamp = time.Now()

		_, err := tweetsCollection.InsertOne(context.Background(), newTweet)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error inserting tweet into database"})
			return
		}

		c.JSON(200, gin.H{"success": "Tweet successfully posted"})
	}
}