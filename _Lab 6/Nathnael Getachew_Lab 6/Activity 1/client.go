package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run client.go <hostname:port> <key:value>")
		return
	}

	server := os.Args[1]
	message := os.Args[2]

	conn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Send message to the server
	fmt.Fprintf(conn, message+"\n")

	// Read response from the server (optional)
	reader := bufio.NewReader(conn)
	response, _ := reader.ReadString('\n')
	fmt.Println("Response from server:", response)
}
