package main

import (
	"context"
	"log"

	"github.com/Sampath942/ecommerce/config"
	"github.com/Sampath942/ecommerce/db"
	authhandler "github.com/Sampath942/ecommerce/internal/auth/handler"
	"github.com/Sampath942/ecommerce/internal/user/bgworker"
	userhandler "github.com/Sampath942/ecommerce/internal/user/handler"
	"github.com/Sampath942/ecommerce/internal/user/middleware"
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
		DB:          database,
		AuthHandler: authHandler,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go bgworker.RemoveExpiredTokens(ctx, &userHandler)
	r := gin.Default()

	// TODO: Add unit tests to test all endpoints
	// TODO: See ways to add redis caching
	// TODO: Currently, every user is registered as customer, there should be an option for a user to promote himself
	// 		 as admin. This can be done only if user does not have any order placed as a customer at any point in time.
	// 		 Once the user clicks on promote to admin, we check this condition and then we proceed to add him as admin.
	//		 For adding or removing an admin, there needs to be a superadmin who gets the notification to promote a user
	// 		 to admin and then the super admin clicks approve, the user becomes amin user.
	r.GET("/verify-email", userHandler.VerifyEmail)
	r.GET("/resend-verification-email", middleware.AuthMiddleware(&userHandler), userHandler.ResendVerificationEmail)
	r.POST("/login", userHandler.LoginUser)
	r.POST("/user", userHandler.AddUser)
	r.PUT("/user", userHandler.UpdateUser)
	r.DELETE("/user", userHandler.DeleteUser)
	r.GET("/auth/google/login", userHandler.GoogleLogin)
	r.GET("/auth/google/callback", authhandler.ValidateGoogleCallback, userHandler.GoogleCallback)
	r.Run(":" + config.AppConfig.Port)
}
