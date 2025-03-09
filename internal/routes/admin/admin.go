package adminRoutes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leonibeldev/askme/internal/controllers"
)

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

func User(c *gin.Context) {
	userEmail, exist := c.Get("email")
	if !exist {
		return
	}

	// get user from db
	user, _ := controllers.GetUser(fmt.Sprintf("%s", userEmail))

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
