package main

import (
	"net/http"
	"time"

	"github.com/MegeKaplan/megebase-identity-service/database"
	"github.com/MegeKaplan/megebase-identity-service/handlers"
	"github.com/MegeKaplan/megebase-identity-service/messaging"
	"github.com/MegeKaplan/megebase-identity-service/middleware"
	"github.com/MegeKaplan/megebase-identity-service/repositories"
	"github.com/MegeKaplan/megebase-identity-service/services"
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

	rabbitMQService, err := messaging.NewRabbitMQService("megebase.topic")
	if err != nil {
		panic(err.Error())
	}
	defer rabbitMQService.Close()

	redisClient, err := database.ConnectRedis()
	if err != nil {
		panic(err.Error())
	}
	
	userRepo := repositories.NewUserGormRepository(db)
	// otpRepo := repositories.NewInMemoryOTPRepository()
	otpRepo := repositories.NewRedisOTPRepository(redisClient, 5*time.Minute)

	authService := services.NewAuthService(userRepo, otpRepo, rabbitMQService)
	userService := services.NewUserService(userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register())
		authRoutes.POST("/register/send-otp", authHandler.SendOTP())
		authRoutes.POST("/register/verify-otp", authHandler.VerifyOTP())
		authRoutes.POST("/login", authHandler.Login())
	}

	userRoutes := r.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("/", userHandler.GetUsers())
		userRoutes.GET("/me", userHandler.GetMe())
		userRoutes.GET("/:id", userHandler.GetUserByID())
		userRoutes.PUT("/:id", userHandler.UpdateUser())
		userRoutes.DELETE("/:id", userHandler.DeleteUser())
	}

	r.Run(":8080")
}
