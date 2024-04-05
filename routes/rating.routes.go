package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/controllers"
	"github.com/pulkit-singhall/go-edtech/middlewares"
)

func RatingRoutes(incomingRoute *gin.Engine) {
	rating := incomingRoute.Group("/api/v1/ratings")
	{
		rating.POST("/create/:courseID", middlewares.AuthCheck(), controllers.CreateRating())
		rating.DELETE("/delete/:ratingID", middlewares.AuthCheck(), controllers.DeleteRating())
		rating.GET("/getCourseRatings/:courseID", controllers.GetCourseRatings())
		rating.GET("/getUserRatings", middlewares.AuthCheck(), controllers.GetUserRatings())
		rating.PATCH("/update/:ratingID", middlewares.AuthCheck(), controllers.UpdateRating())
		rating.GET("/rating/:ratingID", middlewares.AuthCheck(), controllers.GetRating())
	}
}
