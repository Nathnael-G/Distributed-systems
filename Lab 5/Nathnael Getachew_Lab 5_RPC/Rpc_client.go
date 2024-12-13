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
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to RPC server:", err)
	}
	defer client.Close()

	// Example of performing operations
	args := Args{A: 10, B: 5}

	// Perform addition
	var addReply int
	err = client.Call("Calculator.Add", &args, &addReply)
	if err != nil {
		log.Fatal("Error calling Add RPC:", err)
	}
	fmt.Printf("Result of %d + %d = %d\n", args.A, args.B, addReply)

	// Get last result
	var lastResult int
	err = client.Call("Calculator.GetLastResult", &struct{}{}, &lastResult)
	if err != nil {
		log.Fatal("Error calling GetLastResult RPC:", err)
	}
	fmt.Printf("Last Result: %d\n", lastResult)

	// Perform subtraction
	var subReply int
	err = client.Call("Calculator.Subtract", &args, &subReply)
	if err != nil {
		log.Fatal("Error calling Subtract RPC:", err)
	}
	fmt.Printf("Result of %d - %d = %d\n", args.A, args.B, subReply)

	// Get last result
	err = client.Call("Calculator.GetLastResult", &struct{}{}, &lastResult)
	if err != nil {
		log.Fatal("Error calling GetLast Result RPC:", err)
	}
	fmt.Printf("Last Result: %d\n", lastResult)

	// Perform multiplication
	var mulReply int
	err = client.Call("Calculator.Multiply", &args, &mulReply)
	if err != nil {
		log.Fatal("Error calling Multiply RPC:", err)
	}
	fmt.Printf("Result of %d * %d = %d\n", args.A, args.B, mulReply)

	// Get last result
	err = client.Call("Calculator.GetLastResult", &struct{}{}, &lastResult)
	if err != nil {
		log.Fatal("Error calling GetLastResult RPC:", err)
	}
	fmt.Printf("Last Result: %d\n", lastResult)

	// Perform division
	args.B = 2 // Change B to a non-zero value for division
	var divReply float64
	err = client.Call("Calculator.Divide", &args, &divReply)
	if err != nil {
		log.Fatal("Error calling Divide RPC:", err)
	}
	fmt.Printf("Result of %d / %d = %f\n", args.A, args.B, divReply)

	// Get last result
	err = client.Call("Calculator.GetLastResult", &struct{}{}, &lastResult)
	if err != nil {
		log.Fatal("Error calling GetLastResult RPC:", err)
	}
	fmt.Printf("Last Result: %d\n", lastResult)
}
