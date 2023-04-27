package getuserinfo

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserInfo struct {
	PhoneNumber string `bson:"phoneNumber"`
	Nickname    string `bson:"nickname"`
}

func GetUserInfoHandler(usersCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userid")

		// Convert string to ObjectID
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Find user in the database
		filter := bson.M{"_id": objID}
		var userInfo UserInfo
		err = usersCollection.FindOne(context.Background(), filter).Decode(&userInfo)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(404, gin.H{"error": "User not found"})
			} else {
				c.JSON(500, gin.H{"error": "Internal server error"})
			}
			return
		}

		// Return the requested user info
		c.JSON(200, userInfo)
	}
}
