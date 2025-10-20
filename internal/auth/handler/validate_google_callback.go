package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Sampath942/ecommerce/config"
	"github.com/Sampath942/ecommerce/internal/auth/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func ValidateGoogleCallback(c *gin.Context) {
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

	state, err := c.Cookie("state")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Google authentication failed",
			ResponseDetails: "The state cookie is missing",
		})
		c.Abort()
		return
	}

	stateFromQuery := c.Query("state")
	if state != stateFromQuery {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Google authentication failed",
			ResponseDetails: "Invalid state passed",
		})
		c.Abort()
		return
	}

	code := c.Query("code")
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Google authentication failed",
			ResponseDetails: "Code exchange failed",
		})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Google authentication failed",
			ResponseDetails: "Invalid token being passed to login",
		})
		c.Abort()
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var userInfo map[string]any
	json.Unmarshal(data, &userInfo)
	c.Set("userinfo", userInfo)
	c.Next()
}
