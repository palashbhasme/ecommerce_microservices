package request

import "time"

type UserRequest struct {
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	DOB       time.Time `json:"dob" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	Addresses []Address `json:"addresses"`
	Account   Account   `json:"account"`
}

type Address struct {
	Line1     string `json:"line1" binding:"required"`
	Line2     string `json:"line2,omitempty"`
	City      string `json:"city" binding:"required"`
	State     string `json:"state" binding:"required"`
	Country   string `json:"country" binding:"required"`
	ZipCode   string `json:"zip_code" binding:"required"`
	IsDefault bool   `json:"is_default"`
}

type Account struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
