package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

// Args holds the arguments for arithmetic operations
type Args struct {
	A, B int
}

// Calculator provides methods for arithmetic operations
type Calculator int

// Add adds two integers and returns the result
func (c *Calculator) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

// Subtract subtracts two integers and returns the result
func (c *Calculator) Subtract(args *Args, reply *int) error {
	*reply = args.A - args.B
	return nil
}

// Divide divides two integers and returns the result
func (c *Calculator) Divide(args *Args, reply *float64) error {
	if args.B == 0 {
		return errors.New("division by zero is not allowed")
	}
	*reply = float64(args.A) / float64(args.B)
	return nil
}

// Multiply multiplies two integers and returns the result
func (c *Calculator) Multiply(args *Args, reply *int) error {
	if args.A == 0 || args.B == 0 {
		return errors.New("multiplication by zero is not allowed")
	}
	*reply = args.A * args.B
	return nil
}

func main() {
	// Register the Calculator service
	calc := new(Calculator)
	rpc.Register(calc)

	// Start listening for incoming RPC connections
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}
	fmt.Println("RPC server is listening on port 1234...")

	// Accept incoming connections in a separate goroutine
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}
			go rpc.ServeConn(conn) // Handle each connection concurrently
		}
	}()

	// Block forever
	select {}
}
