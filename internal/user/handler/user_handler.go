package handler

import (
	"github.com/Sampath942/ecommerce/db"
	"github.com/Sampath942/ecommerce/internal/auth/handler"
)

type UserHandler struct {
	DB          *db.Database
	AuthHandler handler.AuthHandler
}

func NewUserHandler(authHandler handler.AuthHandler) *UserHandler {
	return &UserHandler{
		AuthHandler: authHandler,
	}
}
