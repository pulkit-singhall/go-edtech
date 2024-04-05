package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pulkit-singhall/go-edtech/db"
	"github.com/pulkit-singhall/go-edtech/routes"
)

func main() {
	// env file variables
	envErr := godotenv.Load(".env")
	if envErr != nil {
		panic(envErr)
	}
	PORT := os.Getenv("PORT")

	// database connect
	_, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	if PORT == "" {
		PORT = "8000"
	}

	// gin framework
	router := gin.New()

	routes.UserRoutes(router)
	routes.CourseRoutes(router)
	routes.RatingRoutes(router)
	routes.CategoryRoutes(router)
	routes.CommentRoutes(router)
	routes.VideoRoutes(router)
	routes.PurchaseRoutes(router)
	routes.CartRoutes(router)

	erro := router.Run(":" + PORT)
	if erro != nil {
		panic(erro)
	}
}
