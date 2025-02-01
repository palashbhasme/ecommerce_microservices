package repository

import (
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	UpdateUser(id primitive.ObjectID, user *models.User) error
	DeleteUser(id string) error
	GetUserById(id primitive.ObjectID) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
}
