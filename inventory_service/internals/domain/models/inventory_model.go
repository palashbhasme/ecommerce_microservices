package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Products    []Product `gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
type Product struct {
	ID            string           `gorm:"type:uuid;primaryKey"`
	Name          string           `gorm:"type:varchar(255);not null"`
	Description   string           `gorm:"type:text"`
	CategoryID    *string          `gorm:"type:uuid"` // Nullable foreign key
	Category      *Category        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SKU           string           `gorm:"type:varchar(100);unique;not null"`
	StockQuantity int              `gorm:"default:0"`
	Brand         string           `gorm:"type:varchar(255)"`
	CreatedAt     time.Time        `gorm:"autoCreateTime"`
	UpdatedAt     time.Time        `gorm:"autoUpdateTime"`
	Variants      []ProductVariant `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"` // Cascade delete for variants
}

type ProductVariant struct {
	ID            string    `gorm:"type:uuid;primaryKey"`
	ProductID     string    `gorm:"type:uuid;not null"`
	Color         string    `gorm:"type:varchar(100)"`
	Size          string    `gorm:"type:varchar(50)"`
	Price         *float64  `gorm:"type:decimal(10,2);not null"`
	StockQuantity int       `gorm:"default:0"`
	SKU           string    `gorm:"type:varchar(100);unique;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Category{}, &Product{}, &ProductVariant{})
	return err
}
