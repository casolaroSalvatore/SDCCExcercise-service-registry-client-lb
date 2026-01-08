# service-registry-client-lb

# Progetto 1: Service Registry & Client-Side Load Balancing (Versione Locale)

Questo progetto implementa un sistema distribuito in Go composto da un Service Registry, Server multipli (RPC) e un Client che effettua il bilanciamento del carico (Load Balancing).

## Specifica del Sistema
Il sistema rispetta i seguenti requisiti:
1.  **Service Registry:** Mantiene una lista dinamica degli indirizzi dei server attivi.
2.  **Server RPC:** Offre un servizio di somma; si registra automaticamente all'avvio e si deregistra alla chiusura.
3.  **Client:**
    - **Caching:** Scarica la lista dei server dal Registry una sola volta all'inizio della sessione.
    - **Load Balancing:** Utilizza un algoritmo Round-Robin (Stateless) per distribuire le richieste tra i server.
    - **Stato:** Non utilizza database esterni per la persistenza della lista server.

## Architettura dei File
Il sistema è strutturato nei seguenti componenti:
- **common/types.go**: Strutture dati condivise (necessarie per l'esecuzione del protocollo RPC).
- **registry/main.go**: Il server registry centrale, che gestisce la mappa degli indirizzi.
- **server/main.go**: Il nodo che offre il servizio di calcolo.
- **client/main.go**: Il client che funge da load balancer.

## Guida all'Esecuzione
Poiché l'esecuzione è locale e manuale (non abbiamo virtualizzazione fornita da Docker), è necessario aprire 4 terminali separati.

### 1. Prerequisiti
- Go installato (`go version`).
- Modulo inizializzato: Eseguire `go mod init service-registry-client-lb` nella cartella radice.

### 2. Passi da eseguire

1.  **Terminale 1 (Registry):**
    Aprire un terminale nella cartella radice ed eseguire:

    ```bash
    go run registry/main.go
    ```

    Così facendo avviamo il service registry, che rimane in ascolto sulla porta 8000.
    *Output atteso:*
    `2026/01/07 20:05:28 [Registry] Service Registry running on port :8000`

2.  Aprire ora due (o più) nuovi terminali per simulare due (o più) server diversi. In questo caso utilizziamo 2 server.

    **Terminale 2 (Server A):**
    Eseguire:
    ```bash
    go run server/main.go -port 9001
    ```
    Avvia il primo server di calcolo sulla porta 9001.
    *Output atteso:*
    ```text
    2026/01/07 20:06:48 [Server] Successfully registered to Registry
    2026/01/07 20:06:48 [Server] Server running on localhost:9001. Waiting for requests...
    ```

3.  **Terminale 3 (Server B):**
    Eseguire:
    ```bash
    go run server/main.go -port 9002
    ```
    Avvia il secondo server di calcolo sulla porta 9002.
    *Output atteso:*
    ```text
    2026/01/07 20:07:05 [Server] Successfully registered to Registry
    2026/01/07 20:07:05 [Server] Server running on localhost:9002. Waiting for requests...
    ```

    Sul terminale del **Registry**, vedremo le conferme di registrazione:
    ```text
    2026/01/07 20:06:48 [Registry] Server registered: localhost:9001
    2026/01/07 20:07:05 [Registry] Server registered: localhost:9002
    ```

4.  **Terminale 4 (Client):**
    Aprire ora un ultimo terminale per simulare il client ed eseguire:

    ```bash
    go run client/main.go
    ```

### 3. Risultato
Nel terminale del client si può osservare che le richieste vengono smistate alternativamente tra i server.

*Output atteso:*
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
