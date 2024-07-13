package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BernardoDenkvitts/PostgresAPP/internal/infra"
	"github.com/BernardoDenkvitts/PostgresAPP/internal/utils"
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
			log.Println("(POSTGRESSQL APP) No latest user informations available at " + time.Now().Format(time.RFC3339))
			continue
		}

		log.Println("Users to be sent to exchange: ")
		for _, user := range latestUsers {
			fmt.Println(user)
		}

		latestUsersJSON, err := json.Marshal(latestUsers)
		utils.FailOnError(err, "(POSTGRESSQL APP) Fail to marshal the data")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		retrys := 0
		for retrys < 3 {
			publishConfirmation, err := rmq.publishUsers(ctx, latestUsersJSON)
			if err != nil {
				log.Println("(POSTGRESSQL APP) Failed to publish the data -> " + err.Error())
				retrys += 1
				continue
			}

			if _, err := publishConfirmation.WaitContext(ctx); err != nil {
				log.Println("(POSTGRESSQL APP) Failed to receive publish confirmation -> " + err.Error())
				retrys += 1
				continue
			}
			break
		}
		if retrys == 3 {
			panic("(POSTGRESSQL APP) Error to publish the data")
		}

		log.Println("(POSTGRESSQL APP) Latest users sent to exchange !!")
	}

}

func (rmq *RabbitMQProducer) publishUsers(ctx context.Context, latestUsers []byte) (*amqp.DeferredConfirmation, error) {
	path, _ := os.Getwd()
	err := godotenv.Load(path + "/../.env")
	utils.FailOnError(err, "Failed to load env file")

	return rmq.amqpChannel.PublishWithDeferredConfirmWithContext(
		ctx,
		os.Getenv("PostgresSQLExchange"),
		"fanout",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        latestUsers,
		},
	)
}
