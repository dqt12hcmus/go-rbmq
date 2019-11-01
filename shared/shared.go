package shared

import (
	"log"
)

type Configuration struct {
	AMQPConnectionURL string
}

type AddTask struct {
	Number1 int
	Number2 int
}

var Config = Configuration{
	AMQPConnectionURL: "amqp://guest:guest@rabbitmq:5672/",
}

func HandleError(err error, msg string) {
	if (err != nil) {
		log.Fatalf("%s : %s", msg, err)
	}
}