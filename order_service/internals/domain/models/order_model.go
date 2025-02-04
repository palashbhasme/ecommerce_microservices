package models

import (
	"time"

	"gorm.io/gorm"
)

// Orders Model
type Order struct {
	OrderID    string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` // UUID as Primary Key
	UserID     string      `gorm:"index"`                                          // Index for faster queries
	Quantity   int         `gorm:"not null"`
	Status     OrderStatus `gorm:"type:varchar(20);not null"`
	CreatedAt  time.Time   `gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;"` // Relationship with OrderItems
}

// OrderItems Model
type OrderItem struct {
	ID        string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID   string  `gorm:"not null;index"`
	ProductID string  `gorm:"not null"`
	Price     float64 `gorm:"not null"` // Changed from int to float for better precision
	Quantity  int     `gorm:"not null"`
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(Order{}, OrderItem{})
	return err
}
