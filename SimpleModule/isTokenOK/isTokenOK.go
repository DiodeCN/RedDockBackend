package isTokenOK

import (
	"context"
	"log"
	"strings"

	"github.com/DiodeCN/RedDockBackend/SimpleModule/iwantatoken"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	PhoneNumber string `bson:"phoneNumber"`
	UserStatus  int    `bson:"userStatus"`
}

func TokenHandler(usersCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		secretKey := iwantatoken.GetTokenSecretKey()
		decryptedToken, err := iwantatoken.Decrypt(token, []byte(secretKey))
		if err != nil {
			c.JSON(400, gin.H{"error": "Decryption error"})
			return
		}

		phoneNumber := strings.Split(decryptedToken, "|")[0]

		var user User
		err = usersCollection.FindOne(context.Background(), bson.M{"phoneNumber": phoneNumber}).Decode(&user)
		if err != nil {
			c.JSON(400, gin.H{"error": "User not found"})
			return
		}

		switch user.UserStatus {
		case 0:
			c.Status(200)
		case 1:
			c.JSON(403, gin.H{"error": "Account banned"})
		default:
			c.JSON(400, gin.H{"error": "Invalid token"})
		}
	}
}

func IsUserValidByToken(token string, usersCollection *mongo.Collection) bool {
	secretKey := iwantatoken.GetTokenSecretKey()
	decryptedToken, err := iwantatoken.Decrypt(token, []byte(secretKey))
	if err != nil {
		return false
	}

	log.Println(decryptedToken)

	phoneNumber := strings.Split(decryptedToken, "|")[0]

	var user User
	err = usersCollection.FindOne(context.Background(), bson.M{"phoneNumber": phoneNumber}).Decode(&user)

	if err != nil || user.UserStatus != 0 {
		return false
	}

	return true
}
