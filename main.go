package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pulkit-singhall/go-edtech/db"
	"github.com/pulkit-singhall/go-edtech/routes"
)

func main() {
	// env file variables
	godotenv.Load(".env")
	PORT := os.Getenv("PORT")

	// database connect
	client, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())

	if PORT == "" {
		PORT = "8000"
	}

	// gin framework
	router := gin.New()

	routes.UserRoutes(router)
	routes.CourseRoutes(router)

	erro := router.Run(":" + PORT)
	if erro != nil {
		panic(erro)
	}
}
