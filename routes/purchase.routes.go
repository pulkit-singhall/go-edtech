package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/controllers"
	"github.com/pulkit-singhall/go-edtech/middlewares"
)

func PurchaseRoutes(incomingRoute *gin.Engine) {
	purchases := incomingRoute.Group("/api/v1/purchases")
	{
		purchases.POST("/purchase/:courseID", middlewares.AuthCheck(), controllers.PurchaseCourse())
		purchases.GET("/userPurchases", middlewares.AuthCheck(), controllers.GetUserPurchasedCourses())
	}
}
