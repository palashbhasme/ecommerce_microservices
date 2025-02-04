package api

import (
	"github.com/gin-gonic/gin"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/api/handlers"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/domain/repository"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/rabbitmq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RunServer(logger *zap.Logger, db *gorm.DB) error {
	router := gin.Default()
	// Start RabbitMQ consumer in a separate goroutine

	// Configure routes here if needed
	repo := repository.NewPostgresRepository(db)
	go func() {
		rabbitmq.ConsumeInventoryCheck(logger, *repo)
	}()

	handlers.NewCategoryHandler(router, repo, logger)
	handlers.NewProductHandler(router, repo, logger)

	err := router.Run(":8081")
	if err != nil {
		return err
	}
	return nil
}
