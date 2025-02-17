package request

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/api/dto/validators"
)

type UserRequest struct {
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	DOB       time.Time `json:"dob" validate:"required"`
	Phone     string    `json:"phone" validate:"required,len=10"`
	Addresses []Address `json:"addresses"`
	Account   Account   `json:"account"`
}

type Address struct {
	Line1     string `json:"line1" validate:"required"`
	Line2     string `json:"line2,omitempty"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	Country   string `json:"country" validate:"required"`
	ZipCode   string `json:"zip_code" validate:"required"`
	IsDefault bool   `json:"is_default"`
}

type Account struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Email     string    `json:"email,omitempty" validate:"omitempty,email"`
	DOB       time.Time `json:"dob,omitempty"`
	Phone     string    `json:"phone,omitempty" validate:"omitempty,len=10"`
}

// NewValidator creates a new validator instance and registers custom validations.
func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("customEmail", validators.ValidateEmail)
	return v
}
