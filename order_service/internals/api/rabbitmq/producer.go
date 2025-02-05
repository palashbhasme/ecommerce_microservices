package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/palashbhasme/ecommerce_microservices/common"
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
func PublishInventoryCheck(orderID string, items []request.OrderItemReq, logger *zap.Logger, conn *amqp.Connection) error {

	client, err := common.NewRabbitMQClient(conn)
	if err != nil {
		logger.Error("Failed to get a client", zap.Error(err))
		return err
	}
	defer client.Close()

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

	err = client.Send(context.TODO(), "inventory_check", "inventory_check_key",
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

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
