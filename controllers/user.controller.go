package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/utils"
)

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		if userID == "uti" {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.QueryParamMissing.Error()})
		}

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
