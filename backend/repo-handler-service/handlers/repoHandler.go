package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/onkarr19/haven/repo-handler-service/models"
	"github.com/onkarr19/haven/repo-handler-service/services"
)

type RepoHandler struct {
	repoService services.RepoService
}

func NewRepoHandler(repoService services.RepoService) *RepoHandler {
	return &RepoHandler{repoService: repoService}
}

func (h *RepoHandler) CreateRepo(c *gin.Context) {
	var repo models.Repo
	if err := c.ShouldBindJSON(&repo); err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	repo.ID = uuid.New()

	if err := h.repoService.CreateRepo(&repo); err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusCreated, repo)
}
