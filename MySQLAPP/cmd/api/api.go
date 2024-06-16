package api

import (
	"net/http"

	"github.com/BernardoDenkvitts/MySQLApp/storage"
)

type APIServer struct {
	Address string
	Storage storage.Storage
}

func NewAPIServer(address string, storage storage.Storage) *APIServer {
	return &APIServer{
		Address: address,
		Storage: storage,
	}
}

func (api *APIServer) Run() {
	router := http.NewServeMux()

	router.Handle("/mysql/api/v1/", http.StripPrefix("/v1", router))

}
