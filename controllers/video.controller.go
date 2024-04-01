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

var videoCollection = db.GetCollection("videos")

func CreateVideo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DeleteVideo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetCourseVideos() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUserVideos() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetVideo() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		videoId := c.Param("videoID")
		email := c.Keys["email"]
		var user *models.User
		decErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(412, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		if videoId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "video ID is required"})
			return
		}
		vId, hexErr := primitive.ObjectIDFromHex(videoId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var video *models.Video
		decVErr := videoCollection.FindOne(context.Background(), bson.M{"_id": vId}).Decode(&video)
		if decVErr != nil {
			c.AbortWithStatusJSON(402, gin.H{"error": utils.InternalServerError.Error(), "detail": decVErr.Error()})
			return
		}
		if video.OwnerID != user.ID {
			c.AbortWithStatusJSON(412, gin.H{"error": utils.AuthorizeError.Error(), "detail": "user cant see this video"})
			return
		}
		c.JSON(200, gin.H{"video": video})
	}
}