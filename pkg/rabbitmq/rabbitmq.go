package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url, exchange string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	rabbit := &RabbitMQ{
		conn:    conn,
		channel: ch,
	}

	// Ensure the queue exists
	_, err = rabbit.channel.QueueDeclare(
		"user_messages", // queue name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return rabbit, nil
}

// func (r *RabbitMQ) Publish(exchange, message string) error {
// 	return r.channel.Publish(exchange, "", false, false, amqp.Publishing{
// 		ContentType: "text/plain",
// 		Body:        []byte(message),
// 	})
// }

func (r *RabbitMQ) Publish(message string) error {
	err := r.channel.Publish(
		"",             // exchange (default)
		"user_messages", // routing key (queue name)
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Println("Message published successfully:", message)
	return nil
}

// Consume method to allow external access to the message consumer
func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
