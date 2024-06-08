package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/onkarr19/haven/deployment-handler-service/models"
	"github.com/onkarr19/haven/deployment-handler-service/services"
)

type DeploymentHandler struct {
	deploymentService services.DeploymentService
}

func NewDeploymentHandler(deploymentService services.DeploymentService) *DeploymentHandler {
	return &DeploymentHandler{deploymentService: deploymentService}
}

func (h *DeploymentHandler) CreateDeployment(c *gin.Context) {
	var deployment models.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	deployment.ID = uuid.New()

	if err := h.deploymentService.CreateDeployment(&deployment); err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusCreated, deployment)
}
