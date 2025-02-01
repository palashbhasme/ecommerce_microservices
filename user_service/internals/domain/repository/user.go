package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/palashbhasme/user_service/internals/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

var _ UserRepository = (*MongoUserRepository)(nil)

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

func (r *MongoUserRepository) UpdateUser(id primitive.ObjectID, user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Log user data
	log.Printf("User data: %+v", user)

	// Build the update fields
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
		return nil
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

func (r *MongoUserRepository) GetUserById(id primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user *models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	if err != nil {
		return nil, errors.New("user for given id not found")
	}

	return user, nil

}

func (r *MongoUserRepository) GetAllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []*models.User

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("users not found")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err = cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
