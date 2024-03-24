package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Rating struct {
	ID       primitive.ObjectID `json:"_id"`
	Rate     int                `json:"rate"`
	OwnerID  primitive.ObjectID `json:"ownerID"`
	CourseID primitive.ObjectID `json:"courseID"`
}
