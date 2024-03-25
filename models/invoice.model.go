package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	UserID    primitive.ObjectID `json:"userID" bson:"userID"`
	CourseID  primitive.ObjectID `json:"courseID" bson:"courseID"`
	Amount    int                `json:"amount"`
	Type      string             `json:"type" validate:"required,eq=CARD|eq=UPI|eq=EMI"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
