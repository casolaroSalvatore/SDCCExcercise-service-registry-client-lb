
```markdown
# Project 1: Service Registry & Client-Side Load Balancing (Local Version)

This project implements a distributed system in Go composed of a Service Registry, multiple RPC Servers, and a Client performing Load Balancing.

## System Specification
The system adheres to the following requirements:
1.  **Service Registry:** Maintains a dynamic list of active server addresses.
2.  **RPC Server:** Provides a summation service; automatically registers at startup and deregisters upon shutdown.
3.  **Client:**
    - **Caching:** Downloads the server list from the Registry only once at the beginning of the session.
    - **Load Balancing:** Uses a Round-Robin (Stateless) algorithm to distribute requests among servers.
    - **State:** Does not use external databases for server list persistence.

## File Architecture
The system is structured into the following components:
- **common/types.go**: Shared data structures (required for RPC protocol execution).
- **registry/main.go**: The central registry server managing the address map.
- **server/main.go**: The node offering the calculation service.
- **client/main.go**: The client acting as the load balancer.

## Execution Guide
Since the execution is local and manual (no virtualization via Docker provided), it is necessary to open 4 separate terminals.

### 1. Prerequisites
- Go installed (`go version`).
- Module initialized: Run `go mod init service-registry-client-lb` in the root folder.

### 2. Steps to Run

1.  **Terminal 1 (Registry):**
    Open a terminal in the root folder and run:

    ```bash
    go run registry/main.go
    ```

    This starts the service registry, listening on port 8000.
    *Expected output:*
    `2026/01/07 20:05:28 [Registry] Service Registry running on port :8000`

2.  Now open two (or more) new terminals to simulate different servers. In this case, we use 2 servers.

    **Terminal 2 (Server A):**
    Run:
    ```bash
    go run server/main.go -port 9001
    ```
    Starts the first calculation server on port 9001.
    *Expected output:*
    ```text
    2026/01/07 20:06:48 [Server] Successfully registered to Registry
    2026/01/07 20:06:48 [Server] Server running on localhost:9001. Waiting for requests...
    ```

3.  **Terminal 3 (Server B):**
    Run:
    ```bash
    go run server/main.go -port 9002
    ```
    Starts the second calculation server on port 9002.
    *Expected output:*
    ```text
    2026/01/07 20:07:05 [Server] Successfully registered to Registry
    2026/01/07 20:07:05 [Server] Server running on localhost:9002. Waiting for requests...
    ```

    On the **Registry** terminal, you will see the successful registration logs:
    ```text
    2026/01/07 20:06:48 [Registry] Server registered: localhost:9001
    2026/01/07 20:07:05 [Registry] Server registered: localhost:9002
    ```

4.  **Terminal 4 (Client):**
    Open a final terminal to simulate the client and run:

    ```bash
    go run client/main.go
    ```

### 3. Result 
In the client terminal, you can observe requests being distributed alternately among the servers.

*Expected output:*
```text
2026/01/07 20:08:36 Found servers: [localhost:9001 localhost:9002]. Starting Load Balancing Session...
2026/01/07 20:08:36 ---------------------------------------------------------
2026/01/07 20:08:36 [Request 1] Selected Server: localhost:9001
2026/01/07 20:08:36  -> Result from localhost:9001: 0 + 0 = 0
2026/01/07 20:08:37 [Request 2] Selected Server: localhost:9002
2026/01/07 20:08:37  -> Result from localhost:9002: 1 + 2 = 3
2026/01/07 20:08:38 [Request 3] Selected Server: localhost:9001
2026/01/07 20:08:38  -> Result from localhost:9001: 2 + 4 = 6
...
2026/01/07 20:08:45 [Request 10] Selected Server: localhost:9002
2026/01/07 20:08:45  -> Result from localhost:9002: 9 + 18 = 27
This confirms that Client-Side Load Balancing (Round-Robin) is working correctly using the cached list.
