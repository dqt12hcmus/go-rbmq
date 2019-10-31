package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/dqt12hcmus/go-rbmq/shared"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial(shared.Config.AMQPConnectionURL)
	shared.HandleError(err, "Cannot connect to AMQP")
	defer conn.Close()
	channel, err := conn.Channel()
	shared.HandleError(err, "Cannot create channel")
	defer channel.Close()
	queue, err := channel.QueueDeclare(
		"add",
		true,
		false,
		false,
		false,
		nil,
	)
	shared.HandleError(
		err,
		"Cannot declare `add` queue",
	)
	for {
		rand.Seed(time.Now().UnixNano())
		addTask := shared.AddTask{
			rand.Intn(999),
			rand.Intn(999),
		}
		body, err := json.Marshal(addTask)
		if err != nil {
			shared.HandleError(err, "Error encoding JSON")
		}
		channel.Publish(
			"",
			queue.Name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         body,
			},
		)
		if err != nil {
			log.Fatalf("Error publishing message: %s", err)
		}

		log.Printf("AddTask: %d+%d", addTask.Number1, addTask.Number2)
	}
}
