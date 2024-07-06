package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/highonsemicolon/haven/deployment-service/handlers"
	"github.com/highonsemicolon/haven/deployment-service/models"
	"github.com/highonsemicolon/haven/deployment-service/repositories"
	"github.com/highonsemicolon/haven/deployment-service/services"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var logger *logrus.Logger
var rds *redis.Client

func init() {
	logger = logrus.New()

	err := godotenv.Load()
	if err != nil {
		logger.Errorf("Error loading .env file: %v", err)
	}

	addr, err := redis.ParseURL(os.Getenv("REDIS_URI"))
	if err != nil {
		logger.Panicf("Could not parse Redis URI: %v", err)
	}

	rds = redis.NewClient(addr)

	_, err = rds.Ping(context.Background()).Result()
	if err != nil {
		logger.Panicf("Could not connect to Redis: %v", err)
	}
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

	dsn := os.Getenv("DSN_URI")
	sql := sqlite.Open(dsn)
	db, sqlDB, err := repositories.ConnectDatabase(sql, &gorm.Config{Logger: gormLogger.Default.LogMode(gormLogger.Silent)}, &models.Deployment{})
	if err != nil {
		logger.Fatalf("failed to connect database: %+v", err)
	} else {
		logger.Info("Database connected")
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Errorf("Error closing database connection: %v", err)
		}
		if err := rds.Close(); err != nil {
			logger.Errorf("Error closing Redis client: %v", err)
		}
		defer logger.Writer().Close()
	}()

	deploymentRepository := repositories.NewDeploymentRepository(db, rds)
	deploymentService := services.NewDeploymentService(deploymentRepository)
	deploymentHandler := handlers.NewDeploymentHandler(deploymentService, logger)

	r.POST("/project", deploymentHandler.CreateDeployment)
	r.GET("/project/:name", deploymentHandler.GetDeployment)

	r.GET("/ws/logs/:id", deploymentHandler.HandleWebSocket)

	r.Run("localhost:8080")
}
