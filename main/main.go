package main

import (
	"RabbitMQPractice/src"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	//Connect to RabbitMQ server
	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	src.Send("amqp://guest:guest@localhost:5672/", "message", "message body")
}
