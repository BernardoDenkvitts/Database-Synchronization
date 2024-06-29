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
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", user, RBMQpassword, host, port))
	if err != nil {
		utils.FailOnError(err, "Error to connect to rabbitmq")
	}

	channel, err := conn.Channel()
	if err != nil {
		utils.FailOnError(err, "Failed to open channel")
	}

	if err := channel.Confirm(false); err != nil {
		utils.FailOnError(err, "Failed to put channel in confirmation mode")
	}
	declareRabbitMQExchange(channel)

	log.Println("RabbitMQ exchange declared")

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
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		utils.FailOnError(err, "Failed to declare exchange")
	}

}
