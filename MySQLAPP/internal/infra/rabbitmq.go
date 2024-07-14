package infra

import (
	"fmt"
	"log"
	"os"

	"github.com/BernardoDenkvitts/MySQLApp/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", os.Getenv("user"), os.Getenv("RBMQpassword"), os.Getenv("host"), os.Getenv("port")))
	utils.FailOnError(err, "Error to connect to rabbitmq")

	channel, err := conn.Channel()
	utils.FailOnError(err, "Failed to open channel")

	err = channel.Confirm(false)
	utils.FailOnError(err, "Failed to put channel in confirmation mode")

	declareQueue(channel)
	log.Printf("%s Queue declared", os.Getenv("MySQLQueueName"))

	bindQueue(channel, os.Getenv("MongoDBExchange"))
	log.Printf("Queue binded to %s Exchange", os.Getenv("MongoDBExchange"))

	bindQueue(channel, os.Getenv("PostgresSQLExchange"))
	log.Printf("Queue binded to %s Exchange", os.Getenv("PostgresSQLExchange"))

	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.Connection.Close()
	r.Channel.Close()
}

func declareQueue(channel *amqp.Channel) {
	_, err := channel.QueueDeclare(
		os.Getenv("MySQLQueueName"),
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare queue")
}

func bindQueue(channel *amqp.Channel, exchange string) {
	err := channel.QueueBind(
		os.Getenv("MySQLQueueName"),
		"",
		exchange,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to bind queue "+os.Getenv("MySQLQueueName"))
}
