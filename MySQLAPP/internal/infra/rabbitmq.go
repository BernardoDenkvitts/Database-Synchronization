package infra

import (
	"fmt"
	"log"

	"github.com/BernardoDenkvitts/MySQLApp/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	user         = "guest"
	RBMQpassword = "guest"
	host         = "localhost"
	port         = "5672"
	exchangeName = "MYSQL"
	queueName    = "MYSQL"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", user, RBMQpassword, host, port))
	utils.FailOnError(err, "Error to connect to rabbitmq")

	channel, err := conn.Channel()
	utils.FailOnError(err, "Failed to open channel")
	err = channel.Confirm(false)
	utils.FailOnError(err, "Failed to put channel in confirmation mode")

	declareRabbitMQExchange(channel)
	log.Println("RabbitMQ exchange declared")

	declareQueue(channel)
	log.Println("Queue declared")

	bindQueue(channel)
	log.Println("Queue binded")

	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.Connection.Close()
	r.Channel.Close()
}

func declareRabbitMQExchange(channel *amqp.Channel) {
	err := channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare exchange")

}

func declareQueue(channel *amqp.Channel) {
	_, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare queue")
}

func bindQueue(channel *amqp.Channel) {
	err := channel.QueueBind(
		queueName,
		"",
		exchangeName, // TODO temporary name
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to bind queue")
}
