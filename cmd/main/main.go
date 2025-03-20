package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/leonibeldev/askme/db"
	adminRoutes "github.com/leonibeldev/askme/internal/routes/admin"
	authRoutes "github.com/leonibeldev/askme/internal/routes/auth"
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
	db.CreateTables()

	// Routes for Home / portfolio
	r.GET("/github/:username", blog.GetGitHubRepos)

	// Routes for Auth API
	auth := r.Group("/auth")

	auth.POST("/signup", authRoutes.Signup)
	auth.POST("/login", authRoutes.Login)

	// Routes for Admin 🔒
	admin := r.Group("/admin")
	admin.Use(authRoutes.Handler())

	admin.GET("/home", adminRoutes.Home) // secure routes 🔒
	admin.GET("/user", adminRoutes.User)

	// Routes for Blog API and Portfolio
	read := r.Group("/blog")

	read.GET("/:id", blog.Read)
	read.POST("/new", authRoutes.Handler(), blog.Write)

	// Handle 404 routes
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "The route you are looking for does not exist. Please check the URL and try again.",
			"status":  http.StatusNotFound,
		})
	})

	// Run application
	//port := os.Getenv("PORT")
	r.Run(fmt.Sprintf(":%s", "3000"))
}
