package main

import (
	"fmt"
	"log"

	"github.com/BernardoDenkvitts/MySQLApp/cmd/api"
)

func main() {

	api := api.NewAPIServer(":8080")
	if err := api.Run(); err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error to initialize server")
	}

}
