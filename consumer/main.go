package main

import (
	"encoding/json"
	"os"
	"log"
	"github.com/streadway/amqp"
	"github.com/dqt12hcmus/go-rbmq/shared"
)

func handleError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s : %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial(shared.Config.AMQPConnectionURL)
	handleError(err, "Cannot connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Cannot create a amqpChannel")
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare(
		"add",
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Cannot declare 'add' queue")

	err = amqpChannel.Qos(
		1,
		0,
		false,
	)

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Cannot register consumer")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for message := range messageChannel {
			log.Printf("Received a message: %s", message.Body)
			addTask := new(shared.AddTask)
			err := json.Unmarshal(message.Body, addTask)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}
			var number1 = addTask.Number1
			var number2 = addTask.Number2
			log.Printf(
				"Result of %d + %d = %d",
				number1,
				number2,
				number1 + number2,
			)

			if err := message.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else { 
				log.Printf("Acknowledged message")
			}

		}
	}()

	<-stopChan
}