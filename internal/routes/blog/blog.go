package blog

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leonibeldev/askme/internal/controllers"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

// Read godoc
// @Summary Get a single blog post by ID
// @Tags Blog
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /blog/{id} [get]
func Read(c *gin.Context) {
	id, _ := c.Params.Get("id")

	post, _ := controllers.GetOnePostFromDB(id)

	c.JSON(200, gin.H{
		"post": post,
	})
}

// Write godoc
// @Summary Create a new blog post
// @Tags Blog
// @Accept json
// @Produce json
// @Param post body models.Post true "Blog Post"
// @Success 200 models.Post{} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @Router /blog/new [post]
func Write(c *gin.Context) {

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// set post date
	post.Date = time.Now()

	// set post author
	author, _ := c.Get("email")
	post.Author = author.(string)

	// save post
	postId, err := controllers.SavePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return response
	c.JSON(http.StatusOK, gin.H{
		"postId": postId,
	})

}

// GetAllPosts godoc
// @Summary Get all blog posts with optional pagination
// @Tags Blog
// @Produce json
// @Param offset query string false "Offset for pagination"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /blog [get]
func GetAllPosts(c *gin.Context) {
	// get page
	offset := c.Query("offset")
	// limit := c.Query("limit")

	// get all posts from db
	posts, err := controllers.GetAllPostsFromDB(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})

}

// GetPostsByTags godoc
// @Summary Get blog posts filtered by tag
// @Tags Blog
// @Produce json
// @Param tag query string true "Tag to filter posts"
// @Success 200 models.Post{} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /blog/tag [get]
func GetPostsByTags(c *gin.Context) {
	tag := c.Query("tag")

	if tag == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Query 'tag' is required",
		})
		return
	}

	posts, err := controllers.GetPostsByTags(tag)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Posts not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}
