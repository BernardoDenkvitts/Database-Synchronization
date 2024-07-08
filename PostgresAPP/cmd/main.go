package main

import (
	"fmt"

	"github.com/BernardoDenkvitts/PostgresAPP/internal/infra"
)

func main() {
	_, err := infra.NewPostgresStore()
	if err != nil {
		panic(err)
	}
	fmt.Println("funcionando")

}
