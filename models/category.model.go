package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID    primitive.ObjectID `json:"_id"`
	Field string             `json:"field"`
}
