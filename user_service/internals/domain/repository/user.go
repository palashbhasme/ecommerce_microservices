package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := r.collection.Database().Client().Ping(ctx, nil)
	if err != nil {
		return errors.New("MongoDB client is disconnected")
	}
	_, err = r.collection.InsertOne(ctx, user)

	return err
}

func (r *MongoUserRepository) UpdateUser(id models.MyObjectID, user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateFields := bson.M{}
	if user.FirstName != "" {
		updateFields["first_name"] = user.FirstName
	}
	if user.LastName != "" {
		updateFields["last_name"] = user.LastName
	}
	if user.Email != "" {
		updateFields["email"] = user.Email
	}
	if user.Phone != "" {
		updateFields["phone"] = user.Phone
	}

	// Log the update fields
	log.Printf("Update fields: %+v", updateFields)

	// If no fields to update, return early
	if len(updateFields) == 0 {
		return errors.New("no valid field provided to update")
	}

	update := bson.M{"$set": updateFields}

	// Perform the update
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *MongoUserRepository) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *MongoUserRepository) GetUserById(id models.MyObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	if err != nil {
		return nil, errors.New("user for given id not found")
	}

	return &user, nil

}

func (r *MongoUserRepository) GetAllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []*models.User

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("error decoding user: %w", err)
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}

func (r *MongoUserRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User // Declare as a value, not a pointer
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // No user found is not an error
		}
		return nil, fmt.Errorf("database error: %w", err) // Handle other DB errors
	}

	return &user, nil // Return the address of the user struct
}
