package main

import (
	"net/http"

	"github.com/MegeKaplan/megebase-identity-service/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", handlers.Login())
		authRoutes.POST("/register", handlers.Register())
	}
	
	r.Run(":8080")
}