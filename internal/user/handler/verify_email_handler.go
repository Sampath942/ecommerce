package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Sampath942/ecommerce/config"
	"github.com/Sampath942/ecommerce/db"
	"github.com/Sampath942/ecommerce/internal/user/models"
	"github.com/Sampath942/ecommerce/internal/user/repository"
	"github.com/Sampath942/ecommerce/internal/user/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func SendVerificationEmail(toEmail, token string) {
	link := fmt.Sprintf("http://localhost:%s/verify-email?token=%s", config.AppConfig.Port, token)
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.OtpSendingEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Verify your email")
	m.SetBody("text/plain", fmt.Sprintf("Please click on the following link to verify your email %s", link))
	d := gomail.NewDialer("smtp.gmail.com", 587, config.AppConfig.OtpSendingEmail, config.AppConfig.OtpSendingEmailPassword)
	err := d.DialAndSend(m)
	if err != nil {
		log.Printf("Sending a mail has failed. Error is: %s", err.Error())
	}
}

func GenerateVerificationToken(userid int, database *db.Database) (string, error) {
	token := uuid.NewString()
	verification_token := models.VerificationToken{
		UserID:    userid,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
		Used:      false,
	}
	if repository.AddVerificationToken(verification_token, database) != nil {
		return "", errors.New("error adding a token to the database")
	}
	return token, nil
}

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "The passed token is empty",
		})
		return
	}
	verificationToken, err := repository.GetVerificationDetailsFromToken(token, h.DB)
	fmt.Println(token, verificationToken)
	fmt.Println(err)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, utils.Response{
				ResponseMessage: "Request Failed",
				ResponseDetails: "The passed token is invalid",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, utils.Response{
				ResponseMessage: "Request Failed",
				ResponseDetails: "Unable to query the database, try later",
			})
			return
		}
	}
	user, err := repository.GetUserById(verificationToken.UserID, h.DB)
	fmt.Println("user is: ", user)
	fmt.Println("err is: ", err)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, utils.Response{
				ResponseMessage: "Request Failed",
				ResponseDetails: "The passed token is invalid",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, utils.Response{
				ResponseMessage: "Request Failed",
				ResponseDetails: "Unable to query the database, try later",
			})
			return
		}
	}
	if user.IsEmailVerified {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "The user email is already verified. No need to verify again",
		})
		return
	}
	if verificationToken.Used {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "The passed token is already used",
		})
		return
	}
	if time.Now().After(verificationToken.ExpiresAt) {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "The passed token is expired",
		})
		return
	}
	err = repository.SetVerificationTokenToUsed(verificationToken, h.DB)
	// fmt.Println("Error is: ", err.Error())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "Unable to update status of the token, try later",
		})
		return
	}
	err = repository.SetUserEmailVerified(user, h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "Unable to update the email verified status, try later",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		ResponseMessage: "Request is successful",
		ResponseDetails: "Email successfully verified",
	})
}

func (h *UserHandler) ResendVerificationEmail(c *gin.Context) {
	u, _ := c.Get("user")
	user := u.(models.User)
	if user.IsEmailVerified {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "The user is already verified. No need to verify again",
		})
		return
	}
	_, err := repository.GetValidVerificationDetailsFromUserID(user.ID, h.DB)
	if err == nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "A token is already passed. Check your inbox",
		})
		return
	}
	token, err := GenerateVerificationToken(user.ID, h.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request Failed",
			ResponseDetails: "Verification Token cannot be generated.",
		})
		return
	} else {
		go SendVerificationEmail(user.Email, token)
	}
	c.JSON(http.StatusOK, utils.Response{
		ResponseMessage: "Request successful",
		ResponseDetails: "Please check your email for further instructions",
	})
}
