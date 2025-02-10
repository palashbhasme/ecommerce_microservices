package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/palashbhasme/order_service/internals/api/dto/mapper"
	"github.com/palashbhasme/order_service/internals/api/dto/request"
	"github.com/palashbhasme/order_service/internals/api/rabbitmq"
	"github.com/palashbhasme/order_service/internals/domain/repository"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type OrderHandler struct {
	repo       repository.OrdersRepository
	logger     *zap.Logger
	connRabbit *amqp091.Connection
}

func InitializeOrderHandler(router *gin.Engine, repo repository.OrdersRepository, logger *zap.Logger, config rabbitmq.RabbitMQConn) {
	orderHandler := OrderHandler{
		repo:       repo,
		logger:     logger,
		connRabbit: config.Conn,
	}

	api := router.Group("/api")
	{
		orderRoutes := api.Group("/orders/v1")
		{
			orderRoutes.POST("/", orderHandler.CreateOrder)
			orderRoutes.GET("/:id", orderHandler.GetOrderByID)
		}

	}
}

// creates an order and publishes and inventory check event on rabbitmq inventory_check exchange
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var orderRequest request.OrderRequest
	err := c.ShouldBindJSON(&orderRequest)
	if err != nil {
		h.logger.Error("error binding request", zap.Error(err))
		c.JSON(400, gin.H{"message": "invalid request body"})
		return
	}

	// var items = make([]models.OrderItem, len(orderRequest.OrderItems))

	// for i, reqItem := range orderRequest.OrderItems {
	// 	items[i] = mapper.ToItemModel(reqItem)
	// }

	// order := models.Order{
	// 	UserID:     orderRequest.UserID,
	// 	Quantity:   orderRequest.Quantity,
	// 	OrderItems: items,
	// 	Status:     models.OrderPending,
	// }

	order := mapper.ToOrderModel(orderRequest)
	orderID, err := h.repo.CreateOrder(order)
	if err != nil {
		h.logger.Error("error creating order", zap.Error(err))
		c.JSON(400, gin.H{"message": "invalid request body"})
		return
	}

	err = rabbitmq.PublishInventoryCheck(orderID, orderRequest.OrderItems, h.logger, h.connRabbit)
	if err != nil {
		h.logger.Error("error publishing order", zap.Error(err))
		c.JSON(500, gin.H{"message": "internal server error"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message":  "Order request received, processing...",
		"order_id": orderID,
		"status":   "pending",
	})
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {

	id := c.Param("id")
	h.logger.Info("Fetching order by id", zap.String("id", id))

	order, err := h.repo.GetOrderByID(id)

	if err != nil {
		h.logger.Error("error failed to fetch order", zap.Error(err))
		c.JSON(500, gin.H{"error": "falied to fetch order"})
		return
	}

	orderResponse := mapper.ToOrderResponse(order)
	c.JSON(200, gin.H{"order": orderResponse})

}
