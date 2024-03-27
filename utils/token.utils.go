package utils

import (
	"context"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/pulkit-singhall/go-edtech/db"
	"github.com/pulkit-singhall/go-edtech/models"
	"go.mongodb.org/mongo-driver/bson"
)

func VerifyToken(token string) (*jwt.Token, error) {
	godotenv.Load(".env")
	key := os.Getenv("TOKEN_KEY")
	tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	return tk, err
}

func VerifyRefreshToken(token string) (*jwt.Token, error) {
	godotenv.Load(".env")
	key := os.Getenv("REFRESH_TOKEN_KEY")
	tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	return tk, err
}

func GenerateNewTokens(email string) (string, string, error) { // token, refresh, err
	userCollection := db.GetCollection("users")
	var user *models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", "", err
	}
	token, tokErr := user.GenerateToken()
	if tokErr != nil {
		return "", "", tokErr
	}
	refresh, refErr := user.GenerateRefreshToken()
	if refErr != nil {
		return "", "", refErr
	}
	_,updErr:=userCollection.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"token": token, "refresh_token": refresh}})
	if updErr!=nil{
		return "","", updErr
	}
	return token, refresh, nil
}
