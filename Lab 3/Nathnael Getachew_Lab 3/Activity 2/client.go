package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Start a goroutine to listen for incoming messages
	go receiveMessages(conn)

	// Read messages from stdin and send to server
	for {
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprintf(conn, message)
	}
}

// Function to receive messages from the server
func receiveMessages(conn net.Conn) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: ", message)
	}
}
