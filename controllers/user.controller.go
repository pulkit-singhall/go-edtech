package controllers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pulkit-singhall/go-edtech/db"
	"github.com/pulkit-singhall/go-edtech/models"
	"github.com/pulkit-singhall/go-edtech/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Validator = validator.New()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user *models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(), "detail": err.Error()})
			return
		}
		validateErr := Validator.Struct(user)
		if validateErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.ValidationFailed.Error(), "detail": validateErr.Error()})
			return
		}
		userCollection := db.GetCollection("users")
		count, existErr := userCollection.CountDocuments(context.Background(), bson.M{"email": user.Email})
		if existErr != nil {
			c.AbortWithStatusJSON(402, gin.H{"error": utils.InternalServerError.Error(), "detail": existErr.Error()})
			return
		}
		if count > 0 {
			c.AbortWithStatusJSON(410, gin.H{"error": utils.UserAlreadyExists.Error()})
			return
		}
		hashErr := user.HashPassword()
		if hashErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": hashErr.Error()})
			return
		}
		user.CreatedAt = time.Now().Local()
		user.UpdatedAt = time.Now().Local()
		user.ID = primitive.NewObjectID()
		id, createErr := userCollection.InsertOne(context.Background(), user)
		if createErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": createErr.Error()})
			return
		}
		c.JSON(200, gin.H{"success": "true", "id": id.InsertedID})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userId := c.Param("userID")
		if userId == "" {
			c.JSON(400, gin.H{"error": utils.QueryParamMissing, "detail": "User ID is required"})
			return
		}
		userCollection := db.GetCollection("users")
		id, hexErr := primitive.ObjectIDFromHex(userId)
		if hexErr != nil {
			c.JSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var user *models.User
		userErr := userCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
		if userErr != nil {
			c.JSON(500, gin.H{"error": utils.UserNotFound, "detail": userErr.Error()})
			return
		}
		c.JSON(200, gin.H{
			"user": user,
		})
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userCollection := db.GetCollection("users")
		cur, err := userCollection.Find(context.Background(), bson.M{})
		if err != nil {
			c.JSON(500, gin.H{
				"error":  utils.InternalServerError,
				"detail": err.Error(),
			})
		}
		var allUser []*models.User
		for cur.Next(context.Background()) {
			var user *models.User
			decodeErr := cur.Decode(&user)
			if decodeErr != nil {
				c.JSON(500, gin.H{"error": utils.InternalServerError, "detail": decodeErr.Error()})
				return
			}
			allUser = append(allUser, user)
		}
		c.JSON(200, gin.H{
			"users": allUser,
		})
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
