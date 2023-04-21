package hash

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

// 加密函数
func HashString(str string) (string, error) {
	key := os.Getenv("HASH_KEY")

	if len(key) == 0 {
		return "", fmt.Errorf("Hash key not found in .env file")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(key+str), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// 解密函数
func CheckHash(hashed, str string) bool {
	key := os.Getenv("HASH_KEY")

	if len(key) == 0 {
		fmt.Println("Hash key not found in .env file")
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(key+str))

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
