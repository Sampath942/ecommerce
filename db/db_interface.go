package db

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
	Redis *redis.Client
}
