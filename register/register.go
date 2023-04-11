package register

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegistrationData struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

func RegisterHandler(c *gin.Context) {
	var registrationData RegistrationData

	if err := c.ShouldBindJSON(&registrationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 在这里处理用户注册逻辑，例如存储用户数据、验证等。

	// 发送短信验证码
	sendSMSVerification(registrationData.PhoneNumber)

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
	})
}

func sendSMSVerification(phoneNumber string) {
	// 在这里调用您之前提供的短信发送代码
}
