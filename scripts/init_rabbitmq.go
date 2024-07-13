/*

	This script is necessary to declare rabbitmq exchanges
	before to start the others applications, otherwise will result in
	applications encountering undefined exchanges references

*/

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	/*

		Trys was used because it tries to connect
		to rabbitmq before it is fully ready.
		After some secods the connection is done

	*/

	trys := 0

	for trys < 3 {
		time.Sleep(5 * time.Second)

		path, _ := os.Getwd()
		err := godotenv.Load(path + "/.env")
		if err != nil {
			panic("Failed to load env file ->" + err.Error())
		}

		host := "rabbitmq"

		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", os.Getenv("user"), os.Getenv("RBMQpassword"), host, os.Getenv("port")))
		if err != nil {
			fmt.Println("Error to connect to rabbitmq ->" + err.Error())
			trys += 1
			continue
		}

		channel, err := conn.Channel()
		if err != nil {
			panic("Failed to open channel ->" + err.Error())
		}

		declareExchange(channel, os.Getenv("MySQLExchange"))
		declareExchange(channel, os.Getenv("MongoDBExchange"))
		declareExchange(channel, os.Getenv("PostgresSQLExchange"))

		return
	}

	panic("Error to connect to rabbitmq")
}

func declareExchange(channel *amqp.Channel, exchangeName string) {
	log.Printf("Declaring exchange %s", exchangeName)
	err := channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic("Failed to declare exchange ->" + err.Error())
	}

	log.Printf("Exchange %s declared", exchangeName)
}
