package repository

import (
	"github.com/palashbhasme/order_service/internals/domain/models"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateOrder(order *models.Order) (string, error) {
	result := r.db.Create(order)
	if result.Error != nil {
		return "", result.Error
	}
	return order.OrderID, nil
}

func (r *PostgresRepository) GetOrderByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderItems").Model(models.Order{}).First(&order, "order_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *PostgresRepository) UpdateOrder(id string, order *models.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Order{}).Where("order_id = ?", id).Updates(order).Error
		if err != nil {
			return err
		}

		for _, orderItems := range order.OrderItems {
			if err := tx.Save(&orderItems).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *PostgresRepository) CreateOrderItem(orderItem *models.OrderItem) error {
	return r.db.Create(&orderItem).Error
}

func (r *PostgresRepository) UpdateOrderStatus(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.Order{}).Where("order_id = ?", id).Updates(updates).Error
}
