package rabbitmq

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/palashbhasme/ecommerce_microservices/common"
	"github.com/palashbhasme/order_service/internals/domain/repository"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type UpdateOrder struct {
	OrderID string `json:"order_id" binding:"required"`
	Status  string `json:"status" binding:"required"`
}

func UpdateOrderConsumer(logger *zap.Logger, repo repository.OrdersRepository, conn *amqp.Connection) error {
	// Open a channel
	client, err := common.NewRabbitMQClient(conn)
	if err != nil {
		logger.Error("Failed to get a client", zap.Error(err))
	}

	// Declare the exchange (should be "direct" if using a specific routing key)
	err = client.CreateExchange("order_update", "direct", true, false, false, false)
	if err != nil {
		logger.Error("Failed to declare exchange", zap.Error(err))
	}
	// Declare queue
	err = client.CreateQueue("order_update", true, false)
	if err != nil {
		logger.Error("error declaring queue", zap.Error(err))
	}

	// Bind queue to exchange
	err = client.CreateBinding("order_update", "order_update_key", "order_update")
	if err != nil {
		logger.Error("error binding queue", zap.Error(err))
	}

	msgs, err := client.Consume("order_update", "order_update_consumer", false)
	if err != nil {
		logger.Error("failed to start consuming messages", zap.Error(err))
		return err
	}

	for msg := range msgs {
		var order UpdateOrder
		if err := json.Unmarshal(msg.Body, &order); err != nil {
			logger.Error("failed to parse message", zap.Error(err))
			msg.Nack(false, false)
			continue
		}

		// Initialize the map
		status := map[string]interface{}{
			"status": order.Status,
		}

		// Call update order and check error
		err := repo.UpdateOrderStatus(order.OrderID, status)
		if err != nil {
			logger.Error("error updating order status", zap.Error(err))
			msg.Nack(false, false)
			continue
		} else {
			msg.Ack(false)
		}
	}
	defer client.Close()

	// Handle SIGINT (Ctrl+C) and SIGTERM (Docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // Wait for termination signal

	logger.Info("Shutting down consumer...")

	return nil
}
