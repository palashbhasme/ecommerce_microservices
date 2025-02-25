package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/palashbhasme/ecommerce_microservices/common/middlewares"
	"github.com/palashbhasme/ecommerce_microservices/common/models"

	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/api/dto/mapper"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/api/dto/request"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/api/dto/response"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/domain/repository"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	categoryRepo repository.CategoryRepository
	logger       *zap.Logger
}

func NewCategoryHandler(router *gin.Engine, repo repository.CategoryRepository, logger *zap.Logger) {
	categoryHandler := &CategoryHandler{
		categoryRepo: repo,
		logger:       logger,
	}
	authconfig := models.NewAuthConfig(os.Getenv("JWT_SECRET"))

	api := router.Group("/api")
	{
		categoryRoutes := api.Group("/category/v1")
		categoryRoutes.Use(middlewares.AuthMiddleware(*authconfig))
		{
			categoryRoutes.GET("/:id", categoryHandler.GetCategory)
			categoryRoutes.GET("/getall", categoryHandler.GetAllCategories)

			protectedRoutes := categoryRoutes.Group("/")
			protectedRoutes.Use(middlewares.AdminMiddleware())
			{
				protectedRoutes.PUT("/update/:id", categoryHandler.UpdateCategory)
				protectedRoutes.POST("/", categoryHandler.CreateCategory)
				protectedRoutes.DELETE("/delete/:id", categoryHandler.DeleteCategory)
			}

		}
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var categoryRequest request.CategoryRequest

	if err := c.ShouldBindJSON(&categoryRequest); err != nil {
		h.logger.Error("failed to bind request", zap.Error(err))
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	id := uuid.New().String()

	new_category := mapper.MapCategoryToRequest(&categoryRequest, id)

	err := h.categoryRepo.CreateCategory(&new_category)
	if err != nil {
		h.logger.Error("failed to create category", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusCreated,
		gin.H{"message": "category created successfully"})

}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := h.categoryRepo.GetCategoryByID(id)
	if err != nil {
		h.logger.Error("failed to get category", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to get category"})
		return
	}

	categoryResponse := mapper.MapCategoryToResponse(category)
	c.JSON(http.StatusOK,
		gin.H{"category": categoryResponse})
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	var categoryResponses []response.CategoryResponse

	categories, err := h.categoryRepo.GetAllCategories()
	if err != nil {
		h.logger.Error("failed to get categories", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to get categories"})
		return
	}

	for _, category := range categories {
		categoryResponses = append(categoryResponses, mapper.MapCategoryToResponse(&category))
	}
	c.JSON(http.StatusOK,
		gin.H{"categories": categoryResponses})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var categoryRequest request.CategoryRequest

	if err := c.ShouldBindJSON(&categoryRequest); err != nil {
		h.logger.Error("failed to bind request", zap.Error(err))
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	category := mapper.MapCategoryToRequest(&categoryRequest, id)

	err := h.categoryRepo.UpdateCategory(id, &category)
	if err != nil {
		h.logger.Error("failed to update category", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to update category"})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"message": "category updated successfully"})

}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	err := h.categoryRepo.DeleteCategory(id)
	if err != nil {
		h.logger.Error("failed to delete category", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to delete category"})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"message": "category deleted successfully"})
}

func (h *CategoryHandler) GetCategoryByName(c *gin.Context) {
	name := c.Param("name")
	category, err := h.categoryRepo.GetCategoryByName(name)
	if err != nil {
		h.logger.Error("failed to get category", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to get category"})
		return
	}

	categoryResponse := mapper.MapCategoryToResponse(category)
	c.JSON(http.StatusOK,
		gin.H{"category": categoryResponse})
}
