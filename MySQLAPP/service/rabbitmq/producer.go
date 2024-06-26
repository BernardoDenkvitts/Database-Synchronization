package rabbitmq

import (
	"context"
	"time"

	"github.com/BernardoDenkvitts/MySQLApp/storage"
	amqp "github.com/rabbitmq/amqp091-go"
)

const exchangeName = "MYSQL"

type IProducer interface {
	Produce()
}

type RabbitMQProducer struct {
	userStorage storage.Storage
	amqpChannel *amqp.Channel
}

func NewRabbitMQProducer(userStorage storage.Storage, amqpChannel *amqp.Channel) *RabbitMQProducer {
	return &RabbitMQProducer{
		userStorage: userStorage,
		amqpChannel: amqpChannel,
	}
}

func (rmq *RabbitMQProducer) Produce() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	latestUsers, err := rmq.userStorage.GetLatestUserInformations()
	if err != nil {
		panic(err)		
	}

	err := rmq.amqpChannel.PublishWithContext(
		ctx,
		exchangeName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: ,
		}
	)
}
