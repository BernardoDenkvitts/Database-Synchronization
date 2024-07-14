package main

import (
	"os"

	"github.com/BernardoDenkvitts/PostgresAPP/cmd/api"
	"github.com/joho/godotenv"
)

func main() {
	path, _ := os.Getwd()
	err := godotenv.Load(path + "/../.env")
	if err != nil {
		panic("Error to load env file")
	}

	api := api.NewAPIServer(":8282")
	if err := api.Run(); err != nil {
		panic("Error to initialize server")
	}
}
