package main

import (
	"log"

	"github.com/Sampath942/ecommerce/config"
	"github.com/Sampath942/ecommerce/db"
	userhandler "github.com/Sampath942/ecommerce/internal/user/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database, err := db.NewProdDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}
	userHandler := userhandler.UserHandler{
		DB: database,
	}
	r := gin.Default()
	r.GET("/user", userHandler.GetUser)
	r.POST("/user", userHandler.AddUser)
	r.PUT("/user", userHandler.UpdateUser)
	r.DELETE("/user", userHandler.DeleteUser)
	r.Run(":" + config.AppConfig.Port)
}
