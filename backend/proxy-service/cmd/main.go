package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/highonsemicolon/haven/proxy-service/handlers"
	"github.com/highonsemicolon/haven/proxy-service/repositories"
	"github.com/highonsemicolon/haven/proxy-service/services"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	defer logger.Writer().Close()
}

func ErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		c.JSON(-1, c.Errors.Last())
	}
}

func main() {
	r := gin.Default()
	r.Use(ErrorHandler)

	base_path := os.Getenv("BASE_PATH")

	s3Repo := repositories.NewProxyRepository(base_path)
	requestService := services.NewProxyService(s3Repo)
	requestHandler := handlers.NewProxyHandler(requestService, logger)

	r.NoRoute(requestHandler.HandleProxy)

	r.Run("localhost:8080")
}
