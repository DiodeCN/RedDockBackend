package getuserinfo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserInfo struct {
	PhoneNumber string `bson:"phoneNumber"`
	Nickname    string `bson:"nickname"`
}

func GetUserInfoHandler(usersCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		hexUserID := c.Param("userid")

		// Convert hex string to int32
		decodedUserID, err := strconv.ParseInt(hexUserID, 16, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID format"})
			return
		}
		userID := int32(decodedUserID)

		fmt.Printf("UserID: %d\n", userID)

		// Find user in the database
		filter := bson.M{"_id": userID}
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
