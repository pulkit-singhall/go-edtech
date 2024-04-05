package controllers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pulkit-singhall/go-edtech/db"
	"github.com/pulkit-singhall/go-edtech/middlewares"
	"github.com/pulkit-singhall/go-edtech/models"
	"github.com/pulkit-singhall/go-edtech/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Validator = validator.New()

var userCollection = db.GetCollection("users")

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
		_, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		var user *models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(), "details": err.Error()})
			return
		}
		if user.Email == "" || user.Password == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(),
				"details": "email or password is missing"})
			return
		}
		var existing *models.User
		existErr := userCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existing)
		if existErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UserNotFound.Error(), "detail": existErr.Error()})
			return
		}
		passErr := existing.CheckPassword(user.Password)
		if passErr != nil {
			c.AbortWithStatusJSON(402, gin.H{"error": utils.PasswordWrong.Error(), "detail": passErr.Error()})
			return
		}
		token, refresh, tokenErr := utils.GenerateNewTokens(user.Email)
		if tokenErr != nil {
			c.AbortWithStatusJSON(415, gin.H{"error": utils.TokenError.Error(), "detail": tokenErr.Error()})
			return
		}
		c.SetCookie("token", token, int(time.Hour*48), "/", "localhost", true, true)
		c.SetCookie("refresh_token", refresh, int(time.Hour*240), "/", "localhost", true, true)
		c.JSON(200, gin.H{
			"user":          existing,
			"token":         token,
			"refresh_token": refresh,
		})
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userId := c.Param("userID")
		if userId == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing, "detail": "User ID is required"})
			return
		}
		id, hexErr := primitive.ObjectIDFromHex(userId)
		if hexErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.HexIdError.Error(), "detail": hexErr.Error()})
			return
		}
		var user *models.User
		userErr := userCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
		if userErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UserNotFound, "detail": userErr.Error()})
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
		cur, err := userCollection.Find(context.Background(), bson.M{})
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error":  utils.InternalServerError,
				"detail": err.Error(),
			})
			return
		}
		var allUser []*models.User
		for cur.Next(context.Background()) {
			var user *models.User
			decodeErr := cur.Decode(&user)
			if decodeErr != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError, "detail": decodeErr.Error()})
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
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email := c.Keys["email"]
		_, err := userCollection.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"token": "", "refresh_token": ""}})
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": err.Error()})
			return
		}
		c.SetCookie("token", "", -1, "/", "localhost", true, true)
		c.SetCookie("refresh_token", "", -1, "/", "localhost", true, true)
		c.JSON(200, gin.H{"message": "logout success", "success": "true"})
	}
}

func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email := c.Keys["email"]
		type Password struct {
			Old string `json:"oldPassword" bson:"oldPassword"`
			New string `json:"newPassword" bson:"newPassword"`
		}
		var pass *Password
		err := c.BindJSON(&pass)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": err.Error()})
			return
		}
		if pass.Old == "" || pass.New == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryBodyMissing.Error(), "detail": "passwords are missing"})
			return
		}
		var user *models.User
		findErr := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if findErr != nil {
			c.AbortWithStatusJSON(402, gin.H{"error": utils.UserNotFound.Error(), "detail": findErr.Error()})
			return
		}
		checkErr := user.CheckPassword(pass.Old)
		if checkErr != nil {
			c.AbortWithStatusJSON(405, gin.H{"error": utils.PasswordWrong.Error(), "detail": checkErr.Error()})
			return
		}
		user.Password = pass.New
		hashErr := user.HashPassword()
		if hashErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": hashErr.Error()})
			return
		}
		_, updErr := userCollection.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"password": user.Password}})
		if updErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": updErr.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "password updated successfully", "success": "true"})
	}
}

func UploadUserAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email := c.Keys["email"]
		avatar, _, avatarErr := c.Request.FormFile("avatar")
		if avatarErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.FileError.Error(), "detail": avatarErr.Error()})
			return
		}
		url, pId, uplErr := middlewares.UploadFile(avatar)
		if uplErr != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UploadFileError.Error(), "detail": uplErr.Error()})
			return
		}
		_,updErr:=userCollection.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"avatar": url, "avatarId": pId}})
		if updErr!=nil{
			middlewares.DeleteImageFile(pId)
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": updErr.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "avatar uploaded", "success": "true"})
	}
}

func ChangeUserAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		email:=c.Keys["email"]
		avatar,_,avtErr:=c.Request.FormFile("avatar")
		if avtErr!=nil{
			c.AbortWithStatusJSON(400, gin.H{"error": utils.FileError.Error(), "detail": avtErr.Error()})
			return 
		}
		var user *models.User
		decErr:=userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if decErr!=nil{
			c.AbortWithStatusJSON(412, gin.H{"error": utils.UserNotFound.Error(), "detail": decErr.Error()})
			return
		}
		oldPId:=user.AvatarId
		url,pId,uplErr:=middlewares.UploadFile(avatar)
		if uplErr!=nil{
			c.AbortWithStatusJSON(500, gin.H{"error": utils.UploadFileError.Error(), "detail": uplErr.Error()})
			return
		}
		_,updErr:=userCollection.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"avatar": url, "avatarId": pId}})
		if updErr!=nil{
			middlewares.DeleteImageFile(pId)
			c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": updErr.Error()})
			return
		}
		middlewares.DeleteImageFile(oldPId)
		c.JSON(200, gin.H{"message": "avatar updated successfully", "success": "true"})
	}
}
