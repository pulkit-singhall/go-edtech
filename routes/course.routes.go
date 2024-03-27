package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/controllers"
	"github.com/pulkit-singhall/go-edtech/middlewares"
)

func CourseRoutes(incomingRoute *gin.Engine) {
	course := incomingRoute.Group("/api/v1/courses")
	{
		course.POST("/create", middlewares.AuthCheck(), controllers.CreateCourse())
		course.DELETE("/delete/:courseID", middlewares.AuthCheck(), controllers.DeleteCourse())
		course.GET("/getCourses/:ownerID", controllers.GetCoursesByOwnerID())
		course.GET("/course/:courseID", controllers.GetCourseByID())
	}
}
