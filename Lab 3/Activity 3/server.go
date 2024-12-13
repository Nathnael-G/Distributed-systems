package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

// List to hold active clients
var clients = make(map[net.Conn]bool)
var mu sync.Mutex

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is ready to assign tasks...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		mu.Lock()
		clients[conn] = true
		mu.Unlock()

		// Handle each client connection
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		conn.Close()
	}()

	for {
		// Generate a random number as a task
		task := time.Now().Unix() % 100

		// Send task to client
		fmt.Fprintf(conn, "%d\n", task)

		// Receive result from client
		response, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Received result from client:", response)

		time.Sleep(5 * time.Second) // Simulate task intervals
	}
}
