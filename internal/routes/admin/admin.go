package adminRoutes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leonibeldev/askme/internal/controllers"
	"github.com/leonibeldev/askme/pkg/utils/models"
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
// @Success 200 {object} models.ResponseMessage
// @Router /admin/user [get]
// @Security ApiKeyAuth
func User(c *gin.Context) {
	userEmail, exist := c.Get("email")
	if !exist {
		return
	}

	// get user from db
	user, _ := controllers.GetUserProfile(userEmail.(string))

	c.JSON(http.StatusOK, models.ResponseMessage{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// UpdateProfile godoc
// @Summary Update admin user profile
// @Tags Admin
// @Accept json
// @Produce json
// @Param user body models.UpdateUserProfile true "User Profile Data"
// @Success 200 {object} models.ResponseMessage
// @Router /admin/profile [put]
// @Security ApiKeyAuth
func UpdateProfile(c *gin.Context) {
	userEmail, exist := c.Get("email")
	if !exist {
		return
	}

	var ProfileData models.Profile
	if err := c.ShouldBind(&ProfileData); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	// update user profile in db
	updatedUser, err := controllers.UpdateProfile(userEmail.(string), ProfileData)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Success: false,
			Message: "Error to update user",
			Error:   err.Error(),
		})
	}

	c.JSON(http.StatusCreated, models.ResponseMessage{
		Success: true,
		Message: updatedUser,
	})
}
