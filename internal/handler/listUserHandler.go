package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (handler *Handler) ListUser(c *gin.Context) {
	users, err := handler.model.ListUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    users,
	})
}
