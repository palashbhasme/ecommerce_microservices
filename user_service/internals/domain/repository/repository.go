package repository

import (
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/domain/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	UpdateUser(id models.MyObjectID, user *models.User) error
	DeleteUser(id string) error
	GetUserById(id models.MyObjectID) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}
