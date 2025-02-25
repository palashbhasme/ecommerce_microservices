package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/palashbhasme/ecommerce_microservices/common"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type UpdateOrder struct {
	OrderID string `json:"order_id" binding:"required"`
	Status  string `json:"status" binding:"required"`
}

func UpdateOrderPublisher(orderId string, status string, logger *zap.Logger, conn *amqp.Connection) {
	// Open a channel
	client, err := common.NewRabbitMQClient(conn)
	if err != nil {
		logger.Error("Failed to get a client", zap.Error(err))
	}
	defer client.Close()

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

	order := UpdateOrder{
		OrderID: orderId,
		Status:  status,
	}

	body, err := json.Marshal(order)
	if err != nil {
		logger.Error("Failed to marshal JSON", zap.Error(err))
	}
	err = client.Send(context.TODO(), "order_update", "order_update_key", amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	if err != nil {
		logger.Error("Failed to publish message", zap.Error(err))
	}

	logger.Info("Successfully published update order request",
		zap.String("exchange", "update_order"),
		zap.String("routing_key", "update_order_key"),
		zap.String("order_id", order.OrderID),
	)
}
