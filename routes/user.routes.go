package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/controllers"
)

func UserRoutes(incomingRoute *gin.Engine){
	user:=incomingRoute.Group("/api/v1/users")
	{
		user.POST("/signup", controllers.Signup())
		user.POST("/login", controllers.Login())
		user.GET("/getUser/:userID", controllers.GetUser())
		user.POST("/logout", controllers.Logout())
		user.PATCH("/changePassword", controllers.ChangePassword())
	}
}