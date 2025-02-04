package repository

import (
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/domain/models"
)

type CategoryRepository interface {
	CreateCategory(category *models.Category) error
	GetCategoryByID(id string) (*models.Category, error)
	GetCategoryByName(name string) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	UpdateCategory(id string, category *models.Category) error
	DeleteCategory(id string) error
}

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id string) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
	GetProductsByCategoryID(categoryID string) ([]models.Product, error)
	GetProductsByCategoryName(categoryName string) ([]models.Product, error)
	UpdateProduct(id string, product *models.Product) error
	DeleteProduct(id string) error
	CheckStockLevel(variantID string, quantity int) (int, *float64, error)
	DecrementStockLevel(variantID []string, quantity []int) (bool, float64, error) //rabbit mq functions
}

type ProductVariantRepository interface {
	CreateProductVariant(variant *models.ProductVariant) error
	GetProductVariantByID(id string) (*models.ProductVariant, error)
	GetProductVariantsByProduct(productID string) ([]models.ProductVariant, error)
	UpdateProductVariant(id string, variant *models.ProductVariant) error
	DeleteProductVariant(id string) error
}
