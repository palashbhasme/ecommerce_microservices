package rabbitmq

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/domain/repository"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// InventoryRequest struct
type InventoryRequest struct {
	OrderID string         `json:"order_id"`
	Items   []OrderItemReq `json:"order_items"`
}

// OrderItemReq struct
type OrderItemReq struct {
	ProductID string  `json:"product_id" binding:"required,uuid"`
	Price     float64 `json:"price" binding:"required,gt=0"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
}

func ConsumeInventoryCheck(logger *zap.Logger, repo repository.PostgresRepository) {
	// Establish connection
	conn, err := amqp.Dial("amqp://percy:secret@localhost:5672/backend")
	if err != nil {
		logger.Fatal("failed to establish a connection to RabbitMQ", zap.Error(err))
		return
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("failed to open channel", zap.Error(err))
		return
	}
	defer ch.Close()

	// Declare exchange
	err = ch.ExchangeDeclare(
		"inventory_check", // Exchange name
		"direct",          // Type (use "fanout" if multiple consumers)
		true,              // Durable
		false,             // Auto-deleted
		false,             // Internal
		false,             // No-wait
		nil,               // Arguments
	)
	if err != nil {
		logger.Fatal("error declaring exchange", zap.Error(err))
		return
	}

	// Declare queue
	q, err := ch.QueueDeclare("inventory_check", false, false, false, false, nil)
	if err != nil {
		logger.Fatal("error declaring queue", zap.Error(err))
		return
	}

	// Bind queue to exchange
	err = ch.QueueBind(q.Name, "inventory_check_key", "inventory_check", false, nil)
	if err != nil {
		logger.Fatal("error binding queue", zap.Error(err))
		return
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // Auto-acknowledge (set to false if you want manual acknowledgment)
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,
	)
	if err != nil {
		logger.Fatal("failed to start consuming messages", zap.Error(err))
		return
	}

	go func() {
		for msg := range msgs {
			var request InventoryRequest
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
				continue
			}

			if available {
				logger.Info("Stock available", zap.String("OrderID", request.OrderID), zap.Float64("TotalPrice", totalPrice))
			} else {
				logger.Warn("Stock not available for one or more products", zap.String("OrderID", request.OrderID))
			}
		}
	}()

	// Handle SIGINT (Ctrl+C) and SIGTERM (Docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // Wait for termination signal

	logger.Info("Shutting down consumer...")
}
