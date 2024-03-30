package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/controllers"
)

func CategoryRoutes(incomingRoute *gin.Engine){
	category:=incomingRoute.Group("/api/v1/category")
	{
		category.POST("/create", controllers.CreateCategory())
		category.PATCH("/update/:categoryID", controllers.UpdateCategory())
		category.DELETE("/delete/:categoryID", controllers.DeleteCategory())
	}
}