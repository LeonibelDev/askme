package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/leonibeldev/askme/db"
	adminRoutes "github.com/leonibeldev/askme/internal/routes/admin"
	authRoutes "github.com/leonibeldev/askme/internal/routes/auth"
	"github.com/leonibeldev/askme/internal/routes/blog"
	"github.com/leonibeldev/askme/internal/routes/newsletter"

	_ "github.com/leonibeldev/askme/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title AskMe API
// @version 1.0
// @description API for authentication, blog management, and newsletter subscription.
// @host localhost:3000
// @BasePath /
func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// db connection
	db.CreateTables()

	// newsletter
	r.POST("/newsletter", newsletter.NewUser)
	r.GET("/newsletter/:uuid", newsletter.RemoveUser)

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
	read.GET("/search", blog.GetPostsByTags)
	read.GET("/all", blog.GetAllPosts)
	read.POST("/new", authRoutes.Handler(), blog.Write)

	// Handle 404 routes
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "The route you are looking for does not exist. Please check the URL and try again.",
			"status":  http.StatusNotFound,
		})
	})

	//Heart
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Run application
	port := os.Getenv("PORT")
	r.Run(fmt.Sprintf(":%s", port))
}
