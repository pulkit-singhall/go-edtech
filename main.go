package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	PORT:=os.Getenv("PORT")
	if PORT==""{
		PORT="8000"
	}
	router := gin.New()

	err:=router.Run(":"+PORT)
	if err!=nil{
		panic(err)
	}
}
