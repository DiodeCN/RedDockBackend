package iwantatoken

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
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getTokenSecretKey() string {
	return os.Getenv("TOKEN_SECRET_KEY")
}

func encrypt(plainText []byte, secretKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

func decrypt(cipherText []byte, secretKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

func EncryptFile(inputFilePath, outputFilePath string) error {
	secretKey := []byte(getTokenSecretKey())

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	plainText, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	cipherText, err := encrypt(plainText, secretKey)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	encoder := base64.NewEncoder(base64.StdEncoding, outputFile)
	defer encoder.Close()

	_, err = encoder.Write(cipherText)
	return err
}

func DecryptFile(inputFilePath, outputFilePath string) error {
	secretKey := []byte(getTokenSecretKey())

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	cipherText, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, inputFile))
	if err != nil {
		return err
	}

	plainText, err := decrypt(cipherText, secretKey)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.Write(plainText)
	return err
}
