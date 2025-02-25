package request

type OrderRequest struct {
	UserID     string         `json:"user_id" binding:"required"`          // Must be a valid UUID
	OrderItems []OrderItemReq `json:"order_items" binding:"required,dive"` // Validate each OrderItemReq
	Quantity   int            `json:"quantity" binding:"required,gt=0"`    // Must be greater than 0
}

// OrderItemReq represents individual order items.
type OrderItemReq struct {
	ProductID string  `json:"product_id" binding:"required,uuid"` // Must be a valid UUID
	Price     float64 `json:"price" binding:"required,gt=0"`      // Must be greater than 0
	Quantity  int     `json:"quantity" binding:"required,gt=0"`   // Must be greater than 0
}
