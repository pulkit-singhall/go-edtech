package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	ID       primitive.ObjectID `json:"_id"`
	Content  string             `json:"content"`
	OwnerID  primitive.ObjectID `json:"ownerID"`
	CourseID primitive.ObjectID `json:"courseID"`
}
