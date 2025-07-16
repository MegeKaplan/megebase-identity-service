package main

import (
	"net/http"

	"github.com/MegeKaplan/megebase-identity-service/database"
	"github.com/MegeKaplan/megebase-identity-service/handlers"
	"github.com/MegeKaplan/megebase-identity-service/repositories"
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

	userRepo := repositories.NewUserGormRepository(db)
	otpRepo := repositories.NewInMemoryOTPRepository()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", handlers.Register(userRepo, otpRepo))
		authRoutes.POST("/register/send-otp", handlers.SendOTP(otpRepo))
		authRoutes.POST("/register/verify-otp", handlers.VerifyOTP(otpRepo))
		authRoutes.POST("/login", handlers.Login(userRepo))
	}

	r.Run(":8080")
}
