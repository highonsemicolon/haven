package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/onkarr19/haven/deployment-handler-service/handlers"
	"github.com/onkarr19/haven/deployment-handler-service/models"
	"github.com/onkarr19/haven/deployment-handler-service/repositories"
	"github.com/onkarr19/haven/deployment-handler-service/services"
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

func ConnectDatabase(sql gorm.Dialector, config *gorm.Config, models *models.Deployment) *gorm.DB {
	db, err := gorm.Open(sql, config)
	if err != nil {
		log.Fatalf("failed to connect database: %+v", err)
	}
	if err := db.AutoMigrate(models); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database connected")
	return db
}

func main() {
	r := gin.Default()
	r.Use(ErrorHandler)

	sql := sqlite.Open("test.db")
	db := ConnectDatabase(sql, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}, &models.Deployment{})

	deploymentRepository := repositories.NewDeploymentRepository(db)
	deploymentService := services.NewDeploymentService(deploymentRepository)
	deploymentHandler := handlers.NewDeploymentHandler(deploymentService, logrus.New())

	r.POST("/project", deploymentHandler.CreateDeployment)
	r.GET("/project/:name", deploymentHandler.GetDeployment)

	r.Run("localhost:8080")
}
