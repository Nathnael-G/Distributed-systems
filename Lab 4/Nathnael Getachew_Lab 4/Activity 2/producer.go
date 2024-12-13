package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	tasks := []string{"Task 1", "Task 2", "Task 3"}

	for _, task := range tasks {

		err = ch.Publish(

			"",

			q.Name,

			false,

			false,

			amqp.Publishing{

				DeliveryMode: amqp.Persistent,

				ContentType: "text/plain",

				Body: []byte(task),
			},
		)

		failOnError(err, "Failed to publish a message")

		fmt.Println("Sent:", task)
	}
}
