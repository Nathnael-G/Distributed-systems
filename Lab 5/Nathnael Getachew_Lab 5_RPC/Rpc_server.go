package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
)

// Args holds the arguments for arithmetic operations
type Args struct {
	A, B int
}

// Calculator provides methods for arithmetic operations and maintains state
type Calculator struct {
	lastResult int
	mu         sync.Mutex
	filename   string
}

// LoadLastResult loads the last result from a file
func (c *Calculator) LoadLastResult() error {
	file, err := os.Open(c.filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, initialize lastResult to 0
			c.lastResult = 0
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c.lastResult); err != nil {
		return err
	}
	return nil
}

// SaveLastResult saves the last result to a file
func (c *Calculator) SaveLastResult() error {
	file, err := os.Create(c.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(c.lastResult)
}

// Add adds two integers and stores the result
func (c *Calculator) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	c.mu.Lock()
	c.lastResult = *reply
	c.SaveLastResult() // Save the result to file
	c.mu.Unlock()
	return nil
}

// Subtract subtracts two integers and stores the result
func (c *Calculator) Subtract(args *Args, reply *int) error {
	*reply = args.A - args.B
	c.mu.Lock()
	c.lastResult = *reply
	c.SaveLastResult() // Save the result to file
	c.mu.Unlock()
	return nil
}

// Divide divides two integers and stores the result
func (c *Calculator) Divide(args *Args, reply *float64) error {
	if args.B == 0 {
		return errors.New("division by zero is not allowed")
	}
	*reply = float64(args.A) / float64(args.B)
	c.mu.Lock()
	c.lastResult = int(*reply)
	c.SaveLastResult() // Save the result to file
	c.mu.Unlock()
	return nil
}

// Multiply multiplies two integers and stores the result
func (c *Calculator) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	c.mu.Lock()
	c.lastResult = *reply
	c.SaveLastResult() // Save the result to file
	c.mu.Unlock()
	return nil
}

// GetLastResult returns the last result of the arithmetic operations
func (c *Calculator) GetLastResult(args *struct{}, reply *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	*reply = c.lastResult
	return nil
}

func main() {
	calc := &Calculator{filename: "last_result.json"}
	err := calc.LoadLastResult() // Load the last result from file
	if err != nil {
		fmt.Println("Error loading last result:", err)
		return
	}

	err = rpc.Register(calc)
	if err != nil {
		fmt.Println("Error registering RPC:", err)
		return
	}

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("RPC server is listening on port 1234...")

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

	select {}
}