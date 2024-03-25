package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Content   string             `json:"content" validate:"required,min=20,max=300"`
	OwnerID   primitive.ObjectID `json:"ownerID" bson:"ownerID"`
	CourseID  primitive.ObjectID `json:"courseID" bson:"courseID"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
