package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sampath942/ecommerce/config"
	"github.com/Sampath942/ecommerce/internal/user/models"
	"github.com/Sampath942/ecommerce/internal/user/repository"
	"github.com/Sampath942/ecommerce/internal/user/utils"
	"github.com/gin-gonic/gin"
)

func (h *UserHandler) RequestOTP(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "The user cannot be found. Try later",
		})
		return
	}
	user := u.(models.User)
	if user.IsMobileVerified {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "The user's mobile number is already verified",
		})
		return
	}
	_, err := repository.GetValueRedis(context.Background(), *user.PhoneNumber, h.DB.Redis)
	if err == nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "An OTP is already sent. Please wait for 5 minutes before sending next request",
		})
		return
	}
	otp := utils.GenerateOTP()
	fmt.Println("OTP generated: ", otp)
	status := repository.SetKeyValueRedis(context.Background(), *user.PhoneNumber, otp, 5*time.Minute, h.DB.Redis)
	if status.Err() != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "Caching the otp failed",
		})
		return
	}
	message := "" + otp
	url := "https://api.smsmobileapi.com/sendsms?apikey=" + config.AppConfig.SMSAPIKey + "&recipients=" + *user.PhoneNumber + "&message=" + message
	fmt.Println(url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := http.DefaultClient.Do(request)

	if err != nil || resp.StatusCode != http.StatusOK{
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "The API for sending SMS failed",
		})
		return
	}

	defer resp.Body.Close()
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println(resp.Status)
	fmt.Println("Response body is: ", result)
	c.JSON(http.StatusOK, utils.Response{
		ResponseMessage: "Request successful",
		ResponseDetails: "OTP sent successfully. Please check your mobile",
	})
}

func (h *UserHandler) VerifyOTP(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "The user cannot be found. Try later",
		})
		return
	}
	user := u.(models.User)
	if user.IsMobileVerified {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "The user's mobile number is already verified",
		})
		return
	}
	var req utils.VerifyOTPRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "The otp is not passed in the right format",
		})
		return
	}
	storedOTP, err := repository.GetValueRedis(context.Background(), *user.PhoneNumber, h.DB.Redis)
	if err != nil || storedOTP != req.OTP {
		c.JSON(http.StatusBadRequest, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "The otp entered is invalid",
		})
		return
	}
	err = repository.SetUserPhoneVerified(user, h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "Updating mobile verification status failed. Try later",
		})
		return
	}
	deleted, err := repository.DeleteKeyRedis(context.Background(), *user.PhoneNumber, h.DB.Redis)
	if deleted == 0 || err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			ResponseMessage: "Request failed",
			ResponseDetails: "Internal error occured while clearing the cache. Try later",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		ResponseMessage: "Request Successful",
		ResponseDetails: "User phone number status set to verified",
	})
}
