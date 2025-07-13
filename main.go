package main

import (
	"net/http"

	"github.com/MegeKaplan/megebase-identity-service/database"
	"github.com/MegeKaplan/megebase-identity-service/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}
}

func main() {
	r := gin.Default()

	db, err := database.Connect()
	if err != nil {
		panic(err.Error())
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", handlers.Register(db))
		authRoutes.POST("/login", handlers.Login(db))
	}

	r.Run(":8080")
}
