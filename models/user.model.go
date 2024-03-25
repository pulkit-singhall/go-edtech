package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	First_name    string             `json:"first_name" validate:"required,min=2,max=20"`
	Last_name     string             `json:"last_name" validate:"required,min=2,max=20"`
	Email         string             `json:"email" validate:"email,required"`
	Password      string             `json:"password" validate:"required,min=8"`
	Refresh_token string             `json:"refresh_token"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"` 
}
