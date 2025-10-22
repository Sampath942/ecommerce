package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Sampath942/ecommerce/internal/user/models"
	"github.com/Sampath942/ecommerce/internal/user/repository"
	"github.com/Sampath942/ecommerce/internal/user/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func IsJWTExistingAndValid(authHeader string, h *UserHandler) (int, bool) {
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
		return -1, false
	}

	authToken := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.AuthHandler.ValidateJWTToken(authToken)
	if err != nil {
		return -1, false
	}
	u, exists := claims["user_id"].(float64)
	uid := int(u)
	if !exists {
		return -1, false
	}
	return uid, true
}

func (h *UserHandler) LoginUser(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	var user models.User
	var jwtToken string
	var err error
	uid, exists := IsJWTExistingAndValid(authHeader, h)
	if !exists {
		var req utils.LoginUserRequest
		if c.ShouldBindJSON(&req) != nil {
			c.JSON(http.StatusBadRequest, utils.Response{
				ResponseMessage: "Request Failed",
				ResponseDetails: "Required fields are not provided in request body",
			})
			return
		}
		user, err = repository.GetUserByCredentials(req.Email, req.Password, h.DB)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() != "unable to query the database" {
				c.JSON(http.StatusBadRequest, utils.Response{
					ResponseMessage: "Request Failed",
					ResponseDetails: "User Credentials are invalid",
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
		jwtToken, err = h.AuthHandler.GenerateJWTToken(user.ID, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				ResponseMessage: "Request Failed",
				ResponseDetails: "Generating the JWT token failed. Please try again later",
			})
			return
		}
	} else {
		user, err = repository.GetUserById(uid, h.DB)
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
		jwtToken = strings.TrimPrefix(authHeader, "Bearer ")
	}

	c.JSON(http.StatusOK, utils.Response{
		ResponseMessage: "Successfully logged in user",
		ResponseDetails: map[string]any{
			"User":      user,
			"JWT Token": jwtToken,
		},
	})
}

func (h *UserHandler) AddUser(c *gin.Context) {
	var addUserReq utils.AddUserRequest
	err := c.ShouldBindJSON(&addUserReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "Request body has some issues. Please check",
		})
		return
	}

	if !utils.IsValidEmail(addUserReq.Email) {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "Email entered is not valid. Please check",
		})
		return
	}

	if !utils.IsValidPhoneNumber(addUserReq.PhoneNumber) {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "Phone number entered is not valid. Please check",
		})
		return
	}

	if !utils.IsValidPassword(addUserReq.Password) {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "Password entered is not strong enough. Please check",
		})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(addUserReq.Password), 10)
	addUserReq.Password = string(hashedPassword)
	user, err := repository.PerformAddUserAndCredentials(addUserReq, h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "Database transaction failed. Try again",
		})
		return
	}
	token, err := GenerateVerificationToken(user.ID, h.DB)
	if err != nil {
		log.Printf("Verification Token cannot be generated. Reason: %s", err.Error())
	} else {
		go SendVerificationEmail(user.Email, token)
	}
	jwt, err := h.AuthHandler.GenerateJWTToken(user.ID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request Partial Success",
			ResponseDetails: "JWT token generation failed but user registration successful. Try logging in",
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		ResponseMessage: "Successfully registered user",
		ResponseDetails: map[string]any{
			"User": user,
			"JWT":  jwt,
		},
	})
}

// TODO: Add actual functionality to update the user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully hit put endpoint",
	})
}

// TODO: Add actual functionality to delete the user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully hit delete endpoint",
	})
}
