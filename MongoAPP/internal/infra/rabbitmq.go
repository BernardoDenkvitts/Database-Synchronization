package infra

import (
	"fmt"
	"log"
	"os"

	"github.com/BernardoDenkvitts/MongoAPP/internal/utils"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {
	path, _ := os.Getwd()
	err := godotenv.Load(path + "/../.env")
	utils.FailOnError(err, "Failed to load env file")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", os.Getenv("user"), os.Getenv("RBMQpassword"), os.Getenv("host"), os.Getenv("port")))
	utils.FailOnError(err, "Error to connect to rabbitmq")

	channel, err := conn.Channel()
	utils.FailOnError(err, "Failed to open channel")

	err = channel.Confirm(false)
	utils.FailOnError(err, "Failed to put channel in confirmation mode")

	declareQueue(channel)
	log.Printf("%s Queue declared", os.Getenv("MongoDBQueueName"))

	bindQueue(channel, os.Getenv("MySQLExchange"))
	log.Printf("Queue binded to %s Exchange", os.Getenv("MySQLExchange"))

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
		os.Getenv("MongoDBQueueName"),
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
		os.Getenv("MongoDBQueueName"),
		"",
		exchange,
		false,
		nil,
	)

	utils.FailOnError(err, "Failed to bind queue "+os.Getenv("MongoDBQueueName"))
}
