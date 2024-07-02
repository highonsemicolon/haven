package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/onkarr19/haven/deployment-service/handlers"
	"github.com/onkarr19/haven/deployment-service/models"
	"github.com/onkarr19/haven/deployment-service/repositories"
	"github.com/onkarr19/haven/deployment-service/services"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	gormLogger "gorm.io/gorm/logger"

	"gorm.io/gorm"
)

var logger *logrus.Logger
var rds *redis.Client

func init() {
	logger = logrus.New()

	err := godotenv.Load()
	if err != nil {
		logger.Errorf("Error loading .env file: %v", err)
	}

	rdsConfig := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	rds = redis.NewClient(rdsConfig)

	_, err = rds.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf("Could not connect to Redis: %v", err)
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

	sql := sqlite.Open("test.db")
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
