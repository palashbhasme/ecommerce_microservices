package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/palashbhasme/inventory_service/internals/api"
	"github.com/palashbhasme/inventory_service/internals/db"
	"github.com/palashbhasme/inventory_service/internals/domain/models"
	"github.com/palashbhasme/inventory_service/utils"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	logger, err := utils.InitLogger()
	if err != nil {
		fmt.Println("Error while initializing logger")
	}
	defer logger.Sync()
	logger.Info("Logger initialized successfully")

	config := db.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}
	inventory_db, err := db.ConnectToDb(&config)
	if err != nil {
		log.Fatal("error connecting to database", zap.Error(err))

	}
	err = models.AutoMigrate(inventory_db)
	if err != nil {
		log.Fatal("error migrating models", zap.Error(err))
	}

	err = api.RunServer(logger, inventory_db)
	if err != nil {
		log.Fatal("error running server", zap.Error(err))
	}

}
