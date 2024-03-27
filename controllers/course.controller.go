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

	}
}

func GetCoursesByOwnerID() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetCourseByID() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
