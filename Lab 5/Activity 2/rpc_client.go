package main

import (
	"fmt"
	"log"
	"net/rpc"
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
	defer client.Close()

	// Prepare the arguments
	args := Args{A: 10, B: 5}

	// Call Add method
	var addReply int
	err = client.Call("Calculator.Add", &args, &addReply)
	if err != nil {
		log.Fatal("Error calling Add RPC:", err)
	}
	fmt.Printf("Result of %d + %d = %d\n", args.A, args.B, addReply)

	// Call Subtract method
	var subReply int
	err = client.Call("Calculator.Subtract", &args, &subReply)
	if err != nil {
		log.Fatal("Error calling Subtract RPC:", err)
	}
	fmt.Printf("Result of %d - %d = %d\n", args.A, args.B, subReply)

	// Call Multiply method
	var mulReply int
	err = client.Call("Calculator.Multiply", &args, &mulReply)
	if err != nil {
		log.Fatal("Error calling Multiply RPC:", err)
	}
	fmt.Printf("Result of %d * %d = %d\n", args.A, args.B, mulReply)

	// Call Divide method
	var divReply float64
	err = client.Call("Calculator.Divide", &args, &divReply)
	if err != nil {
		log.Fatal("Error calling Divide RPC:", err)
	}
	fmt.Printf("Result of %d / %d = %f\n", args.A, args.B, divReply)
}
