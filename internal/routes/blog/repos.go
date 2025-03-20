package blog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

func GetGitHubRepos(c *gin.Context) {
	username, exists := c.Params.Get("username")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username parameter is required"})
		return
	}

	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "GitHub API request failed"})
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var repos []models.GitHubRepo
	if err := json.Unmarshal(bodyBytes, &repos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response", "details": err.Error()})
		return
	}

	if len(repos) < 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "This user doesn't have any public repositories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"repos": repos,
	})
}
