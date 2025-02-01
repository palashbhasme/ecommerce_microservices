package repository

import (
	"fmt"

	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/domain/models"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *PostgresRepository) GetCategoryByID(id string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *PostgresRepository) GetCategoryByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *PostgresRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *PostgresRepository) UpdateCategory(id string, category *models.Category) error {
	result := r.db.Model(&models.Category{}).Where("id = ?", id).Updates(category)
	if result.RowsAffected == 0 {
		return fmt.Errorf("error updating provided category %v", id)
	}
	return nil
}

func (r *PostgresRepository) DeleteCategory(id string) error {
	result := r.db.Model(&models.Category{}).Where("id = ?", id).Delete(&models.Category{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("error deleting provided category %v", id)
	}
	return nil
}

func (r *PostgresRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *PostgresRepository) GetProductByID(id string) (*models.Product, error) {
	var product models.Product

	err := r.db.Preload("Variants").Preload("Category").First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *PostgresRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Variants").Preload("Category").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresRepository) GetProductsByCategoryID(categoryID string) ([]models.Product, error) {
	var products []models.Product
	// Preload the associated Variants and Category
	err := r.db.Preload("Variants").Preload("Category").
		Where("category_id = ?", categoryID).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresRepository) GetProductsByCategoryName(categoryName string) ([]models.Product, error) {
	var products []models.Product

	err := r.db.Preload("Variants").Joins("JOIN categories ON products.category_id = categories.id").
		Where("categories.name = ?", categoryName).
		Find(&products).Error

	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresRepository) UpdateProduct(id string, product *models.Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update the product
		if err := tx.Model(&models.Product{}).Where("id = ?", id).Updates(product).Error; err != nil {
			return err
		}

		// Update or create variants
		for _, variant := range product.Variants {
			if err := tx.Save(&variant).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *PostgresRepository) DeleteProduct(id string) error {
	// Perform the delete operation
	result := r.db.Where("id = ?", id).Delete(&models.Product{})

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("product with id %s not found", id)
	}

	// If there was an error during the delete operation
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PostgresRepository) CreateProductVariant(variant *models.ProductVariant) error {
	return r.db.Create(variant).Error
}

func (r *PostgresRepository) GetProductVariantByID(id string) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.First(&variant, id).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *PostgresRepository) GetProductVariantsByProduct(productID string) ([]models.ProductVariant, error) {
	var variants []models.ProductVariant
	err := r.db.Where("product_id = ?", productID).Find(&variants).Error
	if err != nil {
		return nil, err
	}
	return variants, nil
}

func (r *PostgresRepository) UpdateProductVariant(id string, variant *models.ProductVariant) error {
	err := r.db.Model(&models.ProductVariant{}).Where("id = ?", id).Updates(variant).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) DeleteProductVariant(id string) error {
	err := r.db.Delete(&models.ProductVariant{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) CheckStockLevel(variantID string, quantity int) (int, error) {
	var variant models.ProductVariant

	// Query the product variant by its ID
	err := r.db.Where("id = ?", variantID).First(&variant).Error
	if err != nil {
		// Return 0 and the error if no variant is found
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("variant not found")
		}
		return 0, err
	}

	// Check if the available stock is sufficient
	if variant.StockQuantity < quantity {
		return variant.StockQuantity, fmt.Errorf("insufficient stock, available: %d", variant.StockQuantity)
	}

	// Return the available stock level
	return variant.StockQuantity, nil
}

func (r *PostgresRepository) DecrementStockLevel(variantID string, quantity int) (int, error) {

	tx := r.db.Begin()

	var variant models.ProductVariant

	result := tx.Where("id = ?", variantID).First(&variant)

	if result.Error != nil {
		tx.Rollback()
		if result.Error == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("variant not found")
		}
		return 0, result.Error
	}

	if variant.StockQuantity < quantity {
		// Rollback if there's insufficient stock
		tx.Rollback()
		return variant.StockQuantity, fmt.Errorf("insufficient stock, available: %d", variant.StockQuantity)
	}
	variant.StockQuantity -= quantity
	if err := tx.Save(&variant).Error; err != nil {
		// Rollback if the update fails
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return variant.StockQuantity, nil
}
