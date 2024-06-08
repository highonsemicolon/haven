package repositories

import (
	"github.com/onkarr19/haven/deployment-handler-service/models"
	"gorm.io/gorm"
)

type DeploymentRepository interface {
	CreateDeployment(*models.Deployment) error
	GetDeploymentByName(string) (*models.Deployment, error)
}

type deploymentRepository struct {
	db *gorm.DB
}

func NewDeploymentRepository(db *gorm.DB) DeploymentRepository {
	return &deploymentRepository{db: db}
}

func (p *deploymentRepository) CreateDeployment(deployment *models.Deployment) error {
	return p.db.Create(deployment).Error
}

func (p *deploymentRepository) GetDeploymentByName(name string) (*models.Deployment, error) {
	deployment := &models.Deployment{}
	err := p.db.Where("name = ?", name).First(deployment).Error
	return deployment, err
}
