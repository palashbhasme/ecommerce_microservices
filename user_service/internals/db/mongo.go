package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(db_url string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(db_url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.New("error connecting to mongodb")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.New("failed to ping MongoDB: " + err.Error())
	}

	return client, nil
}
