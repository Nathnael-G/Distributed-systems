package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
)

type Replica struct {
	data    map[string]string
	mu      sync.Mutex
	peers   []string // List of peer addresses
	ackLock sync.Mutex
	acks    map[string]int // Track acknowledgments
	quorum  int            // Quorum size
}

type Args struct {
	Key   string
	Value string
}

func (r *Replica) Update(args *Args, reply *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.data == nil {
		r.data = make(map[string]string) // Ensure data map is initialized
	}
	log.Printf("Initiating update: key=%s, value=%s\n", args.Key, args.Value)
	r.data[args.Key] = args.Value
	*reply = true
	return nil
}

func (r *Replica) propagateUpdates(key, value string) {
	r.ackLock.Lock()
	if r.acks == nil {
		r.acks = make(map[string]int) // Ensure acks map is initialized
	}
	r.acks[key] = 0
	r.ackLock.Unlock()
	for _, peer := range r.peers {
		go func(peer string) {
			client, err := rpc.Dial("tcp", peer)
			if err != nil {
				log.Printf("Error connecting to peer: %s, error: %v\n", peer, err)
				return
			}
			defer client.Close()
			args := &Args{Key: key, Value: value}
			var reply bool
			err = client.Call("Replica.Update", args, &reply)
			if err == nil && reply {
				r.ackLock.Lock()
				r.acks[key]++
				r.ackLock.Unlock()
				log.Printf("Acknowledged by peer: %s\n", peer)
			} else {
				log.Printf("Failed to get acknowledgment from peer: %s\n", peer)
			}
		}(peer)
	}
}

func (r *Replica) waitForAcknowledgments(key string) {
	for {
		r.ackLock.Lock()
		if r.acks[key] >= r.quorum {
			r.ackLock.Unlock()
			break
		}
		r.ackLock.Unlock()
	}
	log.Printf("Update committed for key=%s after receiving %d acknowledgments\n", key, r.quorum)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run replica_strong.go <machine_ip:port> <quorum_size> <peer1_ip:port> [<peer2_ip:port>...]")
		return
	}
	// Parse command-line arguments
	machineAddr := os.Args[1]
	quorumSize := os.Args[2]
	peers := os.Args[3:]

	// Convert quorumSize to int
	quorum, err := strconv.Atoi(quorumSize)
	if err != nil {
		log.Fatalf("Invalid quorum size: %s\n", quorumSize)
	}
	// Initialize the replica
	replica := &Replica{
		data:   make(map[string]string),
		peers:  peers,
		acks:   make(map[string]int),
		quorum: quorum,
	}
	rpc.Register(replica)
	// Start the RPC server
	listener, err := net.Listen("tcp", machineAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	log.Printf("Replica RPC Server listening on %s\n", machineAddr)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()
	// Simulate a strong consistency update
	key, value := "key1", "value1"
	var reply bool
	replica.Update(&Args{Key: key, Value: value}, &reply)
	replica.propagateUpdates(key, value)
	replica.waitForAcknowledgments(key)
}
