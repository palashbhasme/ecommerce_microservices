package api

import (
	"github.com/gin-gonic/gin"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/domain/repository"
	"go.uber.org/zap"
)

func Server(log *zap.Logger, repo repository.UserRepository) {
	router := gin.Default()
	InitializeRoutes(router, log, repo)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}

}
