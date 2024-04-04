package middlewares

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pulkit-singhall/go-edtech/db"
	"github.com/pulkit-singhall/go-edtech/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var invoiceCollection = db.GetCollection("invoices")

var Validator = validator.New()

func GenerateInvoice(price int, paymentType string, userId primitive.ObjectID, courseId primitive.ObjectID) (*models.Invoice, error) {
	var invoice models.Invoice
	invoice.ID = primitive.NewObjectID()
	invoice.Amount = price
	invoice.PaymentType = paymentType
	invoice.CreatedAt = time.Now().Local()
	invoice.CourseID = courseId
	invoice.UserID = userId
	valErr := Validator.Struct(invoice)
	if valErr != nil {
		return nil, valErr
	}
	_, insErr := invoiceCollection.InsertOne(context.Background(), invoice)
	if insErr != nil {
		return nil, insErr
	}
	return &invoice, nil
}
