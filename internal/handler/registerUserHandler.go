package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRegisterRequest struct {
	Username string
}

func (handler *Handler) RegisterUser(c *gin.Context) {
	request := &UserRegisterRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := handler.model.RegisterUser(request.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
