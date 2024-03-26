package models

import (
	"time"

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