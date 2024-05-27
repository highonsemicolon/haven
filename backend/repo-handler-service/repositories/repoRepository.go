package repositories

import (
	"github.com/onkarr19/haven/repo-handler-service/models"
)

type RepoRepository interface {
	CreateRepo(repo *models.Repo) error
	UpdateRepo(repo *models.Repo) error
}

type repoRepository struct {
}

func NewRepoRepository() RepoRepository {
	return &repoRepository{}
}

func (p *repoRepository) CreateRepo(repo *models.Repo) error {
	panic("unimplemented")
}

func (p *repoRepository) UpdateRepo(repo *models.Repo) error {
	panic("unimplemented")
}
