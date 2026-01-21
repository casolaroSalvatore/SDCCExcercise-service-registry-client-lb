// The server implements a simple addition service and multiplication service.
// On startup, it registers itself with the Registry and CTRL+C to deregister before shutting down

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

const (
	registryAddr = "localhost:8000"
	dbAddr = "localhost:8001"
)

// Stateless Service (Sum)

// Provides stateless arithmetic operations (for now only Sum)
type StatelessMath struct{}

func (m *StatelessMath) Sum(args *common.SumArgs, reply *int) error {
	log.Printf("[Server] Received request: %d + %d", args.A, args.B)
	*reply = args.A + args.B
	return nil
}

// Stateful Service (Multiplication)

// Provides stateful arithmetic operations (for now only Multiplication)
type StatefulMath struct{}

// Multiply takes a factor from the client, reads the current state from the DB, multiplies them, and saves the new state in the DB
func (s *StatefulMath) Multiply(args *common.MulArgs, reply *int) error {
	
	// Connect to the database
	client, err := rpc.DialHTTP("tcp", dbAddr)
	if err != nil {
		log.Printf("[Server] Error connecting to DB: %v", err)
		return err
	}
	defer client.Close()

	// Read current state from DB. It is used a fixed key "running_product" to simulate a shared global calculation
	key := "running_product"
	var currentVal int
	err = client.Call("KVStore.Get", &common.KeyArgs{Key: key}, &currentVal)
	if err != nil {
		return err
	}

	// If DB is empty (0), we start with 1 to allow multiplication
	if currentVal == 0 {
		currentVal = 1
	}

	newVal := currentVal * args.Factor

	// Write new state back to DB
	var setReply bool
	err = client.Call("KVStore.Set", &common.KeyValueArgs{Key: key, Value: newVal}, &setReply)
	if err != nil {
		return err
	}

	*reply = newVal
	log.Printf("[Server] Multiply: Old(%d) * Factor(%d) = New(%d)", currentVal, args.Factor, newVal)
	return nil
}

func main() {

	p := flag.String("port", "9000", "Port to listen on")
	flag.Parse()
	myAddr := "localhost:" + *p

	// Register both services
	rpc.Register(new(StatelessMath))
	rpc.Register(new(StatefulMath))
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":"+*p)
	if err != nil {
		log.Fatalf("[Server] Error listening on port %s: %v", *p, err)
	}

	// Register to service registry
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