package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/leonibeldev/askme/internal/controllers"
	"github.com/leonibeldev/askme/pkg/utils/functions"
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
			c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Success: false,
				Message: "No token provide in header (Authorization)",
			})
			return
		}

		AuthorizationToken := strings.Split(c.GetHeader("Authorization"), " ")[1]

		// validate token
		claims, err := token.GetClaims(AuthorizationToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ResponseMessage{
				Success: false,
				Message: "Invalid token in header (Authorization)",
			})
			return
		}

		//return claims
		c.Set("email", claims["email"])
		c.Set("fullname", claims["fullname"])
		c.Set("username", claims["username"])
		c.Next()
	}
}

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.Login true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /auth/login [post]
func Login(c *gin.Context) {

	// parse user data input
	var LoginValues models.Login

	if err := c.ShouldBindJSON(&LoginValues); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Success: false,
			Message: "Data binding error",
			Error:   err.Error(),
		})
		return
	}

	LoginValues.Email = strings.ToLower(LoginValues.Email)

	// validate if user exist
	dbUser, err := controllers.GetUser(LoginValues.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// compare password
	matchingPassword := hash.CheckPasswordHash(LoginValues.Password, dbUser.Password)
	if !matchingPassword {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// generate toke if all information is correct ⚠
	stringToken, _ := token.GenerateToken(dbUser.Email, dbUser.Username, dbUser.Fullname)

	fmt.Println(dbUser.Fullname)

	c.JSON(http.StatusOK, models.ResponseMessage{
		Success: true,
		Data:    "Bearer " + stringToken,
		Message: dbUser.Fullname,
	})
}

// Signup godoc
// @Summary Register new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.DBUser true "User registration info"
// @Success 200 {object} models.ResponseMessage{data=string}
// @Failure 400 {object} models.ResponseMessage
//
//	@Example 200 {json} SuccessResponse {
//	  "success": true,
//	  "message": "Welcome, John Doe",
//	  "data": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//	}
//
// @Router /auth/signup [post]
func Signup(c *gin.Context) {

	var userData models.DBUser

	// verify if all data is comming
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Success: false,
			Message: "All field are required to register",
			Error:   err.Error(),
		})
		return
	}

	userData.Email = strings.ToLower(userData.Email)

	// verify if user exist
	user := controllers.UserExist(userData.Email)
	if user {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Success: false,
			Message: "User already exist",
		})
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
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Success: false,
			Message: "Error hashing password",
			Error:   err.Error(),
		})
		return
	}

	userData.Password = string(hash)

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(originalPassword))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Success: false,
			Message: "Error comparing password",
			Error:   err.Error(),
		})
		return
	}

	// set username
	randomId, _ := functions.RandomNumber()
	userData.Username = strings.ToLower(fmt.Sprintf("%s_%d", strings.Split(userData.Fullname, " ")[0], randomId))

	// set time
	userData.Created_at = time.Now()

	// create user
	controllers.CreateUser(userData)

	// generate token
	token, err := token.GenerateToken(userData.Email, userData.Username, userData.Fullname)
	if err != nil {
		panic(err)
	}

	fmt.Println(userData.Fullname)

	// return token and username
	c.JSON(http.StatusCreated, models.ResponseMessage{
		Success: true,
		Data:    "Bearer " + token,
		Message: userData.Fullname,
	})
}
