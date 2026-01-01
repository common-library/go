# TCP Socket Package

Package providing TCP socket client and server implementation for Go.

## Features

- **TCP Socket Client**: Simple connect, read, write operations
- **TCP Socket Server**: Concurrent connection handling and automatic resource management
- **Concurrency Support**: Handle multiple client connections simultaneously
- **Automatic Resource Cleanup**: Automatic connection and resource release
- **Address Information Access**: Query local and remote addresses

## Installation

```bash
go get github.com/common-library/go/socket/tcp
```

## Usage Examples

### TCP Client

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/socket/tcp"
)

func main() {
    client := &tcp.Client{}
    
    // Connect to server
    err := client.Connect("tcp", "localhost:8080")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // Write data
    n, err := client.Write("Hello, Server!")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Sent %d bytes\n", n)
    
    // Read data
    data, err := client.Read(1024)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Received: %s\n", data)
    
    // Address information
    fmt.Printf("Local: %s\n", client.GetLocalAddr())
    fmt.Printf("Remote: %s\n", client.GetRemoteAddr())
}
```

### TCP Server

```go
package main

import (
    "log"
    "github.com/common-library/go/socket/tcp"
)

func main() {
    server := &tcp.Server{}
    
    // Connection handler
    acceptSuccessFunc := func(client tcp.Client) {
        // Read data from client
        data, err := client.Read(1024)
        if err != nil {
            log.Printf("Read error: %v", err)
            return
        }
        
        // Echo response
        _, err = client.Write("Echo: " + data)
        if err != nil {
            log.Printf("Write error: %v", err)
        }
    }
    
    // Connection failure handler
    acceptFailureFunc := func(err error) {
        log.Printf("Accept error: %v", err)
    }
    
    // Start server
    err := server.Start("tcp", ":8080", 100, acceptSuccessFunc, acceptFailureFunc)
    if err != nil {
        log.Fatal(err)
    }
    
    // Wait until server stops
    // (In reality, signal handling etc. should be added)
    select {}
}
```

### Multi-Client Handling Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/socket/tcp"
)

func main() {
    server := &tcp.Server{}
    
    acceptSuccessFunc := func(client tcp.Client) {
        defer client.Close()
        
        // Send welcome message
        client.Write("Welcome to the server!\n")
        
        // Process client requests
        for {
            data, err := client.Read(1024)
            if err != nil {
                break
            }
            
            // Simple protocol processing
            response := processRequest(data)
            client.Write(response)
        }
    }
    
    acceptFailureFunc := func(err error) {
        log.Printf("Connection error: %v", err)
    }
    
    // Handle maximum 1000 concurrent connections
    err := server.Start("tcp", ":8080", 1000, acceptSuccessFunc, acceptFailureFunc)
    if err != nil {
        log.Fatal(err)
    }
    defer server.Stop()
    
    fmt.Println("Server is running on :8080")
    select {}
}

func processRequest(data string) string {
    return fmt.Sprintf("Processed: %s\n", data)
}
```

## API Documentation

### Client

#### `Connect(network, address string) error`
Establishes a connection to a remote address.

- **network**: Network type ("tcp", "tcp4", "tcp6", "unix")
- **address**: Remote address (e.g., "localhost:8080", "192.168.1.1:9000")

#### `Read(recvSize int) (string, error)`
Reads data from the connection.

- **recvSize**: Maximum number of bytes to read (buffer size)
- **return**: Received data string, error

#### `Write(data string) (int, error)`
Writes data to the connection.

- **data**: Text data to write
- **return**: Number of bytes written, error

#### `Close() error`
Closes the connection.

#### `GetLocalAddr() net.Addr`
Returns the local network address.

#### `GetRemoteAddr() net.Addr`
Returns the remote network address.

### Server

#### `Start(network, address string, clientPoolSize int, acceptSuccessFunc func(client Client), acceptFailureFunc func(err error)) error`
Initializes and starts the socket server.

- **network**: Network type ("tcp", "tcp4", "tcp6", "unix")
- **address**: Listen address (e.g., ":8080", "127.0.0.1:8080")
- **clientPoolSize**: Maximum number of concurrent connections to buffer
- **acceptSuccessFunc**: Callback function for each connection
- **acceptFailureFunc**: Callback function for Accept errors

#### `Stop() error`
Safely shuts down the socket server.

#### `GetCondition() bool`
Returns the server running status.

- **return**: true if server is running, false if stopped

## Notes

- Client `Read()` and `Write()` must be used after calling `Connect()`.
- Server `acceptSuccessFunc` runs in a goroutine, so concurrency must be considered.
- `clientPoolSize` determines the buffer size for connections that can be handled simultaneously.
- When stopping the server, `Stop()` waits until all active connections are terminated.

## License

This package follows the project's license.
