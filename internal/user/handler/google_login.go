package handler

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/Sampath942/ecommerce/config"
	"github.com/Sampath942/ecommerce/internal/user/repository"
	"github.com/Sampath942/ecommerce/internal/user/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (h *UserHandler) GoogleLogin(c *gin.Context) {
	var googleOAuthConfig = &oauth2.Config{
		RedirectURL:  config.AppConfig.OAuthRedirectURL,
		ClientID:     config.AppConfig.OAuthClientID,
		ClientSecret: config.AppConfig.OAuthClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "Generating state failed. Please try again later",
		})
		return
	}
	c.SetCookie("state", state, 3600, "/", "localhost", false, true)
	url := googleOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *UserHandler) GoogleCallback(c *gin.Context) {
	userInfo, _ := c.Get("userinfo")
	userInfoMap := userInfo.(map[string]any)
	email := userInfoMap["email"].(string)
	name := userInfoMap["name"].(string)
	googleID := userInfoMap["id"].(string)
	user, err := repository.GetUserByGoogleID(googleID, h.DB)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Register new user if email registered user is not found.
			user, err = repository.GetUserByEmail(email, h.DB)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Register new user now
					err = repository.AddUserFromGoogleID(googleID, email, name, h.DB)
					if err != nil {
						c.JSON(http.StatusInternalServerError, utils.Response{
							ResponseMessage: "Google signin/signup failed",
							ResponseDetails: "Adding user from google id failed",
						})
						return
					}
				} else {
					// Error out in this case
					c.JSON(http.StatusInternalServerError, utils.Response{
						ResponseMessage: "Google signin/signup failed",
						ResponseDetails: "Database query failed. Try again later",
					})
					return
				}
			} else {
				// Login the user and update the googleID
				err = repository.UpdateUserGoogleID(&user, googleID, h.DB)
				if err != nil {
					c.JSON(http.StatusInternalServerError, utils.Response{
						ResponseMessage: "Google signin/signup failed",
						ResponseDetails: "Database updation failed. Try again later",
					})
					return
				}
			}

		} else {
			// Error out in this case
			c.JSON(http.StatusInternalServerError, utils.Response{
				ResponseMessage: "Google signin/signup failed",
				ResponseDetails: "Database query failed. Try again later",
			})
			return
		}
	}
	// Generate JWT token
	token, err := h.AuthHandler.GenerateJWTToken(user.ID, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Google signin/signup failed",
			ResponseDetails: "JWT token generation failed",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		ResponseMessage: "Successfully logged in user",
		ResponseDetails: map[string]any{
			"User":      user,
			"JWT Token": token,
		},
	})
}
