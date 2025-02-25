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

	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/domain/repository"

	"go.uber.org/zap"
)

type ProductHandler struct {
	repo   repository.ProductRepository
	logger *zap.Logger
}

func NewProductHandler(router *gin.Engine, repo repository.ProductRepository, logger *zap.Logger) {
	productHandler := &ProductHandler{
		repo:   repo,
		logger: logger,
	}
	authconfig := models.NewAuthConfig(os.Getenv("JWT_SECRET"))
	api := router.Group("/api")
	{
		productRoutes := api.Group("/products/v1")
		productRoutes.Use(middlewares.AuthMiddleware(*authconfig))
		{
			productRoutes.GET("/:id", productHandler.GetProduct)
			productRoutes.GET("/getall", productHandler.GetAllProducts)
			productRoutes.GET("/getbycategoryid/:categoryID", productHandler.GetProductsByCategoryID)
			productRoutes.GET("/getbycategoryname/:categoryName", productHandler.GetProductsByCategoryName)
			productRoutes.POST("/checkstock/:id", productHandler.CheckStockLevel)
			// productRoutes.POST("/updateStock/:id", productHandler.UpdateStockLevel) //no longer needed as it is handeld by rabbitmq

			protectedRoutes := productRoutes.Group("/")
			protectedRoutes.Use(middlewares.AdminMiddleware())
			{
				protectedRoutes.POST("/", productHandler.CreateProduct)
				protectedRoutes.DELETE("/delete/:id", productHandler.DeleteProduct)

			}
		}
	}
}

// create product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var productRequest request.ProductRequest
	if err := c.ShouldBindJSON(&productRequest); err != nil {
		h.logger.Error("failed to bind request", zap.Error(err))
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// Generate product ID and variant IDs
	productID := uuid.New().String()
	variantIDs := make([]string, len(productRequest.Variants))
	for i := range productRequest.Variants {
		variantIDs[i] = uuid.New().String()
	}

	// Map the request to the product model
	product := mapper.MapProductToRequest(&productRequest, productID, variantIDs)

	// Save the product to the database
	if err := h.repo.CreateProduct(&product); err != nil {
		h.logger.Error("failed to create product", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "product created successfully"})
}

// delete product by ID
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("Delting product with id", zap.String("id", id))
	err := h.repo.DeleteProduct(id)
	if err != nil {
		h.logger.Error("failed to delete product", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("Fetching product by id", zap.String("id", id))

	product, err := h.repo.GetProductByID(id)
	if err != nil {
		h.logger.Error("failed to get product", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to get product"})
		return
	}
	productResponse := mapper.MapProductToResponse(product)
	c.JSON(http.StatusOK, gin.H{"product": productResponse})
}

// Get All Products
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.repo.GetAllProducts()
	if err != nil {
		h.logger.Error("failed to get products", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to get products"})
		return
	}

	productResponses := mapper.MapProductsToResponse(products)
	c.JSON(http.StatusOK, gin.H{"products": productResponses})
}

// Get Products By Category ID
func (h *ProductHandler) GetProductsByCategoryID(c *gin.Context) {
	categoryID := c.Param("categoryID")
	h.logger.Info("Fetching product by category id", zap.String("id", categoryID))

	products, err := h.repo.GetProductsByCategoryID(categoryID)
	if err != nil {
		h.logger.Error("failed to get products by category", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to get products by category"})
		return
	}

	productResponses := mapper.MapProductsToResponse(products)
	c.JSON(http.StatusOK, gin.H{"products": productResponses})
}

// Get Products By category name
func (h *ProductHandler) GetProductsByCategoryName(c *gin.Context) {
	categoryName := c.Param("categoryName")
	h.logger.Info("Fetching product by category name", zap.String("category", categoryName))

	products, err := h.repo.GetProductsByCategoryName(categoryName)
	if err != nil {
		h.logger.Error("failed to get products by category", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to get products by category"})
		return
	}

	productResponses := mapper.MapProductsToResponse(products)
	c.JSON(http.StatusOK, gin.H{"products": productResponses})
}

func (h *ProductHandler) CheckStockLevel(c *gin.Context) {
	var quantityRequest request.UpdateQuantityRequest

	variantID := c.Param("id")

	if err := c.ShouldBindJSON(&quantityRequest); err != nil {
		h.logger.Error("failed to bind request body", zap.Error(err), zap.String("variant_id", variantID))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	quantity := quantityRequest.Quantity

	stockLevel, price, err := h.repo.CheckStockLevel(variantID, quantity)
	if err != nil {
		h.logger.Error("error checking stock leveles", zap.Error(err), zap.String("variant_id", variantID))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error checking stock level", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stock_level": stockLevel, "price": price})
}

// func (h *ProductHandler) UpdateStockLevel(c *gin.Context) {
// 	var quantityRequest request.UpdateQuantityRequest
// 	h.logger.Info("Updating stock level")

// 	variantID := c.Param("id")

// 	if err := c.ShouldBindJSON(&quantityRequest); err != nil {
// 		h.logger.Error("failed to bind request body", zap.Error(err), zap.String("variant_id", variantID))
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
// 		return
// 	}

// 	quantity := quantityRequest.Quantity

// 	newStockLevel, err := h.repo.DecrementStockLevel(variantID, quantity)
// 	if err != nil {
// 		h.logger.Error("error updating stock leveles", zap.Error(err), zap.String("variant_id", variantID))
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating stock leveles", "error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"remaining_stock": newStockLevel})

// }
