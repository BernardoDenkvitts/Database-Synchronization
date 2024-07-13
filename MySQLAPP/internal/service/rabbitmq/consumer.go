package rabbitmq

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/BernardoDenkvitts/MySQLApp/internal/infra"
	"github.com/BernardoDenkvitts/MySQLApp/internal/types"
	"github.com/BernardoDenkvitts/MySQLApp/internal/utils"
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

func NewRabbitMQConsumer(userStorage infra.Storage, ampqChannel *amqp.Channel) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		userStorage: userStorage,
		amqpChannel: ampqChannel,
	}
}

func (rmq *RabbitMQConsumer) Consume() {

	msgs := registerMySQLConsumer(rmq.amqpChannel)

	for newUsers := range msgs {

		//Necessary to avoid send data that came from other application !!
		time.Sleep(30 * time.Second)

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
	path, _ := os.Getwd()
	err := godotenv.Load(path + "/../.env")
	utils.FailOnError(err, "Failed to load env file")

	msgs, err := channel.Consume(
		os.Getenv("MySQLQueueName"),
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "(MYSQL APP) Failed to register consumer")

	return msgs
}
