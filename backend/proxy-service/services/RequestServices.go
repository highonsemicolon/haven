package services

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/proxy-service/repositories"
	"github.com/sirupsen/logrus"
)

type ProxyService interface {
	ProxyRequest(*gin.Context, string)
}

type proxyService struct {
	proxyRepo repositories.ProxyRepository
	logger    *logrus.Logger
}

func NewProxyService(s3repo repositories.ProxyRepository, logger *logrus.Logger) ProxyService {
	return &proxyService{proxyRepo: s3repo, logger: logger}
}

func (s *proxyService) ProxyRequest(c *gin.Context, path string) {
	hostname := c.Request.Host
	subdomain := strings.Split(hostname, ".")[0]

	key := subdomain + path
	resTo := s.proxyRepo.GetObjectURL(key)

	target, err := url.Parse(resTo)
	if err != nil {
		s.logger.Errorf("failed to parse URL: %v", err)
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
		s.logger.Errorf("proxy error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proxy error", "details": err.Error()})
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
