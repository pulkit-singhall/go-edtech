package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	CourseID  primitive.ObjectID `json:"courseID" bson:"courseID"`
	Title     string             `json:"title" validate:"required,min=2,max=200"`
	VideoFile string             `json:"videoFile" validate:"required"`
	Thumbnail string             `json:"thumbnail" validate:"required"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
