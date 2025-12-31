# Socket

Network socket implementations for TCP and UDP protocols.

## Overview

The socket package provides simplified network programming with protocol-specific subpackages:

- **[socket/tcp](tcp/)** - TCP socket client and server
- **[socket/udp](udp/)** - UDP socket client and server

## Choosing a Protocol

| Feature | TCP (socket/tcp) | UDP (socket/udp) |
|---------|------------------|------------------|
| **Connection** | Connection-oriented | Connectionless |
| **Reliability** | Guaranteed delivery | Best effort |
| **Ordering** | In-order delivery | No ordering guarantee |
| **Server** | ✅ Implemented | ✅ Implemented |
| **Client** | ✅ Implemented | ✅ Implemented |
| **Use Cases** | Web servers, databases, file transfer | DNS, gaming, metrics, streaming |

## Installation

```bash
go get -u github.com/common-library/go/socket/tcp
go get -u github.com/common-library/go/socket/udp
```

## Quick Start

### TCP Server

```go
import "github.com/common-library/go/socket/tcp"

server := &tcp.Server{}
err := server.Start("tcp", ":8080", 100,
    func(client tcp.Client) {
        data, _ := client.Read(1024)
        client.Write("Echo: " + data)
    },
    func(err error) {
        log.Printf("Accept error: %v", err)
    },
)
```

### TCP Client

```go
import "github.com/common-library/go/socket/tcp"

client := &tcp.Client{}
client.Connect("tcp", "localhost:8080")
client.Write("Hello")
data, _ := client.Read(1024)
client.Close()
```

### UDP Client

```go
import (
    "time"
    "github.com/common-library/go/socket/udp"
)

client := &udp.Client{}
client.Connect("udp4", "localhost:8080")
client.Send([]byte("Hello"))
data, addr, _ := client.Receive(1024, 5*time.Second)
client.Close()
```

### UDP Server

```go
import (
    "net"
    "github.com/common-library/go/socket/udp"
)

server := &udp.Server{}
err := server.Start("udp4", ":8080", 1024,
    func(data []byte, addr net.Addr, conn net.PacketConn) {
        conn.WriteTo(data, addr)  // Echo back
    },
    true,  // Async handler
    func(err error) {
        log.Printf("Error: %v", err)
    },
)
defer server.Stop()
```

## Features

### TCP (socket/tcp)
- Connection-oriented reliable communication
- Concurrent client handling with goroutine pools
- Automatic resource management
- Graceful shutdown support
- Local and remote address access

### UDP (socket/udp)
- Connectionless packet-based communication
- Packet send/receive with timeout support
- Async/sync packet handler options
- Optional error handling callbacks
- Read/write buffer configuration
- Graceful shutdown with handler completion wait

## Migration from v1.x

Previous versions used a single `socket` package for TCP. For backward compatibility, type aliases are provided:

```go
// Old code (still works with deprecation warnings)
import "github.com/common-library/go/socket"

server := &socket.Server{}  // Deprecated: use tcp.Server
client := &socket.Client{}  // Deprecated: use tcp.Client
```

**Recommended migration:**

```go
// New code
import "github.com/common-library/go/socket/tcp"

server := &tcp.Server{}
client := &tcp.Client{}
```

## Testing

Run all tests:

```bash
go test -v ./...
```

Run protocol-specific tests:

```bash
go test -v ./tcp
go test -v ./udp
```

## License

This package is part of the common-library project and follows the same license.
