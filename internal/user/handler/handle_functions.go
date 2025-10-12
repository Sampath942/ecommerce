package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler)GetUser(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully hit get endpoint",
	})
}

func (h *UserHandler)AddUser(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully hit post endpoint",
	})
}

func (h *UserHandler)UpdateUser(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully hit put endpoint",
	})
}

func (h *UserHandler)DeleteUser(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully hit delete endpoint",
	})
}
