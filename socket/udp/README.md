# UDP Socket Package

Package providing UDP socket client and server implementation for Go.

## Features

- **UDP Socket Client**: Packet send and receive operations
- **UDP Socket Server**: Processing through packet handlers
- **Concurrency Support**: Asynchronous handler option
- **Timeout Support**: Receive timeout setting
- **Error Handling**: Optional error handler
- **Buffer Management**: Read/write buffer size settings

## Installation

```bash
go get github.com/common-library/go/socket/udp
```

## UDP Characteristics

UDP (User Datagram Protocol) is a **connectionless protocol**:
- No packet delivery guarantee (packets may be lost)
- No order guarantee (packets may not arrive in order)
- No duplication prevention (same packet may arrive multiple times)
- Low overhead and fast transmission speed
- Suitable for real-time applications (streaming, gaming, DNS, etc.)

## Usage Examples

### UDP Client

#### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/common-library/go/socket/udp"
)

func main() {
    client := &udp.Client{}
    
    // Connect to server
    err := client.Connect("udp4", "localhost:8080")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // Send data
    data := []byte("Hello, UDP Server!")
    n, err := client.Send(data)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Sent %d bytes\n", n)
    
    // Receive data (5 second timeout)
    received, addr, err := client.Receive(1024, 5*time.Second)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Received from %s: %s\n", addr, received)
}
```

#### Send to Different Address (SendTo)

```go
client := &udp.Client{}
client.Connect("udp4", "localhost:8080")
defer client.Close()

// Send to different address than Connect
n, err := client.SendTo([]byte("Hello"), "192.168.1.100:9000")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Sent %d bytes to different address\n", n)
```

#### Buffer Size Settings

```go
client := &udp.Client{}
client.Connect("udp4", "localhost:8080")
defer client.Close()

// Set read/write buffer sizes
client.SetReadBuffer(65536)   // 64KB
client.SetWriteBuffer(65536)  // 64KB
```

### UDP Server

#### Basic Echo Server

```go
package main

import (
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"
    "github.com/common-library/go/socket/udp"
)

func main() {
    server := &udp.Server{}
    
    // Start echo server
    err := server.Start("udp4", ":8080", 1024,
        func(data []byte, addr net.Addr, conn net.PacketConn) {
            fmt.Printf("Received from %s: %s\n", addr, data)
            // Echo back
            conn.WriteTo(data, addr)
        },
        false,  // Synchronous processing
        nil,    // No error handler
    )
    if err != nil {
        log.Fatal(err)
    }
    defer server.Stop()
    
    fmt.Printf("Server listening on %s\n", server.GetLocalAddr())
    
    // Wait for signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    
    fmt.Println("Shutting down...")
}
```

#### Asynchronous Handler with Error Processing

```go
server := &udp.Server{}

err := server.Start("udp4", ":8080", 1024,
    // Packet handler
    func(data []byte, addr net.Addr, conn net.PacketConn) {
        // Process each packet in a separate goroutine
        processPacket(data, addr)
        
        // Send response
        response := []byte("Acknowledged")
        conn.WriteTo(response, addr)
    },
    true,  // Asynchronous processing (use goroutines)
    // Error handler
    func(err error) {
        log.Printf("Read error: %v", err)
    },
)
if err != nil {
    log.Fatal(err)
}
defer server.Stop()
```

#### Multicast Server Example

```go
package main

import (
    "fmt"
    "log"
    "net"
    "github.com/common-library/go/socket/udp"
)

func main() {
    // Multicast address
    multicastAddr := "224.0.0.1:9999"
    
    server := &udp.Server{}
    
    err := server.Start("udp4", multicastAddr, 1024,
        func(data []byte, addr net.Addr, conn net.PacketConn) {
            fmt.Printf("Multicast from %s: %s\n", addr, data)
        },
        true,
        func(err error) {
            log.Printf("Error: %v", err)
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    defer server.Stop()
    
    select {} // Keep running
}
```

#### Status Check

```go
server := &udp.Server{}
server.Start("udp4", ":8080", 1024, handler, false, nil)

// Check server running status
if server.IsRunning() {
    fmt.Println("Server is running")
    fmt.Printf("Listening on: %s\n", server.GetLocalAddr())
}

// Stop server
server.Stop()

if !server.IsRunning() {
    fmt.Println("Server stopped")
}
```

## API Documentation

### Client

#### `Connect(network, address string) error`
Creates a UDP connection.

- **network**: Network type ("udp", "udp4", "udp6")
- **address**: Remote address (e.g., "localhost:8080")

#### `Send(data []byte) (int, error)`
Sends data to the connected address.

- **data**: Byte array to send
- **return**: Number of bytes sent, error

#### `SendTo(data []byte, address string) (int, error)`
Sends data to a specific address.

- **data**: Byte array to send
- **address**: Destination address
- **return**: Number of bytes sent, error

**Note**: Creates a temporary socket for each call, which may cause performance degradation with frequent use.

#### `Receive(bufferSize int, timeout time.Duration) ([]byte, net.Addr, error)`
Receives data.

- **bufferSize**: Maximum number of bytes to read
- **timeout**: Read timeout (no timeout if 0)
- **return**: Received data, sender address, error

#### `Close() error`
Closes the UDP connection.

#### `GetLocalAddr() net.Addr`
Returns the local network address.

#### `GetRemoteAddr() net.Addr`
Returns the remote network address.

#### `SetReadBuffer(bytes int) error`
Sets the OS receive buffer size.

#### `SetWriteBuffer(bytes int) error`
Sets the OS send buffer size.

### Server

#### `Start(network, address string, bufferSize int, handler PacketHandler, asyncHandler bool, errorHandler ErrorHandler) error`
Starts the UDP server.

- **network**: Network type ("udp", "udp4", "udp6")
- **address**: Bind address (e.g., ":8080", "0.0.0.0:8080")
- **bufferSize**: Maximum size of received packets (bytes)
- **handler**: Packet processing function
- **asyncHandler**: If true, each handler runs in a goroutine
- **errorHandler**: Read error processing function (can be nil)

#### `Stop() error`
Stops the UDP server. Waits until all running handlers are complete.

#### `IsRunning() bool`
Returns the server running status.

- **return**: true if server is running

#### `GetLocalAddr() net.Addr`
Returns the local address the server is listening on.

## Type Definitions

### `PacketHandler`
```go
type PacketHandler func(data []byte, addr net.Addr, conn net.PacketConn)
```
Function type that processes received UDP packets.

### `ErrorHandler`
```go
type ErrorHandler func(err error)
```
Function type that handles errors occurring during packet reception.

## Notes

### UDP Characteristics Related
- UDP does not guarantee packet delivery, order, or duplication prevention
- For critical data, implement acknowledgment and retransmission mechanisms at the application level
- Packets exceeding network MTU may be fragmented or lost

### Client
- `Send()` and `Receive()` must be used after calling `Connect()`
- `SendTo()` creates a temporary socket, so be mindful of performance
- If timeout occurs, an `i/o timeout` error is returned

### Server
- When `asyncHandler=true`, handlers run concurrently, so concurrency must be considered
- When `asyncHandler=false`, if handlers take a long time, next packet reception will be delayed
- `Stop()` blocks until all handlers are complete
- Use error handler for logging or monitoring network errors

## Performance Tips

1. **Buffer size**: Set buffer according to expected maximum packet size
2. **Asynchronous processing**: Use `asyncHandler=true` in high-load environments
3. **OS buffers**: Utilize `SetReadBuffer()`/`SetWriteBuffer()` for handling large volumes of packets
4. **Timeout**: Set appropriate timeout to prevent infinite waiting

## License

This package follows the project's license.
