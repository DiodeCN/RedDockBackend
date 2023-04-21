package iwantatoken

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetTokenSecretKey() string {
	return os.Getenv("TOKEN_SECRET_KEY")
}

func Encrypt(plainText string, secretKey []byte) (string, error) {
	plainTextBytes := []byte(plainText)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTextBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(cipherText string, secretKey []byte) (string, error) {
	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	if len(cipherTextBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherTextBytes, cipherTextBytes)

	return string(cipherTextBytes), nil
}

func CheckForDelimiter(encodedText string) (bool, error) {
	secretKey := GetTokenSecretKey()
	decodedText, err := Decrypt(encodedText, []byte(secretKey))
	if err != nil {
		return false, err
	}

	// Extract the first item before the delimiter
	parts := strings.Split(decodedText, "|")

	if len(parts) == 0 {
		return false, errors.New("no parts found in the decoded text")
	}
	phoneNumber := parts[0]

	// Connect to MongoDB and check if the phoneNumber exists.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return false, err
	}
	defer client.Disconnect(ctx)

	usersCollection := client.Database("RedDock").Collection("Users")

	// Check if the phoneNumber exists.
	filter := bson.M{"phoneNumber": phoneNumber}
	var result bson.M
	err = usersCollection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// If phoneNumber exists, increment the requestCount.
	update := bson.M{"$inc": bson.M{"requestCount": 1}}
	_, err = usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}

	return true, nil
}

