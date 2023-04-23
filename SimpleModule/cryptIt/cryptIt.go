package cryptIt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
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
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < 2*aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	blockMode := cipher.NewCBCDecrypter(block, iv)
	paddedData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(paddedData, ciphertext)

	//log.Println("Padded data: ", paddedData)

	if len(paddedData)%16 != 0 {
		return "", errors.New("padded data length is not a multiple of 16")
	}

	unpaddedData, err := pkcs7.Unpad(paddedData, 16)
	if err != nil {
		return "", err
	}

	//log.Println("解密数据: ", string(unpaddedData))

	return string(unpaddedData), nil
}

func Encrypt(plaintext string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("AES_SECRET_KEY")
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	paddedText := make([]byte, len(plaintext)+padding)
	copy(paddedText, plaintext)
	for i := 0; i < padding; i++ {
		paddedText[len(plaintext)+i] = byte(padding)
	}

	ciphertext := make([]byte, aes.BlockSize+len(paddedText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedText)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
