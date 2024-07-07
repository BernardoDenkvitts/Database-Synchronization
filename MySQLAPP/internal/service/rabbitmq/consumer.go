package rabbitmq

import (
	"encoding/json"
	"log"
	"time"

	"github.com/BernardoDenkvitts/MySQLApp/internal/infra"
	"github.com/BernardoDenkvitts/MySQLApp/internal/types"
	"github.com/BernardoDenkvitts/MySQLApp/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "MYSQL-APP-QUEUE"

type IConsumer interface {
	Consume()
}

type RabbitMQConsumer struct {
	userStorage infra.Storage
	amqpChannel *amqp.Channel
}

func NewRabbitMQConsumer(userStorage infra.Storage, ampqChannel *amqp.Channel) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		userStorage: userStorage,
		amqpChannel: ampqChannel,
	}
}

func (rmq *RabbitMQConsumer) Consume() {

	msgs := registerMySQLConsumer(rmq.amqpChannel)

	for newUsers := range msgs {

		time.Sleep(1 * time.Minute)

		log.Println("(MYSQL APP) Getting new users")

		var users []*types.User
		if err := json.Unmarshal(newUsers.Body, &users); err != nil {
			utils.FailOnError(err, "(MYSQL APP) Failed to unmarshal new users")
		}

		for _, user := range users {
			rmq.userStorage.CreateUserInformation(user)
			log.Printf("New user saved -> %s", *user)
		}

		newUsers.Ack(false)

		log.Println("(MYSQL APP) Latest users saved")

	}
}

func registerMySQLConsumer(channel *amqp.Channel) <-chan amqp.Delivery {
	msgs, err := channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		utils.FailOnError(err, "(MYSQL APP) Failed to register consumer")
	}

	return msgs
}
