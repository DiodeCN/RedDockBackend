package login

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	crypt "github.com/DiodeCN/RedDockBackend/SimpleModule/crypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

type LoginData struct {
	Timestamp string `json:"timestamp"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Encrypted string `json:"encrypted"`
}

func HandleLogin(usersCollection *mongo.Collection) func(c *gin.Context) {
	return func(c *gin.Context) {
		var loginData LoginData

		err := json.NewDecoder(c.Request.Body).Decode(&loginData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}
		log.Println(loginData)

		encrypted, err := crypt.Encrypt("nihao")
		if err != nil {
			fmt.Println("Encryption error:", err)
			return
		}
		fmt.Println("Encrypted data:", encrypted)

		decrypted, err := crypt.Decrypt(loginData.Encrypted)
		if err != nil {
			log.Println("Error during decryption: ", err) // 添加日志输出
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failed"})
			return
		}
		log.Println(decrypted)

		decryptedParts := strings.Split(decrypted, "|")
		if len(decryptedParts) != 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid decrypted data"})
			return
		}

		// 删除以下代码段
		// if decryptedParts[0] == loginData.Timestamp && decryptedParts[1] == loginData.Email && decryptedParts[2] == loginData.Password {
		// 	log.Println("牛逼")
		// 	c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
		// } else {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		// }

		// 查询用户
		filter := bson.M{"phoneNumber": loginData.Email}
		user := bson.M{}
		err = usersCollection.FindOne(context.Background(), filter).Decode(&user)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// 验证密码以及解密后的信息
		if user["password"] == loginData.Password && decryptedParts[0] == loginData.Timestamp && decryptedParts[1] == loginData.Email && decryptedParts[2] == loginData.Password {
			c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password or credentials"})
		}
	}
}
