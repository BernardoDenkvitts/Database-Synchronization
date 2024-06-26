package message

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	user         = "guest"
	password     = "guest"
	host         = "localhost"
	port         = "5672"
	exchangeName = "MYSQL"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func OpenRabbitMQChannel() (*amqp.Channel, func() error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port))
	if err != nil {
		failOnError(err, "Error to connect to rabbitmq")
	}

	channel, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open channel")
	}

	return channel, channel.Close
}

func DeclareRabbitMQExchange(channel *amqp.Channel) error {
	err := channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		failOnError(err, "Failed to declare exchange")
	}
	return nil
}
