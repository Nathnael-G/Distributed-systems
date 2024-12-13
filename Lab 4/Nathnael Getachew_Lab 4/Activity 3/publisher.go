package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to the NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	subjects := []string{"updates.info", "updates.error"}
	for _, subject := range subjects {
		message := fmt.Sprintf("Message for %s", subject)
		if err := nc.Publish(subject, []byte(message)); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sent:", message)
	}
}
