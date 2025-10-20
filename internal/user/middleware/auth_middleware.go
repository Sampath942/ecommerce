package middleware

import (
	"errors"
	"net/http"

	"github.com/Sampath942/ecommerce/internal/user/handler"
	"github.com/Sampath942/ecommerce/internal/user/repository"
	"github.com/Sampath942/ecommerce/internal/user/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(h *handler.UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		uid, exists := handler.IsJWTExistingAndValid(authHeader, h)
		if !exists {
			c.Abort()
			c.JSON(http.StatusBadRequest, utils.Response{
				ResponseMessage: "User request failed",
				ResponseDetails: "JWT token doesn't exist or isn't valid",
			})
			return
		}
		user, err := repository.GetUserById(uid, h.DB)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, utils.Response{
					ResponseMessage: "Request Failed",
					ResponseDetails: "The JWt token passed is invalid",
				})
				return
			} else {
				c.JSON(http.StatusInternalServerError, utils.Response{
					ResponseMessage: "Request Failed",
					ResponseDetails: "Querying the database failed. Please try again later",
				})
				return
			}
		}
		c.Set("user", user)
		c.Next()
	}
}