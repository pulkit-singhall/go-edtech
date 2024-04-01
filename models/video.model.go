package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	CourseID    primitive.ObjectID `json:"courseID" bson:"courseID"`
	OwnerID     primitive.ObjectID `json:"ownerID" bson:"ownerID"`
	Title       string             `json:"title" validate:"required,min=2,max=200"`
	VideoFileID string             `json:"videoFileID"`
	ThumbnailID string             `json:"thumbnailID"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}
