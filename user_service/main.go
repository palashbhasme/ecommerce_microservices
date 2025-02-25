package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/joho/godotenv"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/api"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/db"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/domain/repository"
	"github.com/palashbhasme/ecommerce_microservices/user_service/utils"
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

	collection := database.Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal("Could not create index:", err)
	}

	repo := repository.NewMongoRepository(database)

	api.Server(logger, repo)

}
