package register

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	hash "github.com/DiodeCN/RedDockBackend/RefactoredModule/hash"
	"github.com/DiodeCN/RedDockBackend/SimpleModule/CanSendVerificationCode"
	globalDataManipulation "github.com/DiodeCN/RedDockBackend/SimpleModule/globalDataManipulation"
	iwantatoken "github.com/DiodeCN/RedDockBackend/SimpleModule/iwantatoken"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterRequestData struct {
	Timestamp        string `json:"timestamp"`
	Nickname         string `json:"nickname"`
	Inviter          string `json:"inviter"`
	PhoneNumber      string `json:"phoneNumber"`
	VerificationCode string `json:"verificationCode"`
	Password         string `json:"password"`
}

// GenerateVerificationCode generates a random 6-digit verification code
func GenerateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func VerifyAndRegisterUser(ctx context.Context, usersCollection *mongo.Collection, inviterCollection *mongo.Collection, phoneNumber, verificationCode, nickname, inviter, password string) (bool, error) {
	//搞点哈希
	hashedPassword, err := hash.HashString(password)
	if err != nil {
		return false, fmt.Errorf("hashing_error")
	}

	// 检查邀请人是否存在
	inviterFilter := bson.M{"inviter": inviter}
	inviterDoc := bson.M{}
	err = inviterCollection.FindOne(ctx, inviterFilter).Decode(&inviterDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 邀请人不存在
			return false, fmt.Errorf("inviter_not_found")
		}
		// 其他类型的错误
		return false, fmt.Errorf("database_error")
	}

	// 邀请人存在，检查验证码
	filter := bson.M{"phoneNumber": phoneNumber}
	user := bson.M{}
	err = usersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 用户不存在，但验证码不匹配，则拒绝注册
			return false, nil
		}
		// 其他类型的错误
		return false, err
	}

	// 用户已存在，检查验证码是否正确
	if user["verificationCode"] == verificationCode {
		log.Printf("验证码匹配: 提交的验证码=%s, 数据库中的验证码=%s", verificationCode, user["verificationCode"])

		// 检查用户是否已经设置了昵称、邀请人和密码
		if user["nickname"] != nil && user["inviter"] != nil && user["password"] != nil {
			log.Printf("用户已存在，无需重新注册")
			return false, fmt.Errorf("user_already_exists")
		}

		// 更新邀请人的 userscount
		usersCount := int(inviterDoc["userscount"].(int32)) + 1
		usersCountUpdate := bson.M{"$set": bson.M{"userscount": usersCount}}
		_, err = inviterCollection.UpdateOne(ctx, inviterFilter, usersCountUpdate)
		if err != nil {
			return false, err
		}

		globalDataManipulation.IncrementUsers()
		uid := int(user["_id"].(int32))

		newUser := bson.M{
			"_id":              uid,
			"nickname":         nickname,
			"inviter":          inviter,
			"phoneNumber":      phoneNumber,
			"password":         hashedPassword, // 将原始密码替换为哈希后的密码
			"verificationCode": verificationCode,
		}

		_, err = usersCollection.InsertOne(ctx, newUser)
		if err != nil {
			return false, err
		}

		// 用户注册成功
		return true, nil
	}
	return false, fmt.Errorf("invalid_verification_code")

}

// 修改 RegisterHandler，以便根据 VerifyAndRegisterUser 返回的错误信息设置适当的响应
func RegisterHandler(usersCollection *mongo.Collection, inviterCollection *mongo.Collection) func(c *gin.Context) {
	return func(c *gin.Context) {
		var requestData RegisterRequestData

		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
			return
		}

		nickname := requestData.Nickname
		inviter := requestData.Inviter
		phoneNumber := requestData.PhoneNumber
		verificationCode := requestData.VerificationCode
		password := requestData.Password

		// Add logging statements to print the received form data
		log.Printf("Form data received: nickname=%s, inviter=%s, phoneNumber=%s, verificationCode=%s, password=%s\n", nickname, inviter, phoneNumber, verificationCode, password)

		registerSuccess, err := VerifyAndRegisterUser(context.Background(), usersCollection, inviterCollection, phoneNumber, verificationCode, nickname, inviter, password)
		if err != nil {
			var errorType string
			switch err.Error() {
			case "inviter_not_found":
				errorType = "inviter_not_found"
			case "user_already_exists":
				errorType = "user_already_exists"
			case "database_error":
				errorType = "database_error"
			case "invalid_verification_code":
				errorType = "invalid_verification_code"
			default:
				errorType = "unknown_error"
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": errorType})
			return
		}

		if registerSuccess {
			// 创建Token字符串
			tokenString := phoneNumber + "|" + requestData.Timestamp

			// 使用iwantatoken加密Token
			secretKey := []byte(iwantatoken.GetTokenSecretKey())
			encryptedToken, err := iwantatoken.Encrypt(tokenString, secretKey)
			if err != nil {
				log.Println("Error during token encryption: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Token encryption failed"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "user_registered_successfully", "token": encryptedToken})
		}
	}
}

func SendVerificationCodeHandler(usersCollection *mongo.Collection) func(c *gin.Context) {
	return func(c *gin.Context) {
		var requestData struct {
			PhoneNumber string `json:"phoneNumber"`
		}

		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		phoneNumber := requestData.PhoneNumber
		clientIP := c.ClientIP()

		// 获取用户文档
		filter := bson.M{"phoneNumber": phoneNumber}
		user := bson.M{}
		err := usersCollection.FindOne(context.Background(), filter).Decode(&user)

		if err != nil && err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		if user["daycount"] != nil && int(user["daycount"].(int32)) > 5 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "You have reached the daily limit for sending verification codes. Please try again tomorrow."})
			return
		}

		if !CanSendVerificationCode.CanSendVerificationCode(clientIP, phoneNumber) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please wait 60 seconds before requesting a new verification code."})
			return
		}

		uid, err := globalDataManipulation.GetAndIncrementUsers()
		if err != nil {
			log.Fatal(err)
		}

		// Generate a verification code
		verificationCode := GenerateVerificationCode()

		// Store the verification code and user ID in the database
		err = CanSendVerificationCode.StoreVerificationCodeInDB(context.Background(), usersCollection, phoneNumber, verificationCode, int(uid))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store verification code"})
			return
		}

		RegisterSMS([]string{phoneNumber}, []string{verificationCode})
		c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
	}
}

func RegisterSMS(phoneNumberSet []string, templateParamSet []string) {
	viper.SetConfigFile("sth.config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	secretID := viper.GetString("tencentcloud_secret_id")
	secretKey := viper.GetString("tencentcloud_secret_key")

	credential := common.NewCredential(
		secretID, secretKey,
	)

	cpf := profile.NewClientProfile()

	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	cpf.SignMethod = "HmacSHA1"

	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	request := sms.NewSendSmsRequest()

	request.SmsSdkAppId = common.StringPtr("1400811261")
	request.SignName = common.StringPtr("榆法糖一般空间")
	request.TemplateId = common.StringPtr("1761760")
	request.TemplateParamSet = common.StringPtrs(templateParamSet)
	request.PhoneNumberSet = common.StringPtrs(phoneNumberSet)
	request.SessionContext = common.StringPtr("")
	request.ExtendCode = common.StringPtr("")
	request.SenderId = common.StringPtr("")

	response, err := client.SendSms(request)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}

	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(response.Response)
	fmt.Printf("%s", b)
}
