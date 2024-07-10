package api

import (
	"fmt"
	"net/http"

	"github.com/BernardoDenkvitts/PostgresAPP/internal/infra"
	"github.com/BernardoDenkvitts/PostgresAPP/internal/route"
	"github.com/BernardoDenkvitts/PostgresAPP/internal/service"
	"github.com/BernardoDenkvitts/PostgresAPP/internal/utils"
)

type ApiServer struct {
	Address string
}

func NewAPIServer(address string) *ApiServer {
	return &ApiServer{
		Address: address,
	}
}

func (api *ApiServer) Run() error {
	router := http.NewServeMux()

	storage := setupDatabase()

	router.Handle("/postgres/", http.StripPrefix("/postgres", router))

	userService := service.NewUserServiceImpl(storage)
	userRoute := route.NewUserRoutesImpl(userService)
	userRoute.Routes(router)

	fmt.Println("Postgres APP Server listening at port " + api.Address)
	return http.ListenAndServe(api.Address, router)
}

func setupDatabase() infra.Storage {
	storage, err := infra.NewPostgresStore()
	utils.FailOnError(err, "(POSTGRESS APP) Error to connect to database")
	storage.Init()

	return storage
}