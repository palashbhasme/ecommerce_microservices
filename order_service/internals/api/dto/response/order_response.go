package response

import (
	"time"

	"github.com/palashbhasme/order_service/internals/domain/models"
)

// OrderResponse represents the structure of the order data sent to the user
type OrderResponse struct {
	OrderID    string              `json:"order_id"`
	UserID     string              `json:"user_id"`
	Quantity   int                 `json:"quantity"`
	Status     models.OrderStatus  `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
	OrderItems []OrderItemResponse `json:"order_items"`
}

// OrderItemResponse represents the structure of an order item in the order response
type OrderItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}
