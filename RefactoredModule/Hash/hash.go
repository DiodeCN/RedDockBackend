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
	hashed, err := bcrypt.GenerateFromPassword([]byte(key+str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// 解密函数
func CheckHash(hashed, str string) bool {
	key := os.Getenv("HASH_KEY")
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(key+str))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
