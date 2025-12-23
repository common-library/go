# Socket

TCP/UDP socket client and server implementations.

## Overview

The socket package provides simplified network programming with high-level abstractions for socket servers and clients. It handles connection management, concurrent client handling, and resource cleanup automatically.

## Features

- **TCP/UDP Server** - Accept and handle concurrent connections
- **TCP/UDP Client** - Connect, read, and write operations
- **Connection Pooling** - Configurable client pool size
- **Concurrent Handling** - Automatic goroutine management
- **Resource Cleanup** - Graceful shutdown and connection closing

## Installation

```bash
go get -u github.com/common-library/go/socket
```

## Quick Start

### Server

```go
server := &socket.Server{}
err := server.Start("tcp", ":8080", 100,
    func(client socket.Client) {
        data, _ := client.Read(1024)
        client.Write("Echo: " + data)
    },
    func(err error) {
        log.Printf("Accept error: %v", err)
    },
)
```

### Client

```go
client := &socket.Client{}
client.Connect("tcp", "localhost:8080")
client.Write("Hello")
data, _ := client.Read(1024)
client.Close()
```

## API Reference

### Server Type

```go
type Server struct {
    // Internal fields
}
```

### Server Methods

#### Start

```go
func (s *Server) Start(network, address string, clientPoolSize int,
    acceptSuccessFunc func(client Client), 
    acceptFailureFunc func(err error)) error
```

Starts the socket server.

#### Stop

```go
func (s *Server) Stop() error
```

Gracefully stops the server.

#### GetCondition

```go
func (s *Server) GetCondition() bool
```

Returns server running state.

### Client Type

```go
type Client struct {
    // Internal fields
}
```

### Client Methods

#### Connect

```go
func (c *Client) Connect(network, address string) error
```

Connects to a remote address.

#### Read

```go
func (c *Client) Read(recvSize int) (string, error)
```

Reads data from the connection.

#### Write

```go
func (c *Client) Write(data string) (int, error)
```

Writes data to the connection.

#### Close

```go
func (c *Client) Close() error
```

Closes the connection.

#### GetLocalAddr

```go
func (c *Client) GetLocalAddr() net.Addr
```

Returns the local address.

#### GetRemoteAddr

```go
func (c *Client) GetRemoteAddr() net.Addr
```

Returns the remote address.

## Complete Examples

### TCP Echo Server

```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "github.com/common-library/go/socket"
)

func main() {
    server := &socket.Server{}
    
    err := server.Start("tcp", ":8080", 100,
        func(client socket.Client) {
            log.Printf("Client connected: %v", client.GetRemoteAddr())
            
            for {
                data, err := client.Read(1024)
                if err != nil {
                    log.Printf("Read error: %v", err)
                    break
                }
                
                log.Printf("Received: %s", data)
                client.Write("Echo: " + data)
            }
        },
        func(err error) {
            log.Printf("Accept error: %v", err)
        },
    )
    
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Server started on :8080")
    
    // Wait for interrupt
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    
    log.Println("Shutting down...")
    server.Stop()
}
```

### TCP Client

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/socket"
)

func main() {
    client := &socket.Client{}
    
    err := client.Connect("tcp", "localhost:8080")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    fmt.Printf("Connected to %v\n", client.GetRemoteAddr())
    
    // Send message
    _, err = client.Write("Hello, Server!")
    if err != nil {
        log.Fatal(err)
    }
    
    // Receive response
    data, err := client.Read(1024)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Server response: %s\n", data)
}
```

### Chat Server

```go
package main

import (
    "fmt"
    "log"
    "sync"
    "github.com/common-library/go/socket"
)

var (
    clients   = make(map[string]socket.Client)
    clientsMu sync.Mutex
)

func broadcast(message string, sender socket.Client) {
    clientsMu.Lock()
    defer clientsMu.Unlock()
    
    for _, client := range clients {
        if client.GetRemoteAddr() != sender.GetRemoteAddr() {
            client.Write(message)
        }
    }
}

func main() {
    server := &socket.Server{}
    
    server.Start("tcp", ":8080", 100,
        func(client socket.Client) {
            addr := client.GetRemoteAddr().String()
            
            clientsMu.Lock()
            clients[addr] = client
            clientsMu.Unlock()
            
            log.Printf("Client joined: %s", addr)
            broadcast(fmt.Sprintf("%s joined the chat", addr), client)
            
            defer func() {
                clientsMu.Lock()
                delete(clients, addr)
                clientsMu.Unlock()
                
                broadcast(fmt.Sprintf("%s left the chat", addr), client)
                log.Printf("Client left: %s", addr)
            }()
            
            for {
                data, err := client.Read(1024)
                if err != nil {
                    break
                }
                
                message := fmt.Sprintf("%s: %s", addr, data)
                log.Println(message)
                broadcast(message, client)
            }
        },
        nil,
    )
    
    log.Println("Chat server started on :8080")
    select {}
}
```

## Best Practices

### 1. Set Appropriate Pool Size

```go
// Good: Based on expected load
server.Start("tcp", ":8080", 100, ...) // 100 concurrent connections

// Avoid: Too small
server.Start("tcp", ":8080", 5, ...)   // May block under load
```

### 2. Handle Read/Write Errors

```go
// Good: Check errors
data, err := client.Read(1024)
if err != nil {
    log.Printf("Read error: %v", err)
    return
}

// Avoid: Ignore errors
data, _ := client.Read(1024)
```

### 3. Always Close Connections

```go
// Good: Defer close
client := &socket.Client{}
err := client.Connect("tcp", "localhost:8080")
if err == nil {
    defer client.Close()
}

// Avoid: Forget to close
client.Connect("tcp", "localhost:8080")
// No close - connection leak
```

## Dependencies

- `net` - Go standard library
- `sync` - Go standard library
- `sync/atomic` - Go standard library

## Further Reading

- [Go net Package](https://pkg.go.dev/net)
- [Network Programming with Go](https://tumregels.github.io/Network-Programming-with-Go/)
