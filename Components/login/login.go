package login

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/DiodeCN/RedDockBackend/SimpleModule/decrypt"

	"github.com/gin-gonic/gin"
)

type LoginData struct {
	Timestamp string `json:"timestamp"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Encrypted string `json:"encrypted"`
}

func HandleLogin(c *gin.Context) {
	var loginData LoginData
	log.Println(loginData)

	err := json.NewDecoder(c.Request.Body).Decode(&loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	decrypted, err := decrypt.Decrypt(loginData.Encrypted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failed"})
		return
	}

	decryptedParts := strings.Split(decrypted, "|")
	if len(decryptedParts) != 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid decrypted data"})
		return
	}

	if decryptedParts[0] == loginData.Timestamp && decryptedParts[1] == loginData.Email && decryptedParts[2] == loginData.Password {
		log.Println("牛逼")
		c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
