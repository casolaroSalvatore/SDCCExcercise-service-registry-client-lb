// The server implements a simple addition service. On startup, it registers itself with the Registry.
// It uses os.Signal to intercept CTRL+C and deregister before shutting down.

package main

import (
	"service-registry-client-lb/common"
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

// MathService provides basic arithmetic operations.
type MathService struct{}

// Sum performs addition of two integers.
func (m *MathService) Sum(args *common.Args, reply *int) error {
	log.Printf("[Server] Received request: %d + %d", args.A, args.B)
	*reply = args.A + args.B
	return nil
}

const registryAddr = "localhost:8000"

func main() {

	p := flag.String("port", "9000", "Port to listen on")
	flag.Parse()
	myAddr := "localhost:" + *p

	// Setup RPC Service
	rpc.Register(new(MathService))
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":"+*p)
	if err != nil {
		log.Fatalf("[Server] Error listening on port %s: %v", *p, err)
	}

	// Register to Service Registry
	register(myAddr)

	// Catch CTRL+C to deregister
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("[Server] Shutting down...")
		deregister(myAddr)
		os.Exit(0)
	}()

	log.Printf("[Server] Server running on %s. Waiting for requests...", myAddr)
	http.Serve(listener, nil)
}

func register(addr string) {
	client, err := rpc.DialHTTP("tcp", registryAddr)
	if err != nil {
		log.Fatalf("[Server] Failed to connect to registry: %v", err)
	}
	defer client.Close()

	args := &common.RegistryArgs{Address: addr}
	var reply bool
	err = client.Call("Registry.Register", args, &reply)
	if err != nil || !reply {
		log.Fatalf("[Server] Failed to register: %v", err)
	}
	log.Println("[Server] Successfully registered to Registry")
}

func deregister(addr string) {
	client, err := rpc.DialHTTP("tcp", registryAddr)
	if err != nil {
		log.Printf("[Server] Failed to connect to registry for deregistration: %v", err)
		return
	}
	defer client.Close()

	args := &common.RegistryArgs{Address: addr}
	var reply bool
	err = client.Call("Registry.Deregister", args, &reply)
	if err != nil {
		log.Printf("[Server] Failed to deregister: %v", err)
	} else {
		log.Println("[Server] Successfully deregistered from Registry")
	}
}