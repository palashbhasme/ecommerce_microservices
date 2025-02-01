package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderShipped   OrderStatus = "shipped"
	OrderConfirmed OrderStatus = "confirmed"
	OrderCancelled OrderStatus = "cancelled"
	OrderDelivered OrderStatus = "delivered"
)

// takes in order status
func (s *OrderStatus) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("error invalid data for order status")
	}

	switch OrderStatus(v) {
	case OrderPending, OrderShipped, OrderConfirmed, OrderCancelled, OrderDelivered:
		*s = OrderStatus(v)
		return nil
	default:
		return errors.New("invalid order status value")
	}

}

// takes in order status and it typecasted to string
func (s OrderStatus) Value() (driver.Value, error) {
	return string(s), nil
}
