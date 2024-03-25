package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id"`
	Field string             `json:"field" validate:"required"`
}
