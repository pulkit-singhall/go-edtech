package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/controllers"
	"github.com/pulkit-singhall/go-edtech/middlewares"
)

func VideoRoutes(incomingRoute *gin.Engine) {
	video := incomingRoute.Group("/api/v1/videos")
	{
		video.POST("/create/:courseID", middlewares.AuthCheck(), controllers.CreateVideo())
		video.DELETE("/delete/:videoID", middlewares.AuthCheck(), controllers.DeleteVideo())
		video.GET("/courseVideos/:courseID", middlewares.AuthCheck(), controllers.GetCourseVideos())
		video.GET("/video/:videoID", middlewares.AuthCheck(), controllers.GetVideo())
		video.GET("/userVideos", middlewares.AuthCheck(), controllers.GetUserVideos())
	}
}
