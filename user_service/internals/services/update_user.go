package services

import (
	"errors"

	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/api/dto/request"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/domain/models"
)

func UpdateUser(existingUser *models.User, updates request.UpdateUserRequest) (*models.User, error) {
	if existingUser == nil {
		return nil, errors.New("existing user not found")
	}
	updatedUser := existingUser

	if updates.FirstName != "" {
		updatedUser.FirstName = updates.FirstName
	}
	if updates.LastName != "" {
		updatedUser.LastName = updates.LastName
	}
	if updates.Email != "" {
		updatedUser.Email = updates.Email
	}
	if updates.Phone != "" {
		updatedUser.Phone = updates.Phone
	}

	return existingUser, nil
}
