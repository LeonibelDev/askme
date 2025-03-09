package routes

import (
	"net/http"
	"strings"

	"github.com/leonibeldev/askme/internal/controllers"
	"github.com/leonibeldev/askme/pkg/utils/hash"
	"github.com/leonibeldev/askme/pkg/utils/models"
	"github.com/leonibeldev/askme/pkg/utils/token"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

/*
middleware to validate token before
access to secure routes
*/
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get header (authorization)
		if len(c.GetHeader("Authorization")) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provide (Unauthorized)",
			})
			return
		}

		AuthorizationToken := strings.Split(c.GetHeader("Authorization"), " ")[1]

		// validate token
		claims, err := token.GetClaims(AuthorizationToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid token in header (Authorization)",
			})
			return
		}

		//return claims
		c.Set("email", claims["email"])
		c.Next()
	}
}

func Login(c *gin.Context) {

	// parse user data input
	var LoginValues models.Login
	if err := c.ShouldBindJSON(&LoginValues); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate if user exist
	dbUser, err := controllers.GetUser(LoginValues.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	// compare password
	matchingPassword := hash.CheckPasswordHash(LoginValues.Password, dbUser.HashPassword)
	if !matchingPassword {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "password not matching",
		})
		return
	}

	// generate toke if all information is correct ⚠
	stringToken, _ := token.GenerateToken(dbUser.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": "Bearer " + stringToken,
	})
}

func Signup(c *gin.Context) {

	var userData models.User

	// verify if all data is comming
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*
	 check fields information, to validate not insert SQLInjection in DB
	 write: ⚠
	*/

	originalPassword := userData.Password

	// hash password
	hash, err := hash.HashPassword(userData.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData.Password = string(hash)

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(originalPassword))
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

	// generate token
	token, err := token.GenerateToken(userData.Email)
	if err != nil {
		panic(err)
	}

	// return user data and comparation
	c.JSON(http.StatusOK, gin.H{
		"token": "Bearer " + token,
	})
}
