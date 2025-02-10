package internals

import (
	"github.com/gin-gonic/gin"
	"github.com/palashbhasme/order_service/internals/api/handlers"
	"github.com/palashbhasme/order_service/internals/api/rabbitmq"
	"github.com/palashbhasme/order_service/internals/domain/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Server(logger *zap.Logger, db *gorm.DB, rabbitmqconfig *rabbitmq.RabbitMQConfig) error {

	repo := repository.NewPostgresRepository(db)
	conn, err := rabbitmqconfig.InitializeRabbit()
	if err != nil {
		logger.Error("Failed to connect to rabbit mq", zap.Error(err))
	}

	go func() {
		if err := rabbitmq.UpdateOrderConsumer(logger, repo, conn.Conn); err != nil {
			logger.Error("error calling update order publisher", zap.Error(err))
		}
	}()

	router := gin.Default()
	handlers.InitializeOrderHandler(router, repo, logger, conn)
	err = router.Run(":8082")
	if err != nil {
		return err
	}
	return nil
}
