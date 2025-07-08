package main

import (
	"net/http"

	"github.com/MegeKaplan/megebase-identity-service/database"
	"github.com/MegeKaplan/megebase-identity-service/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := database.Connect()

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