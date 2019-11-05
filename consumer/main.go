package main

import (
	"encoding/json"
	"os"
	"log"
	"github.com/streadway/amqp"
	"github.com/dqt12hcmus/go-rbmq/shared"
)

func getMessageChannel() (<-chan amqp.Delivery, error) {
	conn, err := amqp.Dial(shared.Config.AMQPConnectionURL)
	shared.HandleError(err, "Cannot connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	shared.HandleError(err, "Cannot create a amqpChannel")
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare(
		"add",
		true,
		false,
		false,
		false,
		nil,
	)
	shared.HandleError(err, "Cannot declare 'add' queue")

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
	shared.HandleError(err, "Cannot register consumer")

	return messageChannel, err
}

func pushToWorkerPool(messageChannel <-chan amqp.Delivery) {
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

func main() {
	messageChannel, _ := getMessageChannel()
	pushToWorkerPool(messageChannel)
}