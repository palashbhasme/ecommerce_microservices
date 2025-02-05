package internals

import (
	"github.com/gin-gonic/gin"
	"github.com/palashbhasme/order_service/internals/api/handlers"
	"github.com/palashbhasme/order_service/internals/api/rabbitmq"
	"github.com/palashbhasme/order_service/internals/domain/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Server(logger *zap.Logger, db *gorm.DB) error {
	router := gin.Default()
	repo := repository.NewPostgresRepository(db)

	config, err := rabbitmq.InitializeRabbit()
	go func() {
		if err := rabbitmq.UpdateOrderConsumer(logger, repo, config.Conn); err != nil {
			logger.Error("error calling update order publisher", zap.Error(err))
		}
	}()

	if err != nil {
		logger.Error("Failed to connect to rabbit mq", zap.Error(err))
	}

	handlers.InitializeOrderHandler(router, repo, logger, config)
	err = router.Run(":8082")
	if err != nil {
		return err
	}
	return nil
}
