package blog

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Read(c *gin.Context) {

	id, _ := c.Params.Get("id")

	c.JSON(200, gin.H{
		"Message": "Welcome to askme.dev API *",
		"id":      id,
		"status":  http.StatusOK,
	})
}
