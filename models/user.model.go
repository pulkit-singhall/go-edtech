package models

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	First_name    string             `json:"first_name" validate:"required,min=2,max=20"`
	Last_name     string             `json:"last_name" validate:"required,min=2,max=20"`
	Email         string             `json:"email" validate:"email,required"`
	Password      string             `json:"password" validate:"required,min=8,max=16"`
	Token         string             `json:"token"`
	Refresh_token string             `json:"refresh_token"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func (user *User) HashPassword() error {
	var password = user.Password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}

func (user *User) CheckPassword(incomingPassword string) error {
	var password = user.Password
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(incomingPassword))
	return err
}

func (user *User) GenerateToken() (string, error) {
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      user.Email,
		"first_name": user.First_name,
		"last_name":  user.Last_name,
		"exp":        time.Now().Local().Add(time.Hour * 24).Unix(),
	})
	godotenv.Load(".env")
	key:=os.Getenv("TOKEN_KEY")
	token, err := claim.SignedString([]byte(key))
	if err!=nil{
		return "", err
	}
	user.Token = token
	return token, nil
}


func (user *User) GenerateRefreshToken() (string, error) {
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      user.Email,
		"exp":        time.Now().Local().Add(time.Hour * 120).Unix(),
	})
	godotenv.Load(".env")
	key:=os.Getenv("REFRESH_TOKEN_KEY")
	token, err := claim.SignedString([]byte(key))
	if err!=nil{
		return "", err
	}
	user.Refresh_token = token
	return token, nil
}