package internals

import (
	"github.com/gin-gonic/gin"
	"github.com/palashbhasme/order_service/internals/api/handlers"
	"github.com/palashbhasme/order_service/internals/domain/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Server(logger *zap.Logger, db *gorm.DB) error {
	router := gin.Default()
	repo := repository.NewPostgresRepository(db)
	handlers.InitializeOrderHandler(router, repo, logger)

	err := router.Run(":8082")
	if err != nil {
		return err
	}
	return nil
}
