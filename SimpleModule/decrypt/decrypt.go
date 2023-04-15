package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log"

	"github.com/joho/godotenv"

	"os"
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

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove padding
	unpaddedData, err := removePadding(ciphertext)
	if err != nil {
		return "", err
	}
	log.Println("Decrypted data: ", string(unpaddedData)) // 添加日志输出

	return string(unpaddedData), nil
}

func removePadding(data []byte) ([]byte, error) {
	length := len(data)
	padding := int(data[length-1])

	if padding < 1 || padding > aes.BlockSize {
		return nil, errors.New("invalid padding value")
	}

	return data[:(length - padding)], nil
}
