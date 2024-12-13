package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

// Args holds the arguments for arithmetic operations
type Args struct {
	A, B int
}

func main() {
	// Connect to the RPC server
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to RPC server:", err)
	}
	defer client.Close() // Ensure the client connection is closed when done

	// Prepare the arguments for a division operation that will cause an error
	args := Args{A: 10, B: 0} // Division by zero to trigger an error
	var reply float64         // Change to float64 to match the Divide method's reply type

	// Call RPC method with a timeout
	call := client.Go("Calculator.Divide", &args, &reply, nil)
	select {
	case <-call.Done:
		if call.Error != nil {
			log.Println("RPC error:", call.Error)
		} else {
			fmt.Printf("Result: %f\n", reply)
		}
	case <-time.After(2 * time.Second): // Adjust the timeout duration as needed
		log.Println("RPC call timed out")
	}
}
