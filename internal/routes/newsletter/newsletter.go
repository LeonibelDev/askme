package newsletter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leonibeldev/askme/internal/controllers"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

func NewUser(c *gin.Context) {

	var user models.Newsletter

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := controllers.AddUserNewsletter(user.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user exist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user added to newsletter",
	})

}

func RemoveUser(c *gin.Context) {

	uuid, _ := c.Params.Get("uuid")

	if err := controllers.RemoveUserNewsletter(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user removed from newsletter",
	})
}
