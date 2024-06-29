package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/BernardoDenkvitts/MySQLApp/internal/infra"
	"github.com/BernardoDenkvitts/MySQLApp/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

const exchangeName = "MYSQL"

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

	go func() {

		for {
			time.Sleep(5 * time.Minute)

			latestUsers, _ := rmq.userStorage.GetLatestUserInformations()

			if len(latestUsers) == 0 {
				log.Println("No latest user informations available at " + time.Now().Format(time.RFC3339))
				continue
			}

			fmt.Println(latestUsers)

			latestUsersJSON, err := json.Marshal(latestUsers)
			utils.FailOnError(err, "Fail to marshal the data")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			retrys := 0
			for retrys < 3 {
				publishConfirmation, err := rmq.publishUsers(ctx, latestUsersJSON)
				if err != nil {
					log.Println("Fail to publish the data -> " + err.Error())
					retrys += 1
					continue
				}
				if _, err := publishConfirmation.WaitContext(ctx); err != nil {
					log.Println("Failed to receive publish confirmation -> " + err.Error())
					retrys += 1
					continue
				}
				break
			}
			if retrys == 3 {
				panic("Error to publish the data")
			}

			log.Println("Latest users sent to exchange !!")
		}

	}()

}

func (rmq *RabbitMQProducer) publishUsers(ctx context.Context, latestUsers []byte) (*amqp.DeferredConfirmation, error) {
	return rmq.amqpChannel.PublishWithDeferredConfirmWithContext(
		ctx,
		exchangeName,
		"fanout",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        latestUsers,
		},
	)
}
