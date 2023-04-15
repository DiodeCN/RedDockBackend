package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log"

	"os"

	"github.com/joho/godotenv"
	"github.com/mergermarket/go-pkcs7"
)

func Decrypt(encryptedData string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("AES_SECRET_KEY")

	log.Println(secretKey)

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Use PKCS#5 padding
	blockMode := cipher.NewCBCDecrypter(block, iv)
	paddedData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(paddedData, ciphertext)
	unpaddedData, err := pkcs7.Unpad(paddedData, 16) // 修改为8
	if err != nil {
		return "", err
	}

	log.Println("Decrypted data: ", string(unpaddedData)) // 添加日志输出

	return string(unpaddedData), nil
}
