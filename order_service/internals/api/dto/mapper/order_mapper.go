package mapper

import (
	"github.com/palashbhasme/order_service/internals/api/dto/request"
	"github.com/palashbhasme/order_service/internals/api/dto/response"
	"github.com/palashbhasme/order_service/internals/domain/models"
)

func ToOrderModel() *models.Order {
	var req request.OrderRequest

	order := &models.Order{
		UserID:   req.UserID,
		Quantity: req.Quantity,
		Status:   models.OrderPending, // Default status
	}

	// Preallocate slice memory for better performance
	order.OrderItems = make([]models.OrderItem, len(req.OrderItems))

	for i, item := range req.OrderItems {
		order.OrderItems[i] = ToItemModel(item)
	}

	return order
}

func ToItemModel(item request.OrderItemReq) models.OrderItem {
	return models.OrderItem{
		ProductID: item.ProductID,
		Price:     item.Price,
		Quantity:  item.Quantity,
	}
}

func ToItemResponse(item models.OrderItem) response.OrderItemResponse {
	return response.OrderItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		Price:     item.Price,
		Quantity:  item.Quantity,
	}
}

// ConvertOrderToResponse converts an Order model to an OrderResponse
func ToOrderResponse(order *models.Order) response.OrderResponse {
	orderItems := make([]response.OrderItemResponse, 0, len(order.OrderItems))

	for _, item := range order.OrderItems {
		orderItems = append(orderItems, ToItemResponse(item))
	}

	return response.OrderResponse{
		OrderID:    order.OrderID,
		UserID:     order.UserID,
		Quantity:   order.Quantity,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		OrderItems: orderItems,
	}
}
