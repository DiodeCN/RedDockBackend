package user

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `json:"id" bson:"id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Username  string `json:"username" bson:"username"`
	Password  string `json:"-" bson:"password"` // 不返回密码，仅在数据库中存储
}

func NewUser(id, firstName, lastName, email, username, password string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Username:  username,
		Password:  hashedPassword,
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateUserPassword(user *User, password string) bool {
	return checkPasswordHash(password, user.Password)
}

func GetAllUsers(ctx context.Context, usersCollection *mongo.Collection) ([]User, error) {
	reqCtx, reqCancel := context.WithCancel(ctx)
	defer reqCancel()

	// Check if the collection is empty
	count, err := usersCollection.CountDocuments(reqCtx, bson.D{})
	if err != nil {
		return nil, err
	}

	// If the collection is empty, insert a default user
	if count == 0 {
		defaultUser, err := NewUser(
			"default1",
			"Default",
			"User",
			"default@example.com",
			"default_user",
			"password",
		)
		if err != nil {
			return nil, err
		}

		_, err = usersCollection.InsertOne(reqCtx, defaultUser)
		if err != nil {
			return nil, err
		}
	}

	cur, err := usersCollection.Find(reqCtx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cur.Close(reqCtx); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	var users []User
	for cur.Next(reqCtx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
