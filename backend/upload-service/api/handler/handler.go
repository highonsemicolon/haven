package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
)

func Home(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome to haven's upload service."})
}

func Upload(c *gin.Context) {
	var body struct {
		Url string `json:"q"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s url not found.", body.Url)})
		return
	}

	repo_id := "ndngdgn" // generate a unique id for the repo
	file_path := fmt.Sprintf("/tmp/repo/%s", repo_id)

	if _, err := git.PlainClone(file_path, false, &git.CloneOptions{
		URL: body.Url,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error cloning the repository %s: %s", body.Url, err.Error())})
		return
	}

	// upload to s3

	// delete the repo from /tmp/repo

	c.JSON(http.StatusOK, gin.H{"repo_id": repo_id, "message": "File uploaded successfully."})
}
