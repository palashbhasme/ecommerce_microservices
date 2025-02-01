package mapper

import (
	"time"

	"github.com/palashbhasme/inventory_service/internals/api/dto/request"
	"github.com/palashbhasme/inventory_service/internals/api/dto/response"

	"github.com/palashbhasme/inventory_service/internals/domain/models"
)

// MapProductToResponse maps a Product model to the ProductResponse struct
func MapProductToResponse(product *models.Product) response.ProductResponse {
	variants := make([]response.ProductVariantResponse, len(product.Variants))
	for i, variant := range product.Variants {
		variants[i] = MapVariantToResponse(&variant)
	}

	categoryName := "General" // Default value if Category is nil
	if product.Category != nil {
		categoryName = product.Category.Name
	}

	return response.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		StockLevel:  product.StockQuantity,
		SKU:         product.SKU,
		Category:    categoryName,
		Variants:    variants,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}
}

// MapVariantToResponse maps a ProductVariant model to the ProductVariantResponse struct
func MapVariantToResponse(variant *models.ProductVariant) response.ProductVariantResponse {
	price := 0.0
	if variant.Price != nil {
		price = *variant.Price
	}

	return response.ProductVariantResponse{
		ID:         variant.ID,
		SKU:        variant.SKU,
		Size:       variant.Size,
		Color:      variant.Color,
		Price:      price,
		StockLevel: variant.StockQuantity,
	}
}

// MapProductsToResponse maps a slice of Product models to a slice of ProductResponse structs
func MapProductsToResponse(products []models.Product) []response.ProductResponse {
	productResponses := make([]response.ProductResponse, 0, len(products))
	for _, product := range products {
		productResponses = append(productResponses, MapProductToResponse(&product))
	}
	return productResponses
}

func MapProductToRequest(product *request.ProductRequest, productID string, variantIDs []string) models.Product {
	var categoryID *string
	if product.CategoryID != nil && *product.CategoryID != "" {
		categoryID = product.CategoryID
	}

	return models.Product{
		ID:          productID,
		Name:        product.Name,
		Description: product.Description,
		CategoryID:  categoryID,
		SKU:         product.SKU,
		Brand:       product.Brand,
		Variants:    MapProductVariantsToRequest(product.Variants, variantIDs),
	}
}

func MapProductVariantsToRequest(productvariants []request.ProductVariantRequest, variantIDs []string) []models.ProductVariant {
	var variants []models.ProductVariant
	for i, variant := range productvariants {
		variants = append(variants, MapProductVariantToRequest(&variant, variantIDs[i]))
	}
	return variants
}

func MapProductVariantToRequest(Productvariant *request.ProductVariantRequest, id string) models.ProductVariant {
	return models.ProductVariant{
		ID:            id,
		Color:         Productvariant.Color,
		Size:          Productvariant.Size,
		Price:         Productvariant.Price,
		StockQuantity: Productvariant.StockLevel,
		SKU:           Productvariant.SKU,
	}
}
