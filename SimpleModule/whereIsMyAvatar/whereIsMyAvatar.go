package whereismyavatar

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AvatarHandler(usersCollection *mongo.Collection, cwd string) gin.HandlerFunc {
	return func(c *gin.Context) {
		avatarID := c.Param("filename")

		filePath := filepath.Join(cwd, "PictureStorage", "Avatar", avatarID+".png")
		c.FileAttachment(filePath, avatarID+".png")

	}
}
