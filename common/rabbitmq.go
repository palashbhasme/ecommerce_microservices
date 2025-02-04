package common

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// ConnectRabbitMQ will spawn a Connection
func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// NewRabbitMQClient will connect and return a Rabbitclient with an open connection
func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {

	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

// CreateQueue will create a new queue based on given cfgs
func (rc RabbitClient) CreateQueue(queueName string, durable, autodelete bool) error {
	_, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)
	return err
}

func (rc *RabbitClient) CreateExchange(name string, kind string, durable, autoDelete, internal, noWait bool) error {
	err := rc.ch.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, nil)
	if err != nil {
		return err
	}
	return nil
}

// CreateBinding is used to connect a queue to an Exchange using the binding rule
func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	// leaveing nowait false, having nowait set to false wctxill cause the channel to return an error and close if it cannot bind
	// the final argument is the extra headers, but we wont be doing that now
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

// Send publishes a message with context mandatory: true
func (rc RabbitClient) Send(ctx context.Context, exchange string, routingKey string, options amqp.Publishing) error {
	return rc.ch.PublishWithContext(ctx, exchange, routingKey, true, false, options)
}

// Consume is a wrapper around consume, it will return a Channel that can be used to digest messages
// Queue is the name of the queue to Consume
// Consumer is a unique identifier for the service instance that is consuming, can be used to cancel etc
// autoAck is important to understand, if set to true, it will automatically Acknowledge that processing is done
// This is good, but remember that if the Process fails before completion, then an ACK is already sent, making a message lost
// if not handled properly
func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}
