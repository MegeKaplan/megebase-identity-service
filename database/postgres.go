package database

import (
	"fmt"
	"os"

	"github.com/MegeKaplan/megebase-identity-service/models"
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, response.ErrDBConnection
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, response.ErrDBMigration
	}

	return db, nil
}
