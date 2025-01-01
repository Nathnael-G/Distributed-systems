package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// Args holds the arguments for arithmetic operations
type Args struct {
	A, B int
}

// Calculator provides methods for arithmetic operations and maintains the last result
type Calculator struct {
	lastResult int
	mu         sync.Mutex
}

// Add adds two integers and returns the result
func (c *Calculator) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	c.mu.Lock()
	c.lastResult = *reply
	c.mu.Unlock()
	return nil
}

// Subtract subtracts two integers and returns the result
func (c *Calculator) Subtract(args *Args, reply *int) error {
	*reply = args.A - args.B
	c.mu.Lock()
	c.lastResult = *reply
	c.mu.Unlock()
	return nil
}

// Multiply multiplies two integers and returns the result
func (c *Calculator) Multiply(args *Args, reply *int) error {
	if args.A == 0 || args.B == 0 {
		return errors.New("multiplication by zero is not allowed")
	}
	*reply = args.A * args.B
	c.mu.Lock()
	c.lastResult = *reply
	c.mu.Unlock()
	return nil
}

// Divide divides two integers and returns the result
func (c *Calculator) Divide(args *Args, reply *int) error {
	if args.B == 0 {
		return errors.New("division by zero is not allowed")
	}
	*reply = args.A / args.B
	c.mu.Lock()
	c.lastResult = *reply
	c.mu.Unlock()
	return nil
}

// GetLastResult retrieves the last result stored on the server
func (c *Calculator) GetLastResult(args *Args, reply *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	*reply = c.lastResult
	return nil
}

func main() {
	calc := new(Calculator)
	rpc.Register(calc)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}
	fmt.Println("RPC server is listening on port 1234...")
	rpc.Accept(listener) // Block and serve clients
}
