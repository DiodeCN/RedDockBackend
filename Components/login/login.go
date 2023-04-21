package login

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	crypt "github.com/DiodeCN/RedDockBackend/SimpleModule/cryptIt"
	"github.com/DiodeCN/RedDockBackend/SimpleModule/iwantatoken"
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

		/*
			encrypted, err := crypt.Encrypt("nihao")
			if err != nil {
				fmt.Println("Encryption error:", err)
				return
			}
			fmt.Println("Encrypted data:", encrypted)
		*/

		decrypted, err := crypt.Decrypt(loginData.Encrypted)
		if err != nil {
			log.Println("Error during decryption: ", err) // 添加日志输出
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failed"})
			return
		}

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
			// 创建Token字符串
			tokenString := loginData.Email + "|" + loginData.Timestamp

			// 使用iwantatoken加密Token
			secretKey := []byte(iwantatoken.GetTokenSecretKey())
			encryptedToken, err := iwantatoken.Encrypt(tokenString, secretKey)
			if err != nil {
				log.Println("Error during token encryption: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Token encryption failed"})
				return
			}
/*
			decryptedToken, err := iwantatoken.Decrypt("3r8AZJTLXjVulZv4L03PYaIgChr/blFzhrspkIEveH0ZS38V1jLsxo8mmwAggxJTlXNHSTalWQ==", secretKey)
			if err != nil {
				log.Println("Error during token encryption: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Token encryption failed"})
				return
			}
			log.Println(decryptedToken)
*/
			// 返回登录成功信息和加密Token
			c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": encryptedToken})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password or credentials"})
		}
	}
}
