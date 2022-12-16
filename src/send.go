package src

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type Sender struct {
	connectionAddress string
	message           string
	body              string
}

func (s Sender) ConnectToServer() (conn *amqp.Connection, err error) {
	conn, err = amqp.Dial(s.connectionAddress)
	return
}

func (s Sender) OpenChannel(conn *amqp.Connection) (ch *amqp.Channel, err error) {
	ch, err = conn.Channel()
	return
}

func (s Sender) DeclareQueue(ch *amqp.Channel) (q amqp.Queue, err error) {
	q, err = ch.QueueDeclare(
		s.message,
		false,
		false,
		false,
		false,
		nil)
	return
}

func Send(addr, msg, body string) {
	sender := Sender{addr, msg, body}

	//Connect to RabbitMQ server
	conn, err := sender.ConnectToServer()
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//Open connection channel
	ch, err := sender.OpenChannel(conn)
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := sender.DeclareQueue(ch)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msgBody := sender.body
	err = ch.PublishWithContext(ctx, "", q.Name, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msgBody),
		})
	failOnError(err, "Failed to publish")
	log.Printf("[x] Sent %s\n", msgBody)

}
