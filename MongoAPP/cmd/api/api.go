package api

import (
	"fmt"
	"net/http"

	"github.com/BernardoDenkvitts/MongoAPP/internal/infra"
	"github.com/BernardoDenkvitts/MongoAPP/internal/route"
	"github.com/BernardoDenkvitts/MongoAPP/internal/service"
	"github.com/BernardoDenkvitts/MongoAPP/internal/service/rabbitmq"
	"github.com/BernardoDenkvitts/MongoAPP/internal/utils"
)

type ApiServer struct {
	Address string
}

func NewApiServer(address string) *ApiServer {
	return &ApiServer{
		Address: address,
	}
}

func (api *ApiServer) Run() error {
	router := http.NewServeMux()

	storage := setupDatabase()

	rabbitMQ, err := infra.NewRabbitMQ()
	utils.FailOnError(err, "(MONGO APP) Error to instanciate RabbitMQ")
	defer rabbitMQ.Close()

	rabbitMqProducer := rabbitmq.NewRabbitMQProducer(storage, rabbitMQ.Channel)
	go rabbitMqProducer.Produce()

	rabbitMqConsumer := rabbitmq.NewRabbitMQConsumer(storage, rabbitMQ.Channel)
	go rabbitMqConsumer.Consume()

	router.Handle("/mongodb/", http.StripPrefix("/mongodb", router))

	userService := service.NewUserServiceImpl(storage)
	userRoute := route.NewUserRoutesImpl(userService)
	userRoute.Routes(router)

	fmt.Println("MongoDB APP Server listening at port " + api.Address)
	return http.ListenAndServe(api.Address, router)
}

func setupDatabase() infra.Storage {
	storage, err := infra.NewMongoDBStore()
	utils.FailOnError(err, "(MONGO APP) Error to connect to database")
	storage.Init()

	return storage
}
