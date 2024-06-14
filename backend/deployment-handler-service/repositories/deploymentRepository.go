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

func (r *deploymentRepository) CreateDeployment(deployment *models.Deployment) error {
	deployment.HostedURL = "https://" + deployment.Name + ".haven.app"
	return r.db.Create(deployment).Error
}

func (r *deploymentRepository) GetDeploymentByName(name string) (*models.Deployment, error) {
	deployment := &models.Deployment{}
	err := r.db.Where("name = ?", name).First(deployment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return deployment, nil
}
