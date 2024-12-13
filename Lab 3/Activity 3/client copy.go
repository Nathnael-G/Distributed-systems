package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	for {
		// Receive task (number) from server
		task, _ := bufio.NewReader(conn).ReadString('\n')
		task = strings.TrimSpace(task)

		// Perform task (square the number)
		num, _ := strconv.Atoi(task)
		result := num * num

		// Send result back to server
		fmt.Fprintf(conn, "%d\n", result)
	}
}
