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

	conn, err := rabbitmq.InitializeRabbit()
	if err != nil {
		logger.Error("Failed to connect to rabbit mq", zap.Error(err))
	}
	repo := repository.NewPostgresRepository(db)
	go func() {
		err = rabbitmq.InventoryCheckConsumer(logger, *repo, conn.Conn)
		logger.Error("consume inventory check stopped", zap.Error(err))
	}()

	router := gin.Default()
	handlers.NewCategoryHandler(router, repo, logger)
	handlers.NewProductHandler(router, repo, logger)

	err = router.Run(":8081")
	if err != nil {
		return err
	}
	return nil
}
