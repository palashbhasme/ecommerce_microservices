package rabbitmq

import (
	"github.com/palashbhasme/ecommerce_microservices/common"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	Conn *amqp.Connection
}

func InitializeRabbit() (RabbitMQConfig, error) {
	conn, err := common.ConnectRabbitMQ("percy", "secret", "localhost", "backend")

	if err != nil {
		return RabbitMQConfig{}, err
	}

	return RabbitMQConfig{
		Conn: conn,
	}, err
}
