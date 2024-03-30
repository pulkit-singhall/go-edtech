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

var categoryCollection = db.GetCollection("categories")

func CreateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var cat *models.Category
		bindErr := c.BindJSON(&cat)
		if bindErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": bindErr.Error()})
			return
		}
		valErr := Validator.Struct(cat)
		if valErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(), "detail": valErr.Error()})
			return
		}
		cat.ID = primitive.NewObjectID()
		res, insErr := categoryCollection.InsertOne(context.Background(), cat)
		if insErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": insErr.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "category created successfully", "id": res.InsertedID})
	}
}

func DeleteCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		categoryId := c.Param("categoryID")
		if categoryId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "categoryID is requred"})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(categoryId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		_, delErr := categoryCollection.DeleteOne(context.Background(), bson.M{"_id": cId})
		if delErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": delErr.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Category Deleted successfully"})
	}
}

func UpdateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		type NewCat struct {
			NewField string `json:"newField" bson:"newField"`
		}
		var cat *NewCat
		bindErr := c.BindJSON(&cat)
		if bindErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": bindErr.Error()})
			return
		}
		if cat.NewField == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(), "detail": "new field of category is required"})
			return
		}
		categoryId := c.Param("categoryID")
		if categoryId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error(), "detail": "category ID is required"})
			return
		}
		cId, hexErr := primitive.ObjectIDFromHex(categoryId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		_, updErr := categoryCollection.UpdateOne(context.Background(), bson.M{"_id": cId}, bson.M{"$set": bson.M{"field": cat.NewField}})
		if updErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": updErr.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "category updated successfully", "success": "true"})
	}
}
