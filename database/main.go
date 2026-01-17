package main

import (
	"service-registry-client-lb/common"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

// KVStore represents the in-memory database shared across all servers.
type KVStore struct {
	mu sync.Mutex
	store map[string]int
}

// Set saves a value into the database associated with a specific key
func (k *KVStore) Set(args *common.KeyValueArgs, reply *bool) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.store[args.Key] = args.Value
	log.Printf("[Database] UPDATE: Key='%s' NewValue=%d", args.Key, args.Value)
	*reply = true
	return nil
}

// Get retrieves a value from the database.
func (k *KVStore) Get(args *common.KeyArgs, reply *int) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	val, ok := k.store[args.Key]
	if !ok {
		// Default to 0 if not found. 
		// The server logic will handle converting 0 to 1 for multiplication.
		*reply = 0
	} else {
		*reply = val
	}
	// log.Printf("[Database] READ: Key='%s' -> %d", args.Key, *reply)
	return nil
}

func main() {
	// Initialize the store map
	db := &KVStore{
		store: make(map[string]int),
	}

	rpc.Register(db)
	rpc.HandleHTTP()

	// The Database listens on a fixed port (8001) known to the Worker Servers
	port := ":8001"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[Database] Shared Memory running on port %s", port)
	http.Serve(listener, nil)
}