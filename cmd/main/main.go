package main

import (
	"fmt"
	"log"
	"os"

	routes "github.com/leonibeldev/askme/internal/routes/auth"
	"github.com/leonibeldev/askme/internal/routes/blog"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// Routes for API Auth
	auth := r.Group("/api/v1")

	auth.GET("/", routes.Handler)
	auth.POST("/signup", routes.Signup)
	auth.POST("/login", routes.Login)

	// Routes for API Blog
	read := r.Group("/read")

	read.GET("/:id", blog.Read)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("defaulting to port %s", port)
	}

	r.Run(fmt.Sprintf(":%s", port))
}
