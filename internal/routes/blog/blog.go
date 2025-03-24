package blog

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leonibeldev/askme/internal/controllers"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

func Read(c *gin.Context) {
	id, _ := c.Params.Get("id")

	post, _ := controllers.GetOnePostFromDB(id)

	c.JSON(200, gin.H{
		"post": post,
	})
}

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
	_, err := controllers.SavePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return response
	c.JSON(http.StatusOK, gin.H{
		"post": post,
		"tags": strings.Join(post.Tags, ", "),
	})

}

func GetAllPosts(c *gin.Context) {

	// get all posts from db
	posts, err := controllers.GetAllPostsFromDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, posts)

}
