package db

import (
	"errors"

	"github.com/Sampath942/ecommerce/config"
	"github.com/Sampath942/ecommerce/internal/user/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewTestDatabase() (*Database, error) {
	dsn := config.AppConfig.TestDatabaseURL
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &Database{}, errors.New("couldn't get a connection to the database")
	}
	database.AutoMigrate(&models.User{}, &models.Credentials{}, &models.VerificationToken{})
	return &Database{
		DB: database,
	}, nil
}
