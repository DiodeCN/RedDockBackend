package register

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GenerateVerificationCode generates a random 6-digit verification code
func GenerateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// StoreVerificationCodeInDB stores the verification code in the MongoDB database
func StoreVerificationCodeInDB(ctx context.Context, usersCollection *mongo.Collection, phoneNumber, verificationCode string) error {
	filter := bson.M{"phoneNumber": phoneNumber}
	update := bson.M{"$set": bson.M{"verificationCode": verificationCode, "createdAt": time.Now()}}
	opts := options.Update().SetUpsert(true)

	_, err := usersCollection.UpdateOne(ctx, filter, update, opts)
	return err
}

// VerifyAndRegisterUser verifies the provided verification code and registers the user if it matches
func VerifyAndRegisterUser(ctx context.Context, usersCollection *mongo.Collection, phoneNumber, verificationCode, nickname, inviter, password string) (bool, error) {
	filter := bson.M{"phoneNumber": phoneNumber}
	user := bson.M{}
	err := usersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return false, err
	}

	if user["verificationCode"] == verificationCode {
		update := bson.M{"$set": bson.M{"nickname": nickname, "inviter": inviter, "password": password}}
		_, err = usersCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func RegisterHandler(usersCollection *mongo.Collection) func(c *gin.Context) {
	return func(c *gin.Context) {
		nickname := c.PostForm("nickname")
		inviter := c.PostForm("inviter")
		phoneNumber := c.PostForm("phoneNumber")
		verificationCode := c.PostForm("verificationCode")
		password := c.PostForm("password")

		// Add logging statements to print the received form data
		log.Printf("Form data received: nickname=%s, inviter=%s, phoneNumber=%s, verificationCode=%s, password=%s\n", nickname, inviter, phoneNumber, verificationCode, password)

		registerSuccess, err := VerifyAndRegisterUser(context.Background(), usersCollection, phoneNumber, verificationCode, nickname, inviter, password)
		if err != nil {
			// Add logging statement to print the error
			log.Printf("Error in VerifyAndRegisterUser: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		if registerSuccess {
			c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
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
		verificationCode := GenerateVerificationCode()

		err := StoreVerificationCodeInDB(context.Background(), usersCollection, phoneNumber, verificationCode)
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
