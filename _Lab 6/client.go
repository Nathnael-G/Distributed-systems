package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <replica_ip:port>")
		return
	}

	replicaAddr := os.Args[1]
	conn, err := net.Dial("tcp", replicaAddr)
	if err != nil {
		fmt.Println("Error connecting to replica:", err)
		return
	}
	defer conn.Close()

	// Send a request to get the data
	fmt.Fprintf(conn, "GET\n")

	// Read the response
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(message)
	}
}
