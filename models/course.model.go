package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	ID          primitive.ObjectID `json:"_id"`
	Owner       primitive.ObjectID `json:"owner"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Tag         primitive.ObjectID `json:"tag"`
	Price       int                `json:"price"`
	Students    int                `json:"students"`
	Ratings     int                `json:"ratings"`
}
