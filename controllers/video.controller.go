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

var videoCollection = db.GetCollection("videos")

func CreateVideo() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		videoFile, _, vErr := c.Request.FormFile("videoFile")
		thumbnail, _, tErr := c.Request.FormFile("thumbnail")
		title := c.Request.FormValue("title")
		email := c.Keys["email"]
		courseId := c.Param("courseID")
		if courseId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "course ID is required"})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(courseId)
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
		if title == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(), "detail": "video title is required"})
			return
		}
		if tErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(), "detail": tErr.Error()})
			return
		}
		if vErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(), "detail": vErr.Error()})
			return
		}
		vUrl, vPId, uplVErr := middlewares.UploadFile(videoFile)
		tUrl, tPId, uplTErr := middlewares.UploadFile(thumbnail)
		if uplVErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UploadFileError.Error(), "detail": uplVErr.Error()})
			return
		}
		if uplTErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UploadFileError.Error(), "detail": uplTErr.Error()})
			return
		}
		var video models.Video
		video.ID = primitive.NewObjectID()
		video.CourseID = cId
		video.CreatedAt = time.Now().Local()
		video.ThumbnailUrl = tUrl
		video.ThumbnailID = tPId
		video.VideoFileUrl = vUrl
		video.VideoFileID = vPId
		video.Title = title
		video.OwnerID = user.ID
		vId, insErr := videoCollection.InsertOne(context.Background(), video)
		if insErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": insErr.Error()})
			return
		}
		c.JSON(200, gin.H{"id": vId.InsertedID, "message": "video created successfully"})
	}
}

func DeleteVideo() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		videoId:=c.Param("videoID")
		email:=c.Keys["email"]
		if videoId==""{
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "video ID is required"})
			return
		}
		vId,hexErr:=primitive.ObjectIDFromHex(videoId)
		if hexErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var user *models.User
		decErr:=userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		var video *models.Video
		decVErr:=videoCollection.FindOne(context.Background(), bson.M{"_id": vId}).Decode(&video)
		if decVErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": decVErr.Error()})
			return
		}
		if user.ID != video.OwnerID{
			c.AbortWithStatusJSON(412, gin.H{"error": utils.AuthorizeError.Error(), "detail": "user not authorize to delete this video"})
			return
		}
		tPId := video.ThumbnailID
		vPId := video.VideoFileID
		dTErr:=middlewares.DeleteImageFile(tPId)
		dVErr:=middlewares.DeleteVideoFile(vPId)
		if dTErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.DeleteFileError.Error(), "detail": dTErr.Error()})
			return
		}
		if dVErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.DeleteFileError.Error(), "detail": dVErr.Error()})
			return
		}
		_,delErr:=videoCollection.DeleteOne(context.Background(), bson.M{"_id": vId})
		if delErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": delErr.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "video deleted successfully", "success": "true"})
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
