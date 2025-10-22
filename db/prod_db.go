package db

import (
	"errors"
	"strconv"

	"github.com/Sampath942/ecommerce/config"
	"github.com/Sampath942/ecommerce/internal/user/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initRedis() (*redis.Client, error) {
	db, err := strconv.Atoi(config.AppConfig.RedisDB)
	if err != nil {
		return &redis.Client{}, errors.New("wrong database passed in REDIS_DB parameter")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisURL,
		Password: config.AppConfig.RedisPassword,
		DB:       db,
	})
	return client, nil
}

func NewProdDatabase() (*Database, error) {
	dsn := config.AppConfig.DatabaseURL
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &Database{}, errors.New("couldn't get a connection to the database")
	}
	result := database.Exec("CREATE SCHEMA IF NOT EXISTS users")
	if result.Error != nil {
		return &Database{}, errors.New("couldn't create a new schema")
	}
	database.AutoMigrate(&models.Credentials{}, &models.User{}, &models.VerificationToken{})
	redisClient, err := initRedis()
	if err != nil {
		return &Database{}, err
	}
	return &Database{
		DB:    database,
		Redis: redisClient,
	}, nil
}
