package main

import "github.com/BernardoDenkvitts/MongoAPP/cmd/api"

func main() {

	api := api.NewApiServer(":8181")
	if err := api.Run(); err != nil {
		panic("")
	}
}
