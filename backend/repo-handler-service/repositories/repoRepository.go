package repositories

import (
	"fmt"

	"github.com/onkarr19/haven/repo-handler-service/models"
	"gorm.io/gorm"
)

type RepoRepository interface {
	CreateRepo(repo *models.Repo) error
}

type repoRepository struct {
	db *gorm.DB
}

func NewRepoRepository(db *gorm.DB) RepoRepository {
	return &repoRepository{db: db}
}

func (p *repoRepository) CreateRepo(repo *models.Repo) error {
	fmt.Println("Creating repo...")
	return p.db.Create(repo).Error
}
