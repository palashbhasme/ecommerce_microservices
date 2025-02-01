package api

import (
	"github.com/gin-gonic/gin"
	"github.com/palashbhasme/inventory_service/internals/api/handlers"
	"github.com/palashbhasme/inventory_service/internals/domain/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RunServer(logger *zap.Logger, db *gorm.DB) error {
	router := gin.Default()

	// Configure routes here if needed
	repo := repository.NewPostgresRepository(db)
	handlers.NewCategoryHandler(router, repo, logger)
	handlers.NewProductHandler(router, repo, logger)

	err := router.Run(":8081")
	if err != nil {
		return err
	}
	return nil
}
