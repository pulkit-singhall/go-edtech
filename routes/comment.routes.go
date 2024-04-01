package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/controllers"
	"github.com/pulkit-singhall/go-edtech/middlewares"
)

func CommentRoutes(incomingRoute *gin.Engine) {
	comment := incomingRoute.Group("/api/v1/comments")
	{
		comment.POST("/create/:courseID", middlewares.AuthCheck(), controllers.CreateComment())
		comment.PATCH("/update/:commentID", middlewares.AuthCheck(), controllers.UpdateComment())
		comment.DELETE("/delete/:commentID", middlewares.AuthCheck(), controllers.DeleteComment())
		comment.GET("/courseComments/:courseID", controllers.GetCourseComments())
		comment.GET("/userComments/:userID", controllers.GetUserComments())
		comment.GET("/comment/:commentID", controllers.GetComment())
	}
}
