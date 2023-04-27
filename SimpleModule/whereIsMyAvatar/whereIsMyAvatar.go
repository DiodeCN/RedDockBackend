package whereismyavatar

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AvatarHandler(usersCollection *mongo.Collection, cwd string) gin.HandlerFunc {
	return func(c *gin.Context) {
		avatarID := c.Param("filename")

		filePath := filepath.Join(cwd, "PictureStorage", "Avatar", avatarID+".png")

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
			return
		}

		c.FileAttachment(filePath, avatarID+".png")
	}
}

