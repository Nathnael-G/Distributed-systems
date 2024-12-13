package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

// Slice to hold all connected clients
var clients = make(map[net.Conn]bool)
var mu sync.Mutex // Mutex to ensure safe concurrent access to client list

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Chat server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Add new client to client list
		mu.Lock()
		clients[conn] = true
		mu.Unlock()

		// Handle each client in a separate goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		mu.Lock()
		delete(clients, conn) // Remove client from list when they disconnect
		mu.Unlock()
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		// Read message from client
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", err)
			return
		}

		fmt.Printf("Broadcasting message: %s", message)
		// Broadcast message to all connected clients
		broadcastMessage(message, conn)
	}
}

// Send the message to all clients except the sender
func broadcastMessage(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients {
		if client != sender {
			client.Write([]byte(message))
		}
	}
}
