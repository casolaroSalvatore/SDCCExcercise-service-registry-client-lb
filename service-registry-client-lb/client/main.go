// The client performs the lookup only once at startup (caching),
// then executes a request loop distributing the load in a round-robin strategy

package main

import (
	"bufio"
	"fmt"
	"service-registry-client-lb/common"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

const registryAddr = "localhost:8000"

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
	
	// List of servers
	return list.Servers
}

// Invokes the Sum method
func callSum(address string, a, b int) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Printf("[Client] Error connecting to server %s: %v", address, err)
		return
	}
	defer client.Close()

	args := &common.SumArgs{A: a, B: b}
	var reply int
	
	err = client.Call("StatelessMath.Sum", args, &reply)
	if err != nil {
		log.Printf("[Client] RPC Sum Error: %v", err)
		return
	}

	log.Printf("[Client] Sum Result: %d + %d = %d", a, b, reply)
}

// Invokes the Multiply method
func callMultiply(address string, factor int) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Printf("[Client] Error connecting to server %s: %v", address, err)
		return
	}
	defer client.Close()

	args := &common.MulArgs{Factor: factor}
	var reply int

	err = client.Call("StatefulMath.Multiply", args, &reply)
	if err != nil {
		log.Printf("[Client] RPC Multiply Error: %v", err)
		return
	}

	log.Printf(" [Client] Global Accumulator: * %d = %d", factor, reply)
}

func main() {

	// Get servers once 
	servers := getServers()

	if len(servers) == 0 {
		log.Fatal("[Client] No servers available. Exiting.")
	}

	log.Printf("[Client] Found servers: %v", servers)
	log.Println("Starting session. Type 'exit' to quit.")

	// Read inpuut
	reader := bufio.NewReader(os.Stdin)
	
	// Counter to handle Round Robin
	requestCount := 0

	for {
		fmt.Println("\n--- Client MENU ---")
		fmt.Println("1. Execute Stateless Sum (A + B)")
		fmt.Println("2. Execute Stateful Multiplication (GlobalAccumulator * Factor)")
		fmt.Println("3. Exit")
		fmt.Print("[Client] Select an option: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "3" || input == "exit" {
			log.Println("[Client] Exiting session...")
			break
		}

		// Round-robin selection
		serverIndex := requestCount % len(servers)
		targetServer := servers[serverIndex]
		
		switch input {
		case "1":
			// Sum
			fmt.Println("\n[Client] [Stateless Operation Selected]")
			a := readInt(reader, "Enter value for A: ")
			b := readInt(reader, "Enter value for B: ")

			log.Printf("[Client] [Request %d] Routing to Server: %s", requestCount+1, targetServer)
			callSum(targetServer, a, b)
			requestCount++

		case "2":
			// Multiply
			fmt.Println("\n[Client] [Stateful Operation Selected]")
			factor := readInt(reader, "Enter Factor to multiply by: ")

			log.Printf("[Client] [Request %d] Routing to Server: %s", requestCount+1, targetServer)
			callMultiply(targetServer, factor)
			requestCount++

		default:
			fmt.Println("[Client] Invalid option, please try again.")
		}
	}
}

// Helper in order to read chars from keayboard
func readInt(reader *bufio.Reader, prompt string) int {
	for {
		fmt.Print(prompt)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		
		val, err := strconv.Atoi(text)
		if err == nil {
			return val
		}
		fmt.Println("[Client] Invalid number. Please enter an integer.")
	}
}