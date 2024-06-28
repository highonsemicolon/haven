package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/proxy-service/services"
	"github.com/sirupsen/logrus"
)

type ProxyHandler struct {
	requestService services.ProxyService
	logger         *logrus.Logger
}

func NewProxyHandler(requestService services.ProxyService, logger *logrus.Logger) *ProxyHandler {
	return &ProxyHandler{requestService: requestService, logger: logger}
}

func (h *ProxyHandler) HandleProxy(c *gin.Context) {

	path := c.Request.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	hostname := c.Request.Host
	subdomain := strings.Split(hostname, ".")[0]

	key := subdomain + path
	resTo := h.requestService.GetObjectURL(key)

	target, err := url.Parse(resTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		req.URL.RawQuery = target.RawQuery
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proxy error", "details": err.Error()})
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
