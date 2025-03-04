package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/leonibeldev/askme/db"
	routes "github.com/leonibeldev/askme/internal/routes/auth"
	"github.com/leonibeldev/askme/internal/routes/blog"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// db connection
	db.VerifyIfDBExist()

	// Routes for API Auth
	auth := r.Group("/api/v1")

	auth.POST("/signup", routes.Signup)
	auth.POST("/login", routes.Login)
	auth.GET("/home", routes.Handler(), routes.Home) // secure router 🔒

	// Routes for API Blog
	read := r.Group("/read")

	read.GET("/:id", blog.Read)

	// Handle 404 routes
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "The route you are looking for does not exist. Please check the URL and try again.",
			"status":  http.StatusNotFound,
		})
	})

	// Run application
	port := os.Getenv("PORT")
	r.Run(fmt.Sprintf(":%s", port))
}
