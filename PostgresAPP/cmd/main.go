package main

import "github.com/BernardoDenkvitts/PostgresAPP/cmd/api"

func main() {

	api := api.NewAPIServer(":8282")
	if err := api.Run(); err != nil {
		panic("Error to initialize server")
	}
}
