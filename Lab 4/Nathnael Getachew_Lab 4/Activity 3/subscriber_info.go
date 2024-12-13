package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	subject := "updates.info"
	_, err = nc.Subscribe(subject, func(m *nats.Msg) {
		fmt.Printf("Received an info message: %s\n", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Subscribed to updates.info. Waiting for messages...")
	select {} // Keep the program running
}
