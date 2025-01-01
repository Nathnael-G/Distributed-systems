package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Replica struct {
	data  map[string]string
	mu    sync.Mutex
	peers []string // List of peer addresses
}

// Update modifies the replica's data
func (r *Replica) Update(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = value
}

// propagateUpdates sends updates to all peers with a simulated delay
func (r *Replica) propagateUpdates(key, value string, delay time.Duration) {
	for _, peer := range r.peers {
		go func(peer string) {
			time.Sleep(delay) // Simulate network delay
			conn, err := net.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer:", peer, err)
				return
			}
			defer conn.Close()
			fmt.Fprintf(conn, "%s:%s\n", key, value)
		}(peer)
	}
}

// handleConnection processes incoming messages from peers
func handleConnection(conn net.Conn, replica *Replica) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		parts := strings.Split(strings.TrimSpace(message), ":")
		if len(parts) == 2 {
			replica.Update(parts[0], parts[1])
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run replica_eventual.go <machine_ip:port> <peer1_ip:port> [<peer2_ip:port>...]")
		return
	}

	machineAddr := os.Args[1]
	peers := os.Args[2:]

	// Initialize the replica
	replica := &Replica{
		data:  make(map[string]string),
		peers: peers,
	}

	// Start the server
	listener, err := net.Listen("tcp", machineAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Replica listening on %s\n", machineAddr)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleConnection(conn, replica)
		}
	}()

	// Measure convergence time
	startTime := time.Now() // Start timing

	// Simulate an update
	replica.Update("key1", "value1")

	// Simulate propagation with a 2-second delay (adjust as needed)
	replica.propagateUpdates("key1", "value1", 2*time.Second)

	// Wait for a sufficient amount of time for updates to propagate
	time.Sleep(10 * time.Second) // Adjust based on expected propagation time

	// Log convergence time
	duration := time.Since(startTime)
	fmt.Printf("Convergence Time: %v\n", duration)

	// Display the current state of the replica's data
	replica.mu.Lock()
	fmt.Println("Replica Data:", replica.data)
	replica.mu.Unlock()
}
