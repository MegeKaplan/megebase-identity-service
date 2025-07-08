package database

import (
	"fmt"
	"os"

	"github.com/MegeKaplan/megebase-identity-service/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	return db
}
