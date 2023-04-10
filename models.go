package main

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
