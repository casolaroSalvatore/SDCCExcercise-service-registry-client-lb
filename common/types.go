// This file contains the definitions of the data structures used for RPC communication.

package common

// Args represents arguments for the calculation service.
type Args struct {
	A, B int
}

// RegistryArgs represents arguments to register/deregister a server.
type RegistryArgs struct {
	Address string // Format "IP:Port"
}

// ServiceList contains the list of active servers.
type ServiceList struct {
	Servers []string
}

// Constants for RPC method names
const (
	RegistryServiceName    = "Registry"
	CalculationServiceName = "MathService"
)