package queue

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

// RabbitMQHandler define the RabbitMQ handler
type RabbitMQHandler struct {
	conn    *amqp.Connection
	queue   amqp.Queue
	channel *amqp.Channel
}

// NewRabbitMQHandler create new RabbitMQHandler
func NewRabbitMQHandler() *RabbitMQHandler {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	queue, err := channel.QueueDeclare(
		"notify",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &RabbitMQHandler{
		conn:    conn,
		queue:   queue,
		channel: channel,
	}
}

// Conn return the conn property
func (r RabbitMQHandler) Conn() *amqp.Connection {
	return r.conn
}

// Queue return the queue property
func (r RabbitMQHandler) Queue() amqp.Queue {
	return r.queue
}

// Channel return the channel property
func (r RabbitMQHandler) Channel() *amqp.Channel {
	return r.channel
}
