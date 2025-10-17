package main

import (
	"log"

	"github.com/Sampath942/ecommerce/config"
	"github.com/Sampath942/ecommerce/db"
	userhandler "github.com/Sampath942/ecommerce/internal/user/handler"
	authhandler "github.com/Sampath942/ecommerce/internal/auth/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database, err := db.NewProdDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}
	authHandler := authhandler.NewAuthHandler() 
	userHandler := userhandler.UserHandler{
		DB: database,
		AuthHandler: authHandler,
	}
	r := gin.Default()
	r.GET("/verify-email", userHandler.VerifyEmail)
	r.GET("/resend-verification-email", userHandler.AuthMiddleware, userHandler.ResendVerificationEmail)
	r.POST("/login", userHandler.LoginUser)
	r.POST("/user", userHandler.AddUser)
	r.PUT("/user", userHandler.UpdateUser)
	r.DELETE("/user", userHandler.DeleteUser)
	r.Run(":" + config.AppConfig.Port)
}
