package rabbitmq

import (
	"encoding/json"

	"github.com/palashbhasme/order_service/internals/api/dto/request"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Struct for inventory check message
type InventoryRequest struct {
	OrderID string                 `json:"order_id"`
	Items   []request.OrderItemReq `json:"order_items"`
}

// PublishInventoryCheck sends an order inventory check request via RabbitMQ
func PublishInventoryCheck(orderID string, items []request.OrderItemReq, logger *zap.Logger) error {
	// Establish a single RabbitMQ connection (consider moving this to a shared package)
	conn, err := amqp.Dial("amqp://percy:secret@localhost:5672/backend")
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", zap.Error(err))
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel", zap.Error(err))
		return err
	}
	defer ch.Close()

	// Declare the exchange (should be "direct" if using a specific routing key)
	err = ch.ExchangeDeclare(
		"inventory_check", // Exchange name
		"direct",          // Type (use "fanout" if multiple consumers)
		true,              // Durable (ensures persistence)
		false,             // Auto-delete
		false,             // Internal
		false,             // No-wait
		nil,               // Arguments
	)
	if err != nil {
		logger.Error("Failed to declare exchange", zap.Error(err))
		return err
	}

	// Create the inventory request payload
	inventoryRequest := InventoryRequest{
		OrderID: orderID,
		Items:   items,
	}

	body, err := json.Marshal(inventoryRequest)
	if err != nil {
		logger.Error("Failed to marshal JSON", zap.Error(err))
		return err
	}

	// Publish message to the exchange
	err = ch.Publish(
		"inventory_check",     // Exchange name
		"inventory_check_key", // Routing key
		true,                  // Mandatory
		false,                 // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		logger.Error("Failed to publish message", zap.Error(err))
		return err
	}

	logger.Info("Successfully published inventory check",
		zap.String("exchange", "inventory_check"),
		zap.String("routing_key", "inventory_check_key"),
		zap.String("order_id", orderID),
	)

	return nil
}
