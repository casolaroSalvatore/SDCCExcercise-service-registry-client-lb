// This file contains the definitions of the data structures used for RPC communication

package common

// Stateless Service Arguments

// Arguments for the stateless calculation service (sum)
type SumArgs struct {
	A, B int
}

// Statefull Service Arguments

// Arguments for the statefull calculation service (multiplication)
// The client sends a factor, which is multiplied against the global state stored in the DB 
type MulArgs struct {
	Factor int
}

// Database Communication 

// Used to write a value into the external database
type KeyValueArgs struct {
	Key string
	Value int
}

// Used to read a value from the external database
type KeyArgs struct {
	Key string
}

// Registry Arguments 

// Represents arguments to register/deregister a server
type RegistryArgs struct {
	Address string // Format "IP:Port"
}

// Contains the list of active servers
type ServiceList struct {
	Servers []string
}

// Constants for RPC method names
const (
	RegistryServiceName = "Registry"
	DatabaseServiceName = "KVStore"
	
	StatelessServiceName = "StatelessMath" // Sum
	StatefulServiceName = "StatefulMath"   // Multiplication
)