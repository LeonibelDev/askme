package blog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

func Read(c *gin.Context) {

	id, _ := c.Params.Get("id")

	c.JSON(200, gin.H{
		"Message": "Welcome to askme.dev API *",
		"id":      id,
		"status":  http.StatusOK,
	})
}

func Write(c *gin.Context) {

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"Message": "Welcome to askme.dev API *",
		"data":    post,
	})

}
