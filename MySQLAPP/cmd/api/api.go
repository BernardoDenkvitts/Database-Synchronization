package api

import (
	"fmt"
	"net/http"

	"github.com/BernardoDenkvitts/MySQLApp/route"
	"github.com/BernardoDenkvitts/MySQLApp/service"
	"github.com/BernardoDenkvitts/MySQLApp/storage"
)

type APIServer struct {
	Address string
}

func NewAPIServer(address string) *APIServer {
	return &APIServer{Address: address}
}

func (api *APIServer) Run() error {
	router := http.NewServeMux()

	storage, err := storage.NewMySQLStore()
	if err != nil {
		panic("Error to connect to database")
	}

	storage.Init()

	router.Handle("/mysql/", http.StripPrefix("/mysql", router))

	userService := service.NewUserService(storage)
	userRoute := route.NewUserRoute(*userService)
	userRoute.Routes(router)

	fmt.Println("Server listening at port " + api.Address)
	return http.ListenAndServe(api.Address, router)
}
