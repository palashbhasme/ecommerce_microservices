package response

// UserResponse represents the user data sent in API responses.
type UserResponse struct {
	ID        string            `json:"id"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Email     string            `json:"email"`
	DOB       string            `json:"dob"` // Formatted as ISO 8601 string
	Phone     string            `json:"phone"`
	Addresses []AddressResponse `json:"addresses"`  // Nested address responses
	Account   AccountResponse   `json:"account"`    // Nested account response
	CreatedAt string            `json:"created_at"` // ISO 8601 timestamp
	UpdatedAt string            `json:"updated_at"` // ISO 8601 timestamp
}

// AddressResponse represents the address data sent in API responses.
type AddressResponse struct {
	ID        string `json:"id"`
	Line1     string `json:"line1"`
	Line2     string `json:"line2,omitempty"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	ZipCode   string `json:"zip_code"`
	IsDefault bool   `json:"is_default"`
}

// AccountResponse represents the account data sent in API responses.
type AccountResponse struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
