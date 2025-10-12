package db

import (
	"errors"

	"github.com/Sampath942/ecommerce/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewProdDatabase() (*Database, error) {
	dsn := config.AppConfig.DatabaseURL
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &Database{}, errors.New("couldn't get a connection to the database")
	}
	database.AutoMigrate()
	return &Database{
		DB: database,
	}, nil
}
