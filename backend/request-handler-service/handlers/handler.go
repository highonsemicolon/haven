package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/request-handler-service/services"
)

type RequestHandler struct {
	requestService services.RequestService
}

func NewRequestHandler(requestService services.RequestService) *RequestHandler {
	return &RequestHandler{requestService: requestService}
}

func (h *RequestHandler) GetDeployment(c *gin.Context) {
	host := c.Request.Host
	subdomain := h.requestService.GetSubdomain(host)

	content, contentType, err := h.requestService.GetDeploymentContent(context.Background(), subdomain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get deployment content"})
		return
	}
	defer content.Close()

	c.DataFromReader(http.StatusOK, -1, contentType, content, nil)
}
