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

func CreateCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		email := c.Keys["email"]
		var course *models.Course
		err := c.BindJSON(&course)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": err.Error()})
			return
		}
		valErr := Validator.Struct(course)
		if valErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.ValidationFailed.Error(), "detail": valErr.Error()})
			return
		}
		userCollection := db.GetCollection("users")
		var user *models.User
		decErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		course.ID = primitive.NewObjectID()
		course.Owner = user.ID
		course.Students = 0
		course.Ratings = 0
		course.CreatedAt = time.Now().Local()
		course.UpdatedAt = time.Now().Local()
		courseCollection := db.GetCollection("courses")
		id, createErr := courseCollection.InsertOne(context.Background(), course)
		if createErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": createErr.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "course created successfully", "id": id.InsertedID})
	}
}

func DeleteCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		courseId := c.Param("courseID")
		if courseId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "courseID is required"})
			return
		}
		email := c.Keys["email"]
		userCollection := db.GetCollection("users")
		courseCollection := db.GetCollection("courses")
		var user *models.User
		decErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(courseId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var course *models.Course
		cDecErr := courseCollection.FindOne(context.Background(), bson.M{"_id": cId}).Decode(&course)
		if cDecErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": cDecErr.Error()})
			return
		}
		if course.Owner != user.ID {
			c.AbortWithStatusJSON(412, gin.H{"error": utils.AuthorizeError.Error(), "detail": "Not a course of current user"})
			return
		}
		res, delErr := courseCollection.DeleteOne(context.Background(), bson.M{"_id": cId})
		if delErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": delErr.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "course deleted successfully", "count": res.DeletedCount})
	}
}

func GetCoursesByOwnerID() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		ownerId := c.Param("ownerID")
		if ownerId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "owner ID is required"})
			return
		}
		_, hexErr := primitive.ObjectIDFromHex(ownerId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
	}
}

func GetCourseByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
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
		courseCollection := db.GetCollection("courses")
		var course *models.Course
		cDecErr := courseCollection.FindOne(context.Background(), bson.M{"_id": cId}).Decode(&course)
		if cDecErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": cDecErr.Error()})
			return
		}
		c.JSON(200, gin.H{"course": course})
	}
}
