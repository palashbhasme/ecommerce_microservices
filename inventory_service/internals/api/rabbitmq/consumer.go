package rabbitmq

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/palashbhasme/ecommerce_microservices/common"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/domain/repository"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// InventoryRequest struct
type Inventory struct {
	OrderID string      `json:"order_id"`
	Items   []OrderItem `json:"order_items"`
}

// OrderItemReq struct
type OrderItem struct {
	ProductID string  `json:"product_id" binding:"required,uuid"`
	Price     float64 `json:"price" binding:"required,gt=0"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
}

func InventoryCheckConsumer(logger *zap.Logger, repo repository.PostgresRepository, conn *amqp.Connection) error {

	// Open a channel
	client, err := common.NewRabbitMQClient(conn)
	if err != nil {
		logger.Error("Failed to get a client", zap.Error(err))
		return err
	}

	// Declare the exchange (should be "direct" if using a specific routing key)
	err = client.CreateExchange("inventory_check", "direct", true, false, false, false)
	if err != nil {
		logger.Error("Failed to declare exchange", zap.Error(err))
		return err
	}

	// Declare queue
	err = client.CreateQueue("inventory_check", true, false)
	if err != nil {
		logger.Fatal("error declaring queue", zap.Error(err))
		return err
	}

	// Bind queue to exchange
	err = client.CreateBinding("inventory_check", "inventory_check_key", "inventory_check")
	if err != nil {
		logger.Error("error binding queue", zap.Error(err))
		return err
	}

	// Start consuming messages
	msgs, err := client.Consume("inventory_check", "inventory_check_consumer", false)
	if err != nil {
		logger.Error("failed to start consuming messages", zap.Error(err))
		return err
	}

	go func() {
		for msg := range msgs {
			var request Inventory
			if err := json.Unmarshal(msg.Body, &request); err != nil {
				logger.Error("failed to parse message", zap.Error(err))
				continue
			}

			// Extract variant IDs and quantities
			var variantIDs []string
			var quantities []int
			for _, item := range request.Items {
				variantIDs = append(variantIDs, item.ProductID)
				quantities = append(quantities, item.Quantity)
			}

			// Call DecrementStockLevel and check error
			available, totalPrice, err := repo.DecrementStockLevel(variantIDs, quantities)
			if err != nil {
				logger.Error("error updating stock levels", zap.Error(err))
				msg.Nack(false, true)
				continue
			}

			if available {
				logger.Info("Stock available", zap.String("OrderID", request.OrderID), zap.Float64("TotalPrice", totalPrice))
				msg.Ack(false)
				UpdateOrderPublisher(request.OrderID, logger, conn)
			} else {
				logger.Warn("Stock not available for one or more products", zap.String("OrderID", request.OrderID))
				msg.Nack(false, false)
			}
		}
	}()

	defer client.Close()
	// Handle SIGINT (Ctrl+C) and SIGTERM (Docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan // Wait for termination signal

	logger.Info("Shutting down consumer...")

	return nil
}
