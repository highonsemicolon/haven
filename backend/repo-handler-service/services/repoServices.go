package services

import (
	"github.com/onkarr19/haven/repo-handler-service/models"
	"github.com/onkarr19/haven/repo-handler-service/repositories"
)

type RepoService interface {
	CreateRepo(repo *models.Repo) error
	UpdateRepo(repo *models.Repo) error
	IsUniqueRepo(name string) bool
}

type repoService struct {
	repoRepo repositories.RepoRepository
}

func NewRepoService() RepoService {
	return &repoService{
		repoRepo: repositories.NewRepoRepository(),
	}
}

func (s *repoService) UpdateRepo(repo *models.Repo) error {
	panic("unimplemented")
}

func (s *repoService) CreateRepo(repo *models.Repo) error {
	panic("unimplemented")
}

func (s *repoService) IsUniqueRepo(_ string) bool {
	// TODO: check if the repository name already does not exists
	return true
}
