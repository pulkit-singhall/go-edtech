package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cart struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	OwnerID  primitive.ObjectID `json:"ownerID"`
	CourseID primitive.ObjectID `json:"courseID"`
}
