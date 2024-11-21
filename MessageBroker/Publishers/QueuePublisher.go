package publishers

import (
	brokerconstants "NetworkAnalyzer/MessageBroker/Constants"
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueuePublisher struct {
}

func (publisher *QueuePublisher) OpenConnection() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ Queues")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel to RabbitMQ Queues")
	defer channel.Close()

	messageQueue, err := channel.QueueDeclare(
		"NetworkMessage",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	body := "Hello World"
	err = channel.PublishWithContext(context,
		"",
		messageQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: brokerconstants.PlainText,
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message to RabbitMQ Queue")
	log.Printf(" [x] Sent %s\n", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
