package api

import (
	"fmt"
	"net/http"

	"github.com/BernardoDenkvitts/MySQLApp/internal/infra"
	"github.com/BernardoDenkvitts/MySQLApp/internal/route"
	"github.com/BernardoDenkvitts/MySQLApp/internal/service"
	"github.com/BernardoDenkvitts/MySQLApp/internal/service/rabbitmq"
	"github.com/BernardoDenkvitts/MySQLApp/internal/utils"
)

type APIServer struct {
	Address string
}

func NewAPIServer(address string) *APIServer {
	return &APIServer{Address: address}
}

func (api *APIServer) Run() error {
	router := http.NewServeMux()

	storage := setupDatabase()

	rabbitMq, err := infra.NewRabbitMQ()
	utils.FailOnError(err, "(MYSQL APP) Error to instanciate RabbitMQ")
	defer rabbitMq.Close()

	rabbitMqProducer := rabbitmq.NewRabbitMQProducer(storage, rabbitMq.Channel)
	rabbitMqProducer.Produce()

	rabbitMqConsumer := rabbitmq.NewRabbitMQConsumer(storage, rabbitMq.Channel)
	rabbitMqConsumer.Consume()

	router.Handle("/mysql/", http.StripPrefix("/mysql", router))

	userService := service.NewUserService(storage)
	userRoute := route.NewUserRoute(*userService)
	userRoute.Routes(router)

	fmt.Println("MySQL APP Server listening at port " + api.Address)
	return http.ListenAndServe(api.Address, router)
}

func setupDatabase() infra.Storage {
	storage, err := infra.NewMySQLStore()
	utils.FailOnError(err, "(MYSQL APP) Error to connect to database")
	storage.Init()

	return storage
}
