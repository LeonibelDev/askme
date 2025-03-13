package blog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

func GetGitHubRepos(c *gin.Context) {
	username, _ := c.Params.Get("username")

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

	var message strings.Builder
	message.WriteString(fmt.Sprintf("Github repos %s", username))

	c.JSON(http.StatusOK, gin.H{
		"repos": repos,
	})
}
