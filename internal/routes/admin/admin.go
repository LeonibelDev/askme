package adminRoutes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leonibeldev/askme/internal/controllers"
)

// Home godoc
// @Summary Admin home
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]string
// @Router /admin/home [get]
// @Security ApiKeyAuth
func Home(c *gin.Context) {
	userEmail, exist := c.Get("email")
	if !exist {
		return
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf("Hello, %s", userEmail))

	c.JSON(http.StatusOK, gin.H{
		"message": message.String(),
	})
}

// User godoc
// @Summary Get admin user
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/user [get]
// @Security ApiKeyAuth
func User(c *gin.Context) {
	userEmail, exist := c.Get("email")
	if !exist {
		return
	}

	// get user from db
	user, _ := controllers.GetUser(userEmail.(string))

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
