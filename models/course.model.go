package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Owner       primitive.ObjectID `json:"owner" bson:"owner"`
	Title       string             `json:"title" validate:"required,min=2,max=50"`
	Description string             `json:"description" validate:"required,min=30,max=200"`
	Tag         primitive.ObjectID `json:"tag" bson:"tag"`
	Price       int                `json:"price" validate:"required"`
	Students    int                `json:"students"`
	Ratings     int                `json:"ratings"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
