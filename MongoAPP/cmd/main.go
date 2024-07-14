package main

import (
	"os"

	"github.com/BernardoDenkvitts/MongoAPP/cmd/api"
	"github.com/joho/godotenv"
)

func main() {

	path, _ := os.Getwd()
	err := godotenv.Load(path + "/../.env")
	if err != nil {
		panic("Error to load env file")
	}

	api := api.NewApiServer(":8181")
	if err := api.Run(); err != nil {
		panic("")
	}
}
