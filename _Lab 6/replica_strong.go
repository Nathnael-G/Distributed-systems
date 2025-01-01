package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
)

// Replica represents a single replica in the distributed system.
type Replica struct {
	data    map[string]string // Store data as key-value pairs
	mu      sync.Mutex        // Mutex for data access
	peers   []string          // List of peer addresses
	ackLock sync.Mutex        // Mutex for acknowledgment tracking
	acks    map[string]int    // Track acknowledgments for each key
}

// Args holds the arguments for the Update RPC method.
type Args struct {
	Key   string
	Value string
}

// Update handles incoming update requests from peers.
func (r *Replica) Update(args *Args, reply *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[args.Key] = args.Value
	*reply = true
	return nil
}

// propagateUpdates sends updates to all peer replicas.
func (r *Replica) propagateUpdates(key, value string) {
	r.ackLock.Lock()
	r.acks[key] = 0 // Reset acknowledgment count for this key
	r.ackLock.Unlock()

	for _, peer := range r.peers {
		go func(peer string) {
			client, err := rpc.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer:", peer, err)
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
			}
		}(peer)
	}
}

// waitForAcknowledgments waits until the required number of acknowledgments are received.
func (r *Replica) waitForAcknowledgments(key string, required int) {
	for {
		r.ackLock.Lock()
		if r.acks[key] >= required {
			r.ackLock.Unlock()
			break
		}
		r.ackLock.Unlock()
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run replica_strong.go <machine_ip:port> <peer1_ip:port> [<peer2_ip:port>...]")
		return
	}

	machineAddr := os.Args[1]
	peers := os.Args[2:]

	replica := &Replica{
		data:  make(map[string]string),
		peers: peers,
		acks:  make(map[string]int),
	}

	rpc.Register(replica)

	listener, err := net.Listen("tcp", machineAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Replica RPC Server listening on %s\n", machineAddr)

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
	replica.Update(&Args{Key: key, Value: value}, nil)
	replica.propagateUpdates(key, value)
	replica.waitForAcknowledgments(key, len(replica.peers))
	fmt.Println("Update committed after receiving acknowledgments")
}
