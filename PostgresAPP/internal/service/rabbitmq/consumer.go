package rabbitmq

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/BernardoDenkvitts/PostgresAPP/internal/infra"
	"github.com/BernardoDenkvitts/PostgresAPP/internal/types"
	"github.com/BernardoDenkvitts/PostgresAPP/internal/utils"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

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

	msgs := registerPostgresSQLConsumer(rmq.amqpChannel)

	for newUsers := range msgs {

		time.Sleep(30 * time.Second)

		log.Println("(POSTGRESSQL APP) Getting new users")

		var users []*types.User
		if err := json.Unmarshal(newUsers.Body, &users); err != nil {
			utils.FailOnError(err, "(POSTGRESSQL APP) Failed to unmarshal new users")
		}

		for _, user := range users {
			rmq.userStorage.CreateUserInformation(user)
			log.Printf("New user saved -> %s", *user)
		}

		newUsers.Ack(false)

		log.Println("(POSTGRESSQL APP) Latest users saved")
	}
}

func registerPostgresSQLConsumer(channel *amqp.Channel) <-chan amqp.Delivery {
	path, _ := os.Getwd()
	err := godotenv.Load(path + "/../.env")
	utils.FailOnError(err, "Failed to load env file")

	msgs, err := channel.Consume(
		os.Getenv("PostgresSQLQueueName"),
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "(POSTGRESSQL APP) Failed to register consumer")

	return msgs
}
