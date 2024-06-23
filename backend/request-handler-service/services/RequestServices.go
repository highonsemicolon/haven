package services

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/request-handler-service/repositories"
)

type ProxyService interface {
	ProxyRequest(*gin.Context, string)
}

type proxyService struct {
	proxyRepo repositories.ProxyRepository
}

func NewProxyService(s3repo repositories.ProxyRepository) ProxyService {
	return &proxyService{proxyRepo: s3repo}
}

func (s *proxyService) ProxyRequest(c *gin.Context, path string) {
	hostname := c.Request.Host
	subdomain := strings.Split(hostname, ".")[0]

	key := subdomain + path
	resTo := s.proxyRepo.GetObjectURL(key)

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
