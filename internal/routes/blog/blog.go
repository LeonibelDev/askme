package blog

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leonibeldev/askme/internal/controllers"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

// Read godoc
// @Summary      Get a blog post by ID
// @Description  Retrieve a single blog post using its unique identifier
// @Tags         Blog
// @Produce      json
// @Param        id   path      string  true  "Post ID"
// @Success      200  {object}  models.ResponseMessage{data=models.Post}
// @Failure      404  {object}  models.ResponseMessage
// @Router       /posts/{id} [get]
func Read(c *gin.Context) {
	id, _ := c.Params.Get("id")

	post, err := controllers.GetOnePostFromDB(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Success: false,
			Message: "Post not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, models.ResponseMessage{
		Success: true,
		Message: "Post retrieved successfully",
		Data:    post,
	})
}

// Write godoc
// @Summary      Create a new blog post
// @Description  Create a blog post (authentication required)
// @Tags         Blog
// @Accept       json
// @Produce      json
// @Param        post  body      models.Post  true  "Blog post payload"
// @Success      200   {object}  models.ResponseMessage{data=string}
// @Failure      400   {object}  models.ResponseMessage
// @Failure      500   {object}  models.ResponseMessage
// @Security     BearerAuth
// @Router       /posts/new [post]
func Write(c *gin.Context) {

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	// set post date
	post.Date = time.Now()

	// set post author
	author, _ := c.Get("username")
	post.Author = author.(string)

	// set fullname
	fullname, _ := c.Get("fullname")
	post.FullName = fullname.(string)

	// save post
	postId, err := controllers.SavePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Success: false,
			Message: "Error saving post",
			Error:   err.Error(),
		})
		return
	}

	// return response
	c.JSON(http.StatusOK, models.ResponseMessage{
		Success: true,
		Message: "Post created successfully",
		Data:    postId,
	})

}

// GetAllPosts godoc
// @Summary      Get all blog posts
// @Description  Retrieve paginated list of blog posts
// @Tags         Blog
// @Produce      json
// @Param        offset  query     string  true  "Pagination offset"
// @Success      200     {object}  models.ResponseMessage{data=[]models.Post}
// @Failure      400     {object}  models.ResponseMessage
// @Failure      404     {object}  models.ResponseMessage
// @Router       /posts [get]
func GetAllPosts(c *gin.Context) {
	// get page
	offset := c.Query("offset")
	// limit := c.Query("limit")

	if len(offset) == 0 {
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Success: false,
			Message: "Offset query parameter is required",
		})
		return
	}
	// get all posts from db
	posts, err := controllers.GetAllPostsFromDB(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Success: false,
			Message: "Error retrieving posts",
			Error:   err.Error(),
		})
		return
	}

	if len(posts) == 0 {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Success: false,
			Message: "No posts found",
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseMessage{
		Success: true,
		Message: "Posts retrieved successfully",
		Data:    posts,
	})

}

// GetPostsByTags godoc
// @Summary      Get posts by tag
// @Description  Retrieve blog posts filtered by tag
// @Tags         Blog
// @Produce      json
// @Param        tag     query     string  true  "Tag name"
// @Param        offset  query     string  false "Pagination offset"
// @Success      200     {object}  models.ResponseMessage{data=[]models.Post}
// @Failure      404     {object}  models.ResponseMessage
// @Router       /posts/tag [get]
func GetPostsByTags(c *gin.Context) {
	tag, _ := c.Params.Get("tag")
	offset := c.Query("offset")

	if len(offset) == 0 {
		offset = "0"
	}

	if tag == "" {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Success: false,
			Message: "Query 'tag' is required",
		})
		return
	}

	posts, err := controllers.GetPostsByTags(tag, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Success: false,
			Message: "Posts not found",
		})
		return
	}

	if len(posts) == 0 {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Success: false,
			Message: "No posts found for the given tag",
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseMessage{
		Success: true,
		Message: "Posts retrieved successfully",
		Data:    posts,
	})
}

// GetTopPosts godoc
// @Summary      Get top blog posts
// @Description  Retrieve top blog posts ordered by views
// @Tags         Blog
// @Produce      json
// @Success      200  {object}  models.ResponseMessage{data=[]models.Post}
// @Failure      500  {object}  models.ResponseMessage
// @Router       /posts/top [get]
func GetTopPosts(c *gin.Context) {
	posts, err := controllers.GetTopPosts(4)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Success: false,
			Message: "Error retrieving top posts",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseMessage{
		Success: true,
		Message: "Top posts retrieved successfully",
		Data:    posts,
	})
}

// GetPostsByAuthor godoc
// @Summary      Get posts by author
// @Description  Retrieve blog posts written by a specific author
// @Tags         Blog
// @Produce      json
// @Param        author  path      string  true  "Author username"
// @Success      200     {object}  models.ResponseMessage{data=[]models.Post}
// @Failure      404     {object}  models.ResponseMessage
// @Router       /posts/author/{author} [get]
func GetPostsByAuthor(c *gin.Context) {
	author, _ := c.Params.Get("author")
	offset, _ := c.GetQuery("offset")

	if len(offset) == 0 {
		offset = "0"
	}

	posts, err := controllers.GetPostsByAuthor(author, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Success: false,
			Message: "Posts not found",
			Error:   err.Error(),
		})
	}

	c.JSON(http.StatusOK, models.ResponseMessage{
		Success: true,
		Message: fmt.Sprintf("Posts by %s", author),
		Data:    posts,
	})
}

// Update godoc
// @Summary      Update post information
// @Description  Retrieve a single blog post using its unique identifier
// @Tags         Blog
// @Produce      json
// @Param        id   path      string  true  "Post ID"
// @Success      200  {object}  models.ResponseMessage{data=models.Post}
// @Failure      404  {object}  models.ResponseMessage
// @Router       /posts/{id} [put]
func UpdatePost(c *gin.Context) {
	dbPostId, _ := c.Params.Get("id")

	var post models.Post
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	fmt.Printf("post from user: %v\n", post.ID)
	fmt.Printf("post from db: %v\n", dbPostId)

	//author, _ := c.Get("username")

	_, err := controllers.UpdatePost(post, "leonibel_45771")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Success: false,
			Message: "Error updating the post",
		})
		return
	}

	c.JSON(http.StatusOK, models.ResponseMessage{
		Success: true,
		Message: "Post updated successfully",
		Data:    post,
	})

}
