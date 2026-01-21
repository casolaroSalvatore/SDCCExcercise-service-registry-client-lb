// The registry maintains a list of addresses in memory, protected by a mutex to handle concurrency

package main

import (
	"service-registry-client-lb/common"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

// Handles the list of active servers.
type Registry struct {
	mu sync.Mutex
	servers map[string]bool
}

// Adds a server to the registry
func (r *Registry) Register(args *common.RegistryArgs, reply *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the server is not already registered
	if _, ok := r.servers[args.Address]; ok {
		log.Printf("[Registry] Server %s is already registered", args.Address)
		*reply = false
		return nil
	}

	r.servers[args.Address] = true
	log.Printf("[Registry] Server registered: %s", args.Address)
	*reply = true
	return nil
}

// Removes a server from the registry
func (r *Registry) Deregister(args *common.RegistryArgs, reply *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the server is present
	if _, ok := r.servers[args.Address]; !ok {
		log.Printf("[Registry] Server %s not found", args.Address)
		*reply = false
		return nil
	}

	delete(r.servers, args.Address)
	log.Printf("[Registry] Server deregistered: %s", args.Address)
	*reply = true
	return nil
}

// Returns the list of active servers
func (r *Registry) GetProviders(args *struct{}, reply *common.ServiceList) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Println("[Registry] Client requested server list")
	
	list := make([]string, 0, len(r.servers))
	for addr := range r.servers {
		list = append(list, addr)
	}
	reply.Servers = list
	return nil
}

func main() {

	// Initialize the registy
	registry := &Registry{
		servers: make(map[string]bool),
	}

	err := rpc.Register(registry)
	if err != nil {
		log.Fatal("[Registry] Error registering RPC service:", err)
	}
	rpc.HandleHTTP()

	// Registry runs on fixed port 8000
	port := ":8000"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("[Registry] Listener error:", err)
	}

	log.Printf("[Registry] Service Registry running on port %s", port)
	http.Serve(listener, nil)
}