package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/proxy-service/services"
)

type ProxyHandler struct {
	requestService services.ProxyService
}

func NewProxyHandler(requestService services.ProxyService) *ProxyHandler {
	return &ProxyHandler{requestService: requestService}
}

func (h *ProxyHandler) HandleProxy(c *gin.Context) {

	path := c.Request.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	h.requestService.ProxyRequest(c, path)
}
