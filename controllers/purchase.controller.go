package controllers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/db"
	"github.com/pulkit-singhall/go-edtech/middlewares"
	"github.com/pulkit-singhall/go-edtech/models"
	"github.com/pulkit-singhall/go-edtech/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var purchaseCollection = db.GetCollection("purchases")

func PurchaseCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email := c.Keys["email"]
		courseId := c.Param("courseID")
		type paymentType struct {
			PaymentType string `json:"paymentType" bson:"paymentType" validate:"required,eq=CARD|eq=UPI|eq=EMI"`
		}
		var payment *paymentType
		bindErr := c.BindJSON(&payment)
		if bindErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": bindErr.Error()})
			return
		}
		valErr := Validator.Struct(payment)
		if valErr != nil {
			c.AbortWithStatusJSON(412, gin.H{"error": utils.QueryBodyMissing.Error(), "detail": valErr.Error()})
			return
		}
		if courseId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "courseID is required"})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(courseId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var course *models.Course
		decCErr := courseCollection.FindOne(context.Background(), bson.M{"_id": cId}).Decode(&course)
		if decCErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": decCErr.Error()})
			return
		}
		var user *models.User
		decUErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decUErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UserNotFound.Error(), "detail": decUErr.Error()})
			return
		}
		if user.ID == course.Owner {
			c.AbortWithStatusJSON(412, gin.H{"error": utils.AuthorizeError.Error(), "detail": "user already created this course"})
			return
		}
		invoice,invErr:=middlewares.GenerateInvoice(course.Price, payment.PaymentType, user.ID, cId)
		if invErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": invErr.Error()})
			return
		}
		var purchase models.Purchase
		purchase.ID = primitive.NewObjectID()
		purchase.CourseID = cId
		purchase.UserID = user.ID
		purchase.CreatedAt = time.Now().Local()
		purchase.InvoiceID = invoice.ID
		_, updErr := courseCollection.UpdateOne(context.Background(), bson.M{"_id": cId}, bson.M{"$set": bson.M{"students": course.Students + 1}})
		if updErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": updErr.Error()})
			return
		}
		pId, insErr := purchaseCollection.InsertOne(context.Background(), purchase)
		if insErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": insErr.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "course purchased successfully", "id": pId.InsertedID, "invoice": invoice})
	}
}

// func CancelPurchase() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		_,cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()
// 		email:=c.Keys["email"]

// 	}
// }

func GetUserPurchasedCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}
