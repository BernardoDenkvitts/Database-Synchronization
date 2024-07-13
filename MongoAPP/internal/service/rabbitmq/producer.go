package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/BernardoDenkvitts/MongoAPP/internal/infra"
	"github.com/BernardoDenkvitts/MongoAPP/internal/utils"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type IProducer interface {
	Produce()
}

type RabbitMQProducer struct {
	userStorage infra.Storage
	amqpChannel *amqp.Channel
}

func NewRabbitMQProducer(userStorage infra.Storage, amqpChannel *amqp.Channel) *RabbitMQProducer {
	return &RabbitMQProducer{
		userStorage: userStorage,
		amqpChannel: amqpChannel,
	}
}

func (rmq *RabbitMQProducer) Produce() {

	for {
		time.Sleep(30 * time.Second)

		latestUsers, _ := rmq.userStorage.GetLatestUserInformations()

		if len(latestUsers) == 0 {
			log.Println("(MONGODB APP) No latest user informations available at " + time.Now().UTC().String())
			continue
		}

		log.Println("Users to be sent to exchange: ")
		for _, user := range latestUsers {
			log.Println(user)
		}

		latestUsersJSON, err := json.Marshal(latestUsers)
		utils.FailOnError(err, "(MONGODB APP) Failed to marshal the data")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		retrys := 0
		for retrys < 3 {
			publishConfirmation, err := rmq.publishUsers(ctx, latestUsersJSON)
			if err != nil {
				log.Println("(MONGODB APP) Failed to publish the data -> " + err.Error())
				retrys += 1
				continue
			}

			if _, err := publishConfirmation.WaitContext(ctx); err != nil {
				log.Println("(MONGODB APP) Failed to receive publish confirmation -> " + err.Error())
				retrys += 1
				continue
			}
			break
		}
		if retrys == 3 {
			panic("(MONGODB APP) Error to publish the data")
		}

		log.Println("(MONGODB APP) Latest users sent to exchange !!")
	}

}

func (rmq *RabbitMQProducer) publishUsers(ctx context.Context, latestUsers []byte) (*amqp.DeferredConfirmation, error) {
	path, _ := os.Getwd()
	err := godotenv.Load(path + "/../.env")
	utils.FailOnError(err, "Failed to load env file")

	return rmq.amqpChannel.PublishWithDeferredConfirmWithContext(
		ctx,
		os.Getenv("MongoDBExchange"),
		"fanout",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        latestUsers,
		},
	)
}
