// The client performs the lookup only once at startup (caching requirement),
// then executes a request loop distributing the load in a round-robin strategy.

package main

import (
	"service-registry-client-lb/common"
	"log"
	"net/rpc"
	"time"
)

const registryAddr = "localhost:8000"

func main() {

	// Get servers once (caching requirement)
	servers := getServers()

	if len(servers) == 0 {
		log.Fatal("[Client] No servers available. Exiting.")
	}

	log.Printf("[Client] Found servers: %v. Starting Load Balancing Session...", servers)
	log.Println("---------------------------------------------------------")

	// Simulate a session of 10 requests using Round-Robin
	for i := 0; i < 10; i++ {
	
		// Round-robin selection
		serverIndex := i % len(servers)
		targetServer := servers[serverIndex]

		log.Printf("[Client] [Request %d] Selected Server: %s", i+1, targetServer)
		
		callService(targetServer, i, i*2) 

		time.Sleep(1 * time.Second) 
	}
}

func getServers() []string {
	client, err := rpc.DialHTTP("tcp", registryAddr)
	if err != nil {
		log.Fatalf("[Client] Error connecting to registry: %v", err)
	}
	defer client.Close()

	var list common.ServiceList
	
	err = client.Call("Registry.GetProviders", &struct{}{}, &list)
	if err != nil {
		log.Fatalf("[Client] Error fetching providers: %v", err)
	}
	return list.Servers
}

func callService(address string, a, b int) {

	// Connect directly to the chosen server
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Printf("[Client] Error connecting to server %s: %v", address, err)
		return
	}
	defer client.Close()

	args := &common.Args{A: a, B: b}
	var reply int
	err = client.Call("MathService.Sum", args, &reply)
	if err != nil {
		log.Printf("[Client] RPC Error on server %s: %v", address, err)
		return
	}

	log.Printf(" -> Result from %s: %d + %d = %d", address, a, b, reply)
}