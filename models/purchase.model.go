package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Purchase struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	OwnerID  primitive.ObjectID `json:"onwerID"`
	CourseID primitive.ObjectID `json:"courseID"`
}
