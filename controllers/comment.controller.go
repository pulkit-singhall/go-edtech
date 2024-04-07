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

var commentCollection = db.GetCollection("comments")

func CreateComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		courseId := c.Param("courseID")
		email := c.Keys["email"]
		if courseId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error()})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(courseId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var comment *models.Comment
		bindErr := c.BindJSON(&comment)
		if bindErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": bindErr.Error()})
			return
		}
		valErr := Validator.Struct(comment)
		if valErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(), "detail": valErr.Error()})
			return
		}
		var user *models.User
		decErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(402, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		comment.ID = primitive.NewObjectID()
		comment.CourseID = cId
		comment.CreatedAt = time.Now().Local()
		comment.UpdatedAt = time.Now().Local()
		comment.OwnerID = user.ID
		id, insErr := commentCollection.InsertOne(context.Background(), comment)
		if insErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": insErr.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "comment created successfully", "id": id.InsertedID})
	}
}

func UpdateComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email := c.Keys["email"]
		commentId := c.Param("commentID")
		if commentId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "comment ID is required"})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(commentId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		type Com struct {
			NewContent string `json:"newContent" bson:"newContent"`
		}
		var com *Com
		bindErr := c.BindJSON(&com)
		if bindErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": bindErr.Error()})
			return
		}
		if com.NewContent == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(), "detail": "new content is required"})
			return
		}
		var comment *models.Comment
		decCErr := commentCollection.FindOne(context.Background(), bson.M{"_id": cId}).Decode(&comment)
		if decCErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": decCErr.Error()})
			return
		}
		var user *models.User
		decUErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decUErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": decUErr.Error()})
			return
		}
		if user.ID != comment.OwnerID {
			c.AbortWithStatusJSON(410, gin.H{"error": utils.AuthorizeError.Error(), "detail": "user not authorize to update this comment"})
			return
		}
		_, updErr := commentCollection.UpdateOne(context.Background(), bson.M{"_id": cId}, bson.M{"$set": bson.M{"content": com.NewContent}})
		if updErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": updErr.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "comment updated successfully", "success": "true"})
	}
}

func DeleteComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email := c.Keys["email"]
		commentId := c.Param("commentID")
		if commentId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "comment ID is required"})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(commentId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var comment *models.Comment
		decCErr := commentCollection.FindOne(context.Background(), bson.M{"_id": cId}).Decode(&comment)
		if decCErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": decCErr.Error()})
			return
		}
		var user *models.User
		decUErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decUErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": decUErr.Error()})
			return
		}
		if user.ID != comment.OwnerID {
			c.AbortWithStatusJSON(410, gin.H{"error": utils.AuthorizeError.Error(), "detail": "user not authorize to delete this comment"})
			return
		}
		_, delErr := commentCollection.DeleteOne(context.Background(), bson.M{"_id": cId})
		if delErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": delErr.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "comment deleted successfully", "success": "true"})
	}
}

func GetCourseComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		courseId:=c.Param("courseID")
		cId,hexErr:=primitive.ObjectIDFromHex(courseId)
		if hexErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		pipeline:=[]bson.M{}
		match:=bson.M{
			"$match": bson.M{
				"courseID": cId,
			},
		}
		pipeline = append(pipeline, match)
		cur,pipeErr:=commentCollection.Aggregate(context.Background(), pipeline)
		if pipeErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.PipelineError.Error(), "detail": pipeErr.Error()})
			return
		}
		var comments []bson.M
		curErr:=cur.All(context.Background(), &comments)
		if curErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": curErr.Error()})
			return
		}
		c.JSON(200, gin.H{"course comments": comments})
	}
}

func GetUserComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email:=c.Keys["email"]
		var user *models.User
		decErr:=userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr!=nil{
			c.AbortWithStatusJSON(412, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		pipeline:=[]bson.M{}
		match:=bson.M{
			"$match": bson.M{
				"ownerID": user.ID,
			},
		}
		pipeline = append(pipeline, match)
		cur,pipeErr:=commentCollection.Aggregate(context.Background(), pipeline)
		if pipeErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.PipelineError.Error(), "detail": pipeErr.Error()})
			return
		}
		var comments []bson.M
		curErr:=cur.All(context.Background(), &comments)
		if curErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": curErr.Error()})
			return
		}
		c.JSON(200, gin.H{"user comments": comments})
	}
}

func GetComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email := c.Keys["email"]
		commentId := c.Param("commentID")
		if commentId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "comment ID is required"})
			return
		}
		var user *models.User
		decErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(412, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(commentId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var comment *models.Comment
		decCErr := commentCollection.FindOne(context.Background(), bson.M{"_id": cId}).Decode(&comment)
		if decCErr != nil {
			c.AbortWithStatusJSON(412, gin.H{"error": utils.InternalServerError.Error(), "detail": decCErr.Error()})
			return
		}
		if comment.OwnerID != user.ID {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.AuthorizeError.Error(), "detail": "user not authorize to see this comment"})
			return
		}
		c.JSON(200, gin.H{"comment": comment})
	}
}
