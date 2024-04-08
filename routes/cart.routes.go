package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/controllers"
	"github.com/pulkit-singhall/go-edtech/middlewares"
)

func CartRoutes(incomingRoute *gin.Engine){
	cart:=incomingRoute.Group("/api/v1/cart")
	{
		cart.POST("/addCourse/:courseID", middlewares.AuthCheck(), controllers.AddCourseToCart())
		cart.DELETE("/removeCourse/:courseID", middlewares.AuthCheck(), controllers.RemoveCourseFromCart())
		cart.GET("/userCart", middlewares.AuthCheck(), controllers.GetUserCartCourses())
	}
}