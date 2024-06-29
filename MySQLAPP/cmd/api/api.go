package api

import (
	"fmt"
	"net/http"

	"github.com/BernardoDenkvitts/MySQLApp/internal/infra"
	"github.com/BernardoDenkvitts/MySQLApp/internal/route"
	"github.com/BernardoDenkvitts/MySQLApp/internal/service"
	"github.com/BernardoDenkvitts/MySQLApp/internal/service/rabbitmq"
)

type APIServer struct {
	Address string
}

func NewAPIServer(address string) *APIServer {
	return &APIServer{Address: address}
}

func (api *APIServer) Run() error {
	router := http.NewServeMux()

	storage, err := infra.NewMySQLStore()
	if err != nil {
		panic("Error to connect to database")
	}

	storage.Init()

	rabbitMq, err := infra.NewRabbitMQ()
	if err != nil {
		panic("Error to instanciate RabbitMQ")
	}
	defer rabbitMq.Close()

	rabbitMqProducer := rabbitmq.NewRabbitMQProducer(storage, rabbitMq.Channel)

	go rabbitMqProducer.Produce()

	router.Handle("/mysql/", http.StripPrefix("/mysql", router))

	userService := service.NewUserService(storage)
	userRoute := route.NewUserRoute(*userService)
	userRoute.Routes(router)

	fmt.Println("Server listening at port " + api.Address)
	return http.ListenAndServe(api.Address, router)
}
