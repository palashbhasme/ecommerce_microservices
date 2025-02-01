package models

import (
	"time"
)

// Orders Model
type Orders struct {
	OrderID    string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` // UUID as Primary Key
	UserId     string       `gorm:"index"`                                          // Index for faster queries
	Quantity   int          `gorm:"not null"`
	Status     OrderStatus  `gorm:"type:varchar(20);not null"`
	CreatedAt  time.Time    `gorm:"autoCreateTime"`
	UpdatedAt  time.Time    `gorm:"autoUpdateTime"`
	OrderItems []OrderItems `gorm:"foreignKey:OrderId;constraint:OnDelete:CASCADE;"` // Relationship with OrderItems
}

// OrderItems Model
type OrderItems struct {
	ID        string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID   string  `gorm:"not null;index"`
	ProductID string  `gorm:"not null"`
	Price     float64 `gorm:"not null"` // Changed from int to float for better precision
	Quantity  int     `gorm:"not null"`
}
