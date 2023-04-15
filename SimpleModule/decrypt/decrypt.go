package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log"
	"os"
)

func Decrypt(encryptedData string) (string, error) {
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
	unpaddedData := removePadding(ciphertext)

	return string(unpaddedData), nil
}

func removePadding(data []byte) []byte {
	length := len(data)
	padding := int(data[length-1])

	return data[:(length - padding)]
}
