package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/BernardoDenkvitts/MongoAPP/internal/infra"
	"github.com/BernardoDenkvitts/MongoAPP/internal/types"
	"github.com/BernardoDenkvitts/MongoAPP/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "MONGODB-APP-QUEUE"

type IConsumer interface {
	Consume()
}

type RabbitMQConsumer struct {
	userStorage infra.Storage
	amqpChannel *amqp.Channel
}

func NewRabbitMQConsumer(userStorage infra.Storage, amqpChannel *amqp.Channel) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		userStorage: userStorage,
		amqpChannel: amqpChannel,
	}
}

func (rmq *RabbitMQConsumer) Consume() {
	msgs := registerMongoDBConsumer(rmq.amqpChannel)

	for newUsers := range msgs {

		log.Println("(MONGO APP) Getting new users")

		var users []*types.User
		if err := json.Unmarshal(newUsers.Body, &users); err != nil {
			utils.FailOnError(err, "(MONGO APP) Failed to unmarshal new users")
		}

		for _, user := range users {
			rmq.userStorage.CreateUserInformation(user)
			log.Printf("New user saved -> %s", *user)
		}

		newUsers.Ack(false)

		log.Println("(MONGO APP) Latest users saved")

	}

}

func registerMongoDBConsumer(channel *amqp.Channel) <-chan amqp.Delivery {
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
		utils.FailOnError(err, "(MONGO APP) Failed to register consumer")
	}

	return msgs
}
