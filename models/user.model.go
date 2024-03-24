package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	First_name    string             `json:"first_name"`
	Last_name     string             `json:"last_name"`
	Email         string             `json:"email"`
	Password      string             `json:"password"`
	Refresh_token string             `json:"refresh_token"`
}
