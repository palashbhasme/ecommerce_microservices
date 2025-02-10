package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/palashbhasme/ecommerce_microservices/common"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/api"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/api/rabbitmq"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/domain/models"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/utils"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize Logger
	logger, err := utils.InitLogger()
	if err != nil {
		fmt.Println("Error while initializing logger")
	}
	defer logger.Sync()
	logger.Info("Logger initialized successfully")

	// Connect to PostgreSQL
	config := common.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}

	rabbitmqconfig := rabbitmq.RabbitMQConfig{
		Host:     os.Getenv("RABBITMQ_HOST"),
		User:     os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASS"),
		Vhost:    os.Getenv("RABBITMQ_VHOST"),
	}

	inventory_db, err := common.ConnectToDb(&config)
	if err != nil {
		log.Fatal("Error connecting to database", zap.Error(err))
	}

	// Run database migrations
	err = models.AutoMigrate(inventory_db)
	if err != nil {
		log.Fatal("Error migrating models", zap.Error(err))
	}

	// Start the API server
	err = api.RunServer(logger, inventory_db, &rabbitmqconfig)
	if err != nil {
		log.Fatal("Error running server", zap.Error(err))
	}
}
