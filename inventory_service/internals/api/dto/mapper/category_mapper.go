package mapper

import (
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/api/dto/request"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/api/dto/response"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/domain/models"
)

// MapCategoryToResponse maps a Category model to the CategoryResponse struct
func MapCategoryToResponse(category *models.Category) response.CategoryResponse {
	return response.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Products:    MapProductsToResponse(category.Products),
	}
}

func MapCategoryToRequest(category *request.CategoryRequest, id string) models.Category {
	return models.Category{
		ID:          id,
		Name:        category.Name,
		Description: category.Description,
	}
}
