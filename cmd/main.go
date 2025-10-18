package main

import (
	"context"
	"log"

	"github.com/Sampath942/ecommerce/config"
	"github.com/Sampath942/ecommerce/db"
	authhandler "github.com/Sampath942/ecommerce/internal/auth/handler"
	"github.com/Sampath942/ecommerce/internal/user/bgworker"
	userhandler "github.com/Sampath942/ecommerce/internal/user/handler"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go bgworker.RemoveExpiredTokens(ctx, &userHandler)
	r := gin.Default()
	r.GET("/verify-email", userHandler.VerifyEmail)
	r.GET("/resend-verification-email", userHandler.AuthMiddleware, userHandler.ResendVerificationEmail)
	r.POST("/login", userHandler.LoginUser)
	r.POST("/user", userHandler.AddUser)
	r.PUT("/user", userHandler.UpdateUser)
	r.DELETE("/user", userHandler.DeleteUser)
	r.Run(":" + config.AppConfig.Port)
}
