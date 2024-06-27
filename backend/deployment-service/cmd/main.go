package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/onkarr19/haven/deployment-handler-service/handlers"
	"github.com/onkarr19/haven/deployment-handler-service/models"
	"github.com/onkarr19/haven/deployment-handler-service/repositories"
	"github.com/onkarr19/haven/deployment-handler-service/services"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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

	rdsConfig := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	rds := redis.NewClient(rdsConfig)
	defer rds.Close()

	sql := sqlite.Open("test.db")
	db, sqlDB := repositories.ConnectDatabase(sql, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}, &models.Deployment{})
	defer sqlDB.Close()

	deploymentRepository := repositories.NewDeploymentRepository(db)
	deploymentService := services.NewDeploymentService(deploymentRepository, rds)
	deploymentHandler := handlers.NewDeploymentHandler(deploymentService, logrus.New())

	r.POST("/project", deploymentHandler.CreateDeployment)
	r.GET("/project/:name", deploymentHandler.GetDeployment)

	r.Run("localhost:8080")
}
