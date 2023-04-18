package whereismyavatar

import (
	"path/filepath"

	"github.com/DiodeCN/RedDockBackend/SimpleModule/isTokenOK"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AvatarHandler(usersCollection *mongo.Collection, cwd string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		avatarID := c.Param("filename")

		if isTokenOK.IsUserValidByToken(token, usersCollection) {
			filePath := filepath.Join(cwd, "PictureStorage", "Avatar", avatarID+".png")
			c.FileAttachment(filePath, avatarID+".png")
		} else {
			c.JSON(403, gin.H{"error": "Invalid token"})
		}
	}
}
