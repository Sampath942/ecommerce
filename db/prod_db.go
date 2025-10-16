package db

import (
	"errors"
	"fmt"

	"github.com/Sampath942/ecommerce/config"
	"github.com/Sampath942/ecommerce/internal/user/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewProdDatabase() (*Database, error) {
	dsn := config.AppConfig.DatabaseURL
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &Database{}, errors.New("couldn't get a connection to the database")
	}
	fmt.Println("Automigrating tables")
	result := database.Exec("CREATE SCHEMA IF NOT EXISTS users")
	if result.Error != nil {
		return &Database{}, errors.New("couldn't create a new schema")
	}
	database.AutoMigrate(&models.Credentials{}, &models.User{}, &models.VerificationToken{})
	return &Database{
		DB: database,
	}, nil
}
