package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BernardoDenkvitts/MySQLApp/cmd/api"
	"github.com/joho/godotenv"
)

func main() {
	path, _ := os.Getwd()
	err := godotenv.Load(path + "/../.env")
	if err != nil {
		panic("Error to load env file")
	}
	api := api.NewAPIServer(":8080")
	if err := api.Run(); err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error to initialize server")
	}

}
