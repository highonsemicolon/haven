package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/repo-handler-service/models"
	"github.com/onkarr19/haven/repo-handler-service/services"
)

type RepoHandler struct {
	repoService services.RepoService
}

func NewRepoHandler() *RepoHandler {
	return &RepoHandler{
		repoService: services.NewRepoService(),
	}
}

func (h *RepoHandler) CreateRepo(c *gin.Context) {
	repo := models.Repo{}
	if err := c.ShouldBindJSON(&repo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the repository name is unique
	if exists := h.repoService.IsUniqueRepo(repo.Name); exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Repository name already exists"})
		return
	}

	repo.ID = 789

	log.Printf("Creating project: %+v", repo)
	c.JSON(http.StatusCreated, repo)
}
