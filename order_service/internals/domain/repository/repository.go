package repository

import "github.com/palashbhasme/order_service/internals/domain/models"

type OrdersRepository interface {
	CreateOrder(order *models.Order) (string, error)
	GetOrderByID(id string) (*models.Order, error)
	UpdateOrder(id string, order *models.Order) error
	CreateOrderItem(orderItem *models.OrderItem) error
	UpdateOrderStatus(id string, updates map[string]interface{}) error
}
