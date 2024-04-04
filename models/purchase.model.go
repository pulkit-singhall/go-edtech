package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Purchase struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	UserID    primitive.ObjectID `json:"onwerID" bson:"ownerID"`
	CourseID  primitive.ObjectID `json:"courseID" bson:"courseID"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	InvoiceID primitive.ObjectID `json:"invoiceID" bson:"invoiceID"`
}
