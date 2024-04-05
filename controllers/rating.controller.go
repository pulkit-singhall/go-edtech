package controllers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/db"
	"github.com/pulkit-singhall/go-edtech/models"
	"github.com/pulkit-singhall/go-edtech/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ratingCollection = db.GetCollection("ratings")

func CreateRating() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email := c.Keys["email"]
		courseId := c.Param("courseID")
		if courseId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "courseID is required"})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(courseId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var rating *models.Rating
		bindErr := c.BindJSON(&rating)
		if bindErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": bindErr.Error()})
			return
		}
		valErr := Validator.Struct(rating)
		if valErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(), "detail": valErr.Error()})
			return
		}
		var user *models.User
		decErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		rating.ID = primitive.NewObjectID()
		rating.CreatedAt = time.Now().Local()
		rating.UpdatedAt = time.Now().Local()
		rating.OwnerID = user.ID
		rating.CourseID = cId
		var course *models.Course
		decCErr := courseCollection.FindOne(context.Background(), bson.M{"_id": cId}).Decode(&course)
		if decCErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": decCErr.Error()})
			return
		}
		courseCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": cId}, bson.M{"$set": bson.M{"ratings": course.Ratings + 1}})
		rId, insErr := ratingCollection.InsertOne(context.Background(), rating)
		if insErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": insErr.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "rating added successfully", "id": rId.InsertedID})
	}
}

func DeleteRating() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		ratingId := c.Param("ratingID")
		email := c.Keys["email"]
		if ratingId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "ratingID is required"})
			return
		}
		rId, hexErr := primitive.ObjectIDFromHex(ratingId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var user *models.User
		decErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		var rating *models.Rating
		decRErr := ratingCollection.FindOne(context.Background(), bson.M{"_id": rId}).Decode(&rating)
		if decRErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": decRErr.Error()})
			return
		}
		if user.ID != rating.OwnerID {
			c.AbortWithStatusJSON(410, gin.H{"error": utils.AuthorizeError.Error()})
			return
		}
		var course *models.Course
		decCErr := courseCollection.FindOne(context.Background(), bson.M{"_id": rating.CourseID}).Decode(&course)
		if decCErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": decCErr.Error()})
			return
		}
		courseCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": rating.CourseID}, bson.M{"$set": bson.M{"ratings": course.Ratings - 1}})
		_, delErr := ratingCollection.DeleteOne(context.Background(), bson.M{"_id": rId})
		if delErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": delErr.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "rating deleted successfully"})
	}
}

func UpdateRating() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		ratingId := c.Param("ratingID")
		email := c.Keys["email"]
		type Rate struct {
			NewRate int `json:"newRate" bson:"newRate"`
		}
		var rate *Rate
		c.BindJSON(&rate)
		if ratingId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "ratingID is required"})
			return
		}
		rId, hexErr := primitive.ObjectIDFromHex(ratingId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var user *models.User
		decErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		var rating *models.Rating
		decRErr := ratingCollection.FindOne(context.Background(), bson.M{"_id": rId}).Decode(&rating)
		if decRErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": decRErr.Error()})
			return
		}
		if user.ID != rating.OwnerID {
			c.AbortWithStatusJSON(410, gin.H{"error": utils.AuthorizeError.Error()})
			return
		}
		_, updErr := ratingCollection.UpdateOne(context.Background(), bson.M{"_id": rId}, bson.M{"$set": bson.M{"rate": rate.NewRate}})
		if updErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": updErr.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "rating updated successfully"})
	}
}

func GetCourseRatings() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUserRatings() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetRating() gin.HandlerFunc{
	return func(c *gin.Context) {
		_,cancel:=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		email:=c.Keys["email"]
		ratingId:=c.Param("ratingID")
		if ratingId == ""{
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "rating ID is required"})
			return
		}
		var user *models.User
		decErr:=userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr!=nil{
			c.AbortWithStatusJSON(412, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}	
		rId,hexErr:=primitive.ObjectIDFromHex(ratingId)
		if hexErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var rating *models.Rating
		decRErr:=ratingCollection.FindOne(context.Background(), bson.M{"_id": rId}).Decode(&rating)
		if decRErr!=nil{
			c.AbortWithStatusJSON(412, gin.H{"error": utils.InternalServerError.Error(), "detail": decRErr.Error()})
			return
		}
		if rating.OwnerID!=user.ID{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.AuthorizeError.Error(), "detail": "user not authorize to see this rating"})
			return
		}
		c.JSON(200, gin.H{"rating": rating})
	}
}