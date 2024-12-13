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

	// Perform addition
	args := Args{A: 10, B: 5}
	var reply int
	err = client.Call("Calculator.Add", &args, &reply)
	if err != nil {
		log.Fatal("Error calling Add RPC:", err)
	}
	fmt.Printf("Result of %d + %d = %d\n", args.A, args.B, reply)

	// Retrieve last result
	err = client.Call("Calculator.GetLastResult", &args, &reply)
	if err != nil {
		log.Fatal("Error calling GetLastResult RPC:", err)
	}
	fmt.Printf("Last result stored on server: %d\n", reply)

	// Perform subtraction
	args = Args{A: 10, B: 3}
	err = client.Call("Calculator.Subtract", &args, &reply)
	if err != nil {
		log.Fatal("Error calling Subtract RPC:", err)
	}
	fmt.Printf("Result of %d - %d = %d\n", args.A, args.B, reply)

	// Retrieve last result
	err = client.Call("Calculator.GetLastResult", &args, &reply)
	if err != nil {
		log.Fatal("Error calling GetLastResult RPC:", err)
	}
	fmt.Printf("Last result stored on server: %d\n", reply)

	// Perform division
	args = Args{A: 20, B: 4}
	err = client.Call("Calculator.Divide", &args, &reply)
	if err != nil {
		log.Fatal("Error calling Divide RPC:", err)
	}
	fmt.Printf("Result of %d / %d = %d\n", args.A, args.B, reply)

	// Retrieve last result
	err = client.Call("Calculator.GetLastResult", &args, &reply)
	if err != nil {
		log.Fatal("Error calling GetLastResult RPC:", err)
	}
	fmt.Printf("Last result stored on server: %d\n", reply)

	// // Perform multiplication
	args = Args{A: 6, B: 7}
	err = client.Call("Calculator.Multiply", &args, &reply)
	if err != nil {
		log.Fatal("Error calling Multiply RPC:", err)
	}
	fmt.Printf("Result of %d * %d = %d\n", args.A, args.B, reply)

	// Retrieve last result
	err = client.Call("Calculator.GetLastResult", &args, &reply)
	if err != nil {
		log.Fatal("Error calling GetLastResult RPC:", err)
	}
	fmt.Printf("Last result stored on server: %d\n", reply)
}
