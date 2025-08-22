package main

import (
	"net/http"
	"time"

	"github.com/MegeKaplan/megebase-identity-service/database"
	"github.com/MegeKaplan/megebase-identity-service/handlers"
	"github.com/MegeKaplan/megebase-identity-service/messaging"
	"github.com/MegeKaplan/megebase-identity-service/repositories"
	"github.com/MegeKaplan/megebase-identity-service/services"
	"github.com/MegeKaplan/megebase-identity-service/utils"
	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	db, err := database.ConnectPostgres()
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
	refreshTokenRepo := repositories.NewRefreshTokenRepository(redisClient, utils.RefreshTokenTTL())

	authService := services.NewAuthService(userRepo, otpRepo, rabbitMQService, refreshTokenRepo)
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
		authRoutes.POST("/refresh", authHandler.RefreshTokens())
	}

	userRoutes := r.Group("/users")
	// userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("/", userHandler.GetUsers())
		userRoutes.GET("/me", userHandler.GetMe())
		userRoutes.GET("/:id", userHandler.GetUserByID())
		userRoutes.PUT("/:id", userHandler.UpdateUser())
		userRoutes.DELETE("/:id", userHandler.DeleteUser())
	}

	r.Run(":8080")
}
