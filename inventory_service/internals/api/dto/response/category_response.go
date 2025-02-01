package response

// CategoryResponse represents the response structure for category data
type CategoryResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Products    []ProductResponse `json:"products,omitempty"`
}
