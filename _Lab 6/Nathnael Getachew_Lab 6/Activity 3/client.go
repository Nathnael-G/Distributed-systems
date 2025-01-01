package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

// Args holds the arguments for the RPC calls
type Args struct {
	A, B int
}

func main() {
	// Connect to the server
	client, err := rpc.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal("Error connecting to RPC server:", err)
	}

	// Set up arguments for the divide operation
	args := Args{A: 10, B: 0} // Division by zero to trigger an error
	var reply int

	// Call the Divide method on the server with a timeout
	call := client.Go("Calculator.Divide", &args, &reply, nil)
	select {
	case <-call.Done:
		if call.Error != nil {
			log.Println("RPC error:", call.Error)
		} else {
			fmt.Printf("Result: %d\n", reply)
		}
	case <-time.After(2 * time.Second):
		log.Println("RPC call timed out")
	}
}
