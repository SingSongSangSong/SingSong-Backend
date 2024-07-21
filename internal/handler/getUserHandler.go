package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (handler *Handler) GetUser(c *gin.Context) {
	username := c.Param("user")
	user, err := handler.model.GetUser(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    user,
	})
}
