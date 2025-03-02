package routes

import (
	"net/http"

	"github.com/leonibeldev/askme/internal/controllers"
	"github.com/leonibeldev/askme/pkg/utils/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Handler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"Message": "Welcome to askme.dev API *",
		"status":  http.StatusOK,
	})
}

func Login(c *gin.Context) {

	var LoginValues models.Login
	if err := c.ShouldBindJSON(&LoginValues); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginValues)
}

func Signup(c *gin.Context) {

	var userData models.User
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	originaPassword := userData.Password

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userData.Password = string(hash)

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(originaPassword))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verify if user exist
	user := controllers.UserExist(userData.Email)
	if user {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exist"})
		return
	}

	// create user
	controllers.CreateUser(userData)

	// return user data and comparation
	c.JSON(http.StatusOK, gin.H{
		"user":    userData,
		"compare": (err == nil),
	})
}
