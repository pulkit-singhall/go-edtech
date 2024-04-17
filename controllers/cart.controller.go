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

var cartCollection = db.GetCollection("carts")

func AddCourseToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email := c.Keys["email"]
		courseId := c.Param("courseID")
		var user *models.User
		decErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(412, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(courseId)
		if hexErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var course *models.Course
		decCErr := courseCollection.FindOne(context.Background(), bson.M{"_id": cId}).Decode(&course)
		if decCErr != nil {
			c.AbortWithStatusJSON(412, gin.H{"error": utils.InternalServerError.Error(), "detail": decCErr.Error()})
			return
		}
		if course.Owner == user.ID {
			c.AbortWithStatusJSON(500, gin.H{"error": "the course owner is the current user only"})
			return
		}
		var cart models.Cart
		cart.ID = primitive.NewObjectID()
		cart.OwnerID = user.ID
		cart.CourseID = cId
		newC, insErr := cartCollection.InsertOne(context.Background(), cart)
		if insErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": insErr.Error()})
			return
		}
		c.JSON(201, gin.H{"new course added to the cart": newC.InsertedID})
	}
}

func RemoveCourseFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email := c.Keys["email"]
		courseId := c.Param("courseID")
		var user *models.User
		decErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(412, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(courseId)
		if hexErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		_, delErr := cartCollection.DeleteOne(context.Background(), bson.M{"ownerID": user.ID, "courseID": cId})
		if delErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": delErr.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "course removed from cart"})
	}
}

func GetUserCartCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email := c.Keys["email"]
		var user *models.User
		decErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr != nil {
			c.AbortWithStatusJSON(412, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		pipeline := []bson.M{}
		match := bson.M{
			"$match": bson.M{
				"ownerID": user.ID,
			},
		}
		pipeline = append(pipeline, match)
		cur, pipeErr := cartCollection.Aggregate(context.Background(), pipeline)
		if pipeErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.PipelineError.Error(), "detail": pipeErr.Error()})
			return
		}
		cartCourses := []bson.M{}
		curErr := cur.All(context.Background(), &cartCourses)
		if curErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": curErr.Error()})
			return
		}
		c.JSON(200, gin.H{"cart courses": cartCourses})
	}
}
