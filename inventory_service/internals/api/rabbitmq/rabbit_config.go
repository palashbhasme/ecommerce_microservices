package rabbitmq

import (
	"github.com/palashbhasme/ecommerce_microservices/common"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	User     string
	Password string
	Host     string
	Vhost    string
}
type RabbitMQConn struct {
	Conn *amqp.Connection
}

func (config *RabbitMQConfig) InitializeRabbit() (RabbitMQConn, error) {
	conn, err := common.ConnectRabbitMQ(config.User, config.Password, config.Host, config.Vhost)

	if err != nil {
		return RabbitMQConn{}, err
	}

	return RabbitMQConn{
		Conn: conn,
	}, err
}
