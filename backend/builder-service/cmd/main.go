package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/onkarr19/haven/builder-service/handlers"
	"github.com/onkarr19/haven/builder-service/models"
	"github.com/onkarr19/haven/builder-service/repositories"
	"github.com/onkarr19/haven/builder-service/services"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var logger *logrus.Logger
var rds *redis.Client
var ctx = context.Background()

func init() {

	if err := godotenv.Load(); err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}

	logger = logrus.New()

	rdsConfig := &redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	rds = redis.NewClient(rdsConfig)

	pong, err := rds.Ping(ctx).Result()
	if err != nil {
		logger.Fatalf("Could not connect to Redis: %v", err)
	}
	logger.Println(pong)
}

func main() {

	inputQueue := os.Getenv("INPUT_QUEUE")
	outputQueue := os.Getenv("OUTPUT_QUEUE")

	sql := sqlite.Open("test.db")
	db, sqlDB, err := repositories.ConnectDatabase(sql, &gorm.Config{Logger: gormLogger.Default.LogMode(gormLogger.Silent)}, &models.Builder{})
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

	repo := repositories.NewBrokerRepository(db, rds, inputQueue, outputQueue)
	service := services.NewBrokerService(repo)
	handler := handlers.NewHandler(logger, service)

	if err := handler.Listen(ctx); err != nil {
		logger.Fatalf("Error listening: %v", err)
	}

}
