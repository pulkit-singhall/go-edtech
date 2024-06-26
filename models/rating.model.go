package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rating struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Rate      int                `json:"rate" validate:"required"`
	OwnerID   primitive.ObjectID `json:"ownerID" bson:"ownerID"`
	CourseID  primitive.ObjectID `json:"courseID" bson:"courseID"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
