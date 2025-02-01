package request

import "time"

type CategoryRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
