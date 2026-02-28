// @title           mytechblog API
// @version         1.0
// @description     API for mytechblog platform.
// @termsOfService  work in progress

// @contact.name   LeonibelDev
// @contact.email  leonibel.ramirez@gmail.com

// @host      localhost:3000
// @BasePath  /api

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/leonibeldev/askme/db"
	adminRoutes "github.com/leonibeldev/askme/internal/routes/admin"
	authRoutes "github.com/leonibeldev/askme/internal/routes/auth"
	"github.com/leonibeldev/askme/internal/routes/blog"
	"github.com/leonibeldev/askme/internal/routes/newsletter"
	"github.com/leonibeldev/askme/pkg/utils/functions"
	"github.com/leonibeldev/askme/pkg/utils/models"

	_ "github.com/leonibeldev/askme/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	// load .env
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init Uptrace
	/*uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(fmt.Sprintf("https://%s@api.uptrace.dev?grpc=4317", os.Getenv("OPENTRACE_GRPC_ADDR"))),
	)*/

	gin.SetMode(gin.DebugMode)
	app := gin.Default()

	// Metrics
	//app.Use(otelgin.Middleware("mytechblog"))

	// Rate Limiter Middleware
	app.Use(functions.RateLimiter())

	// CORS
	app.Use(functions.Cors())

	// Group for api
	r := app.Group("/api")

	/******************************
	*	Databases Connections
	*******************************/

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

	/******************************
	*	Databases Connections
	*******************************/

	// newsletter
	r.POST("/newsletter", newsletter.NewUser)
	r.GET("/newsletter/:uuid", newsletter.RemoveUser)

	// Routes for Auth API
	auth := r.Group("/auth")

	auth.POST("/signup", authRoutes.Signup)
	auth.POST("/login", authRoutes.Login)

	// Secure Routes
	admin := r.Group("/admin")
	admin.Use(authRoutes.Handler())

	admin.GET("/home", adminRoutes.Home)
	admin.GET("/profile", adminRoutes.User)
	admin.PUT("/profile", adminRoutes.UpdateProfile)

	// Routes for Blog API and Portfolio
	posts := r.Group("/posts")

	posts.GET("/", blog.GetAllPosts)
	posts.GET("/top", blog.GetTopPosts)
	posts.GET("/:id", blog.Read)
	posts.PUT("/:id", blog.UpdatePost)
	posts.GET("/tag/:tag", blog.GetPostsByTags)

	posts.POST("/", authRoutes.Handler(), blog.Write)
	posts.GET("/by/:author", blog.GetPostsByAuthor)

	// Handle 404 routes
	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Success: false,
			Message: "Resource Not Found. Please check the URL and try again.",
		})
	})

	app.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, models.ResponseMessage{
			Success: false,
			Message: "Method Not Allowed. Please check the request method and try again.",
		})
	})

	//Heart
	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run application

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", os.Getenv("PORT")),
		Handler: app,
	}

	go func() {
		log.Printf("Server on port %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	shutdown, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	<-shutdown.Done()

	log.Println("Shutingdown in progress...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Println("Forced shutdown")
	}
	log.Println("Server down successfully")

	/*
		port := os.Getenv("PORT")

		err = app.Run(fmt.Sprintf("0.0.0.0:%s", port))
		if err != nil {
			return
		}*/
}
