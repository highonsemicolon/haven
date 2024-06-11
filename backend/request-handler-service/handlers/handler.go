package handlers

import (
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
	c.JSON(http.StatusOK, gin.H{"subdomain": subdomain})
}
