package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"regexp"
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
	fmt.Println("JWT is: " + authToken)
	claims, err := h.AuthHandler.ValidateJWTToken(authToken)
	if err != nil {
		fmt.Println("ValidateJWTToken failed")
		return -1, false
	}
	uid_float, exists := claims["user_id"].(float64)
	fmt.Println("claims are: ", claims)
	fmt.Println("uid is: ", uid_float)
	if !exists {
		fmt.Println("uid doesn't exist in claims fail")
		return -1, false
	}
	uid := int(uid_float)
	return uid, true
}

func isValidPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}

	upper := regexp.MustCompile(`[A-Z]`)
	lower := regexp.MustCompile(`[a-z]`)
	digit := regexp.MustCompile(`[0-9]`)
	special := regexp.MustCompile(`[!@#$&*]`)

	// At least 2 uppercase, 3 lowercase, 2 digits, and 1 special char
	return upper.MatchString(pw) &&
		lower.MatchString(pw) &&
		digit.MatchString(pw) &&
		special.MatchString(pw)
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


// TODO: Add verify phone number step similar to verify email step.
// TODO: Add resend verification for mobile number endpoints
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

	_, err = mail.ParseAddress(addUserReq.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "Email entered is not valid. Please check",
		})
		return
	}

	re := regexp.MustCompile(`((\+*)((0[ -]*)*|((91 )*))((\d{12})+|(\d{10})+))|\d{5}([- ]*)\d{6}`)
    if !re.MatchString(addUserReq.PhoneNumber) {
        c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "Phone number entered is not valid. Please check",
		})
   		return
	}

	if !isValidPassword(addUserReq.Password) {
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
