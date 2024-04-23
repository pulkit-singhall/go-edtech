package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/controllers"
	"github.com/pulkit-singhall/go-edtech/middlewares"
)

func UserRoutes(incomingRoute *gin.Engine) {
	user := incomingRoute.Group("/api/v1/users")
	{
		user.POST("/signup", controllers.Signup())
		user.POST("/login", controllers.Login())
		user.GET("/getUser/:userID", controllers.GetUser())
		user.POST("/logout", middlewares.AuthCheck(), controllers.Logout())
		user.PATCH("/changePassword", middlewares.AuthCheck(), controllers.ChangePassword())
		user.GET("/getAllUsers", controllers.GetAllUsers())
		user.POST("/uploadAvatar", middlewares.AuthCheck(), controllers.UploadUserAvatar())
		user.POST("/changeAvatar", middlewares.AuthCheck(), controllers.ChangeUserAvatar())
		user.PATCH("/verifyEmail", middlewares.AuthCheck(), controllers.VerifyEmail())
	}
}
