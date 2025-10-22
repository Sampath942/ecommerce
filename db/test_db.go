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

func initTestRedis() (*redis.Client, error) {
	db, err := strconv.Atoi(config.AppConfig.RedisTestDB)
	if err != nil {
		return &redis.Client{}, errors.New("wrong database passed in REDIS_DB parameter")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisTestURL,
		Password: config.AppConfig.RedisTestPassword,
		DB:       db,
	})
	return client, nil
}

func NewTestDatabase() (*Database, error) {
	dsn := config.AppConfig.DatabaseURLTest
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &Database{}, errors.New("couldn't get a connection to the database")
	}
	database.AutoMigrate(&models.User{}, &models.Credentials{}, &models.VerificationToken{})
	redisClient, err := initTestRedis()
	if err != nil {
		return &Database{}, err
	}
	return &Database{
		DB: database,
		Redis: redisClient,
	}, nil
}
