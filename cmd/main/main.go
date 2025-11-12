package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/leonibeldev/askme/db"
	adminRoutes "github.com/leonibeldev/askme/internal/routes/admin"
	authRoutes "github.com/leonibeldev/askme/internal/routes/auth"
	"github.com/leonibeldev/askme/internal/routes/blog"
	"github.com/leonibeldev/askme/internal/routes/newsletter"
	"github.com/leonibeldev/askme/pkg/utils/functions"

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
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gin.SetMode(gin.DebugMode)
	app := gin.Default()

	// Rate Limiter Middleware
	app.Use(functions.RateLimiter())

	// CORS
	app.Use(functions.Cors())

	// Group for api
	r := app.Group("/api")

	// db connection
	err = db.DataBaseConn()
	if err != nil {
		return
	}
	// create tables if not exist
	err = db.CreateTables()
	if err != nil {
		return
	}

	defer db.Conn.Close()

	// init redis
	db.InitRedis()
	defer db.RedisClient.Close()

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
	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "The route you are looking for does not exist. Please check the URL and try again.",
		})
	})

	app.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":   "Method Not Allowed",
			"message": "Please check the URL and try again.",
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
	err = app.Run(fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return
	}
}
