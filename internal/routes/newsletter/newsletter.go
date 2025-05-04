package newsletter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leonibeldev/askme/internal/controllers"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

// NewUser godoc
// @Summary Subscribe to newsletter
// @Tags Newsletter
// @Accept json
// @Produce json
// @Param data body models.Newsletter true "Newsletter User"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /newsletter [post]
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

// RemoveUser godoc
// @Summary Unsubscribe from newsletter
// @Tags Newsletter
// @Produce json
// @Param uuid path string true "User UUID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /newsletter/{uuid} [get]
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
