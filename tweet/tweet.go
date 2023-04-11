package tweet

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func NewTweet(id, name, avatarURL, content string) *Tweet {
	return &Tweet{
		ID:        id,
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

func GetTweets(tweetsCollection *mongo.Collection) []Tweet {
	reqCtx, reqCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer reqCancel()

	// Check if the collection is empty
	count, err := tweetsCollection.CountDocuments(reqCtx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	// If the collection is empty, insert a new "helloworld" tweet
	if count == 0 {
		helloTweet := NewTweet(
			"你好世界！",
			"他妈的",
			"helloworld",
			"如果你看到这个东西，说明数据库已经被remade了。",
		)

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

	return tweets
}

func GetAllTweets(ctx context.Context, tweetsCollection *mongo.Collection) ([]Tweet, error) {
	reqCtx, reqCancel := context.WithTimeout(ctx, 5*time.Second)
	defer reqCancel()

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
