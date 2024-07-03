package api

import (
	"fmt"
	"net/http"

	"github.com/BernardoDenkvitts/MongoAPP/internal/infra"
	"github.com/BernardoDenkvitts/MongoAPP/internal/route"
	"github.com/BernardoDenkvitts/MongoAPP/internal/service"
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

	router.Handle("/mongodb/", http.StripPrefix("/mongodb", router))

	storage, err := infra.NewMongoDBStore()
	if err != nil {
		return err
	}

	storage.InitMongoDB()

	userService := service.NewUserService()
	userRoute := route.NewUserRoute(userService)
	userRoute.Routes(router)

	fmt.Println("MongoDB APP Server listening at port " + api.Address)
	return http.ListenAndServe(api.Address, router)
}
