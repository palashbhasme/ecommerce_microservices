package mapper

import (
	"time"

	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/api/dto/request"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/api/dto/response"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/domain/models"
)

// MapUserToResponse maps a User model to a UserResponse struct.
func MapUserToResponse(user models.User) response.UserResponse {
	addresses := make([]response.AddressResponse, len(user.Addresses))
	for i, addr := range user.Addresses {
		addresses[i] = response.AddressResponse{
			ID:        addr.ID,
			Line1:     addr.Line1,
			Line2:     addr.Line2,
			City:      addr.City,
			State:     addr.State,
			Country:   addr.Country,
			ZipCode:   addr.ZipCode,
			IsDefault: addr.IsDefault,
		}
	}

	return response.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		DOB:       user.DOB.Format(time.RFC3339), // Format DOB as ISO 8601
		Phone:     user.Phone,
		Addresses: addresses,
		Account: response.AccountResponse{
			ID:       user.Account.ID,
			Username: user.Account.Username,
			IsActive: user.Account.IsActive,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// MapUserToRequest maps a UserRequest struct to a User model.
func MapUserToRequest(user *request.UserRequest) models.User {

	return models.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		DOB:       user.DOB,
		Phone:     user.Phone,
		Addresses: mapAddresses(user.Addresses),
		Account:   mapAccount(user.Account),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}

func mapAddresses(reqAddresses []request.Address) []models.Address {
	var addresses []models.Address

	for _, address := range reqAddresses {
		addresses = append(addresses, models.Address{
			Line1:     address.Line1,
			Line2:     address.Line2,
			City:      address.City,
			State:     address.State,
			Country:   address.Country,
			ZipCode:   address.ZipCode,
			IsDefault: address.IsDefault,
		})
	}

	return addresses
}

func mapAccount(reqAccount request.Account) models.Account {
	return models.Account{
		Username:     reqAccount.Username,
		PasswordHash: reqAccount.Password,
		IsActive:     true,
	}
}
