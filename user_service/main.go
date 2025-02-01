package main

import (
	"context"
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/joho/godotenv"
	"github.com/palashbhasme/user_service/internals/api"
	"github.com/palashbhasme/user_service/internals/db"
	"github.com/palashbhasme/user_service/internals/domain/repository"
	"github.com/palashbhasme/user_service/utils"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	logger, err := utils.InitLogger()

	if err != nil {
		log.Fatalln("Failed to initialize logger")
	}

	defer logger.Sync()

	db_url := os.Getenv("DB_URL")

	client, err := db.ConnectMongo(db_url)
	if err != nil {
		log.Fatal("Error starting server", zap.Error(err))
	}

	// Defer the disconnect when the application exits
	defer client.Disconnect(context.Background())

	database := client.Database("users_database")

	repo := repository.NewMongoRepository(database)

	api.Server(logger, repo)

}
