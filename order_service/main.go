package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/palashbhasme/ecommerce_microservices/common"
	internals "github.com/palashbhasme/order_service/internals/api"
	"github.com/palashbhasme/order_service/internals/api/rabbitmq"
	"github.com/palashbhasme/order_service/internals/domain/models"
	"github.com/palashbhasme/order_service/utils"
	"go.uber.org/zap"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	logger, err := utils.InitLogger()
	if err != nil {
		log.Fatal("Failed to initialze logger")
	}
	defer logger.Sync()

	logger.Info("Logger initialized successfully")

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

	orders_db, err := common.ConnectToDb(&config)
	if err != nil {
		log.Fatal("error connecting to database", zap.Error(err))

	}
	err = models.AutoMigrate(orders_db)
	if err != nil {
		log.Fatal("error migrating models", zap.Error(err))
	}
	err = internals.Server(logger, orders_db, &rabbitmqconfig)
	if err != nil {
		log.Fatal("error starting server", zap.Error(err))
	}
}
