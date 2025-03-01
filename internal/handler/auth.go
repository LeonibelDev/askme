package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "auth",
	})
}
