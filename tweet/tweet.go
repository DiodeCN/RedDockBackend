package tweet

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
