# Service Registry & Client-Side Load Balancing

This project implements a distributed system in Go that demonstrates the use of RPC, Service Discovery, Client-Side Load Balancing, and distributed state management via an external Key-Value Store.

## System Specification
The system adheres to the following requirements:
1.  **Service Registry:** Maintains a dynamic list of active server addresses.
2.  **RPC Server:** Offers two types of services:
    - **Stateless:** Sum service (Sum), where the result depends only on the inputs.
    - **Stateful:** Multiplication service with an accumulator (Multiply), where the state is kept consistent across replicas thanks to the external Database.
3.  **External Database (KVStore):** A separate component acting as shared memory to guarantee data persistence between replicas.
4.  **Interactive Client:** Features the following characteristics:
    - **Caching:** Downloads the server list from the Registry only once at the beginning of the session.
    - **Load Balancing:** Uses a Round-Robin algorithm to distribute requests among servers.
    - **Interface:** Offers a command-line menu to choose which service to invoke.

## File Architecture
The system is structured into the following components:
- **common/types.go**: Shared data structures for the RPC protocol (arguments for sum, multiplication, and database).
- **registry/main.go**: The central registry server, which manages the address map.
- **database/main.go**: The in-memory Key-Value Store.
- **server/main.go**: The node offering the calculation service.
- **client/main.go**: The client acting as the load balancer.

## Execution Guide
Since the execution is local and manual (no virtualization provided by Docker is used), it is necessary to open at least 5 separate terminals.

### 1. Prerequisites
- Go installed (`go version`).
- Module initialized: Run `go mod init service-registry-client-lb` and `go mod tidy` in the root folder.

### 2. Steps to Run

1.  **Terminal 1 (Registry):**
    Open a terminal in the root folder and run:

    ```bash
    go run registry/main.go
    ```

    This starts the service registry, listening on port 8000.
    *Expected output:*
    ```text
    2026/01/17 20:04:16 [Registry] Service Registry running on port :8000
    ```

2.  **Terminal 2 (Database):**
    Open a terminal in the root folder and run:

    ```bash
    go run database/main.go
    ```

    This starts the shared memory, listening on port 8001.
    *Expected output:*
    ```text
    2026/01/17 20:06:34 [Database] Shared Memory running on port :8001
    ```

3.  Now open two (or more) new terminals to simulate different servers. In this case, we use 2 servers.

    **Terminal 3 (Server A):**
    Run:
    ```bash
    go run server/main.go -port 9001
    ```
    Starts the first calculation server on port 9001.
    *Expected output:*
    ```text
    2026/01/17 20:07:28 [Server] Successfully registered to Registry
    2026/01/17 20:07:28 [Server] Server running on localhost:9001. Waiting for requests...
    ```

    **Terminal 4 (Server B):**
    Run:
    ```bash
    go run server/main.go -port 9002
    ```
    Starts the second calculation server on port 9002.
    *Expected output:*
    ```text
    2026/01/17 20:07:44 [Server] Successfully registered to Registry
    2026/01/17 20:07:44 [Server] Server running on localhost:9002. Waiting for requests...
    ```

    On the **Registry** terminal, you will see the successful registration logs:
    ```text
    2026/01/17 20:07:28 [Registry] Server registered: localhost:9001
    2026/01/17 20:07:44 [Registry] Server registered: localhost:9002
    ```

4.  **Terminal 5 (Client):**
    Open a final terminal to simulate the client and run:

    ```bash
    go run client/main.go
    ```

### 3. Result
Once the client is started, an interactive menu will be observed:

```text
--- Client MENU ---
1. Execute Stateless Sum (A + B)
2. Execute Stateful Multiplication (GlobalAccumulator * Factor)
3. Exit
[Client] Select an option:
```

### Scenario A: Stateless Test (Sum)
By selecting option **1**, the client prompts the user to enter two values in order to execute the sum operation.

```text
[Client] [Stateless Operation Selected]
Enter value for A:
Enter value for B:
```

Once the values are entered (e.g., A = 1 and B = 1), the sum result is returned by Server A:

```bash
2026/01/17 20:11:01 [Client] [Request 1] Routing to Server: localhost:9001
2026/01/17 20:11:01 [Client] Sum Result: 1 + 1 = 2
```

By executing the sum operation again (still with A = 1 and B = 1), it can be observed that the request is now handled by Server B:

```bash
2026/01/17 20:14:05 [Client] [Request 2] Routing to Server: localhost:9002
2026/01/17 20:14:05 [Client] Sum Result: 1 + 1 = 2
```

Therefore, requests are alternately routed to Server A and Server B according to the Round-Robin strategy, while the computation remains local and isolated for each request.

### Scenario B: Stateful Test (Multiplication)
By selecting option **2**, the user can enter a multiplication factor. The system maintains a Global Accumulator, initially set to 1.

```bash
[Client] [Stateful Operation Selected]
Enter Factor to multiply by:
```

After entering a multiplication factor (e.g., Factor = 2), the multiplication result is returned by Server A:

```bash
2026/01/17 20:15:32 [Client] [Request 3] Routing to Server: localhost:9001
2026/01/17 20:15:32 [Client] Global Accumulator: * 2 = 2
```

By executing the multiplication again (still with Factor = 2), the request is handled by Server B, which reads the global accumulator value from the external database
(in this case, Global Accumulator = 2) and performs the multiplication starting from that value:

```bash
2026/01/17 20:19:06 [Client] [Request 4] Routing to Server: localhost:9002
2026/01/17 20:19:06 [Client] Global Accumulator: * 2 = 4
```

This confirms that the accumulator state is shared and preserved across different servers through the external Key-Value Store.

### Scenario C: Deregistration Verification (Shutdown)
To verify the automatic deregistration requirement:

1. Go to the terminal of Server A.
2. Press **CTRL+C**
3. Observe the output on the Registry terminal:

```bash
2026/01/17 20:20:27 [Server] Shutting down...
2026/01/17 20:20:27 [Server] Successfully deregistered from Registry
```

This confirms that servers correctly deregister from the Service Registry upon shutdown.



