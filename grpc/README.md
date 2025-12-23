# gRPC

Simplified utilities for gRPC client and server in Go.

## Overview

The grpc package provides convenient wrapper functions around Google's gRPC-Go library, simplifying the creation of client connections and server lifecycle management. It reduces boilerplate code while maintaining the full power of gRPC.

## Features

- **Simplified Client Connections** - Easy client connection creation
- **Server Lifecycle Management** - Start/Stop server with simple API
- **Thread-Safe Operations** - Safe concurrent server management
- **Minimal Configuration** - Sensible defaults for quick setup
- **Full gRPC Compatibility** - Works with standard gRPC services

## Installation

```bash
go get -u github.com/common-library/go/grpc
go get -u google.golang.org/grpc
go get -u google.golang.org/protobuf
```

## Quick Start

### 1. Define Protocol Buffers

```protobuf
// greeter.proto
syntax = "proto3";

package greeter;
option go_package = "example.com/greeter/pb";

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

### 2. Generate Go Code

```bash
protoc --go_out=. --go-grpc_out=. greeter.proto
```

### 3. Implement Server

```go
package main

import (
    "context"
    "log"
    
    "github.com/common-library/go/grpc"
    pb "example.com/greeter/pb"
    grpclib "google.golang.org/grpc"
)

type GreeterServer struct {
    pb.UnimplementedGreeterServer
}

func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{
        Message: "Hello, " + req.Name,
    }, nil
}

func (s *GreeterServer) RegisterServer(server *grpclib.Server) {
    pb.RegisterGreeterServer(server, s)
}

func main() {
    var server grpc.Server
    
    log.Println("Starting gRPC server on :50051")
    err := server.Start(":50051", &GreeterServer{})
    if err != nil {
        log.Fatal(err)
    }
}
```

### 4. Implement Client

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/common-library/go/grpc"
    pb "example.com/greeter/pb"
)

func main() {
    conn, err := grpc.GetConnection("localhost:50051")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    client := pb.NewGreeterClient(conn)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    response, err := client.SayHello(ctx, &pb.HelloRequest{
        Name: "Alice",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Response: %s", response.Message)
}
```

## API Reference

### Client

#### GetConnection

```go
func GetConnection(address string) (*grpc.ClientConn, error)
```

Creates a new gRPC client connection.

**Parameters:**
- `address` - Server address in "host:port" format

**Returns:**
- `*grpc.ClientConn` - Client connection for creating service clients
- `error` - Error if connection fails

**Example:**

```go
conn, err := grpc.GetConnection("localhost:50051")
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := pb.NewMyServiceClient(conn)
```

### Server

#### Server Type

```go
type Server struct {
    // Contains unexported fields
}
```

Server manages the gRPC server lifecycle with thread-safe operations.

#### implementServer Interface

```go
type implementServer interface {
    RegisterServer(server *grpc.Server)
}
```

Interface that service implementations must satisfy. The RegisterServer method should register the gRPC service handlers.

#### Start

```go
func (grpcSrv *Server) Start(address string, server implementServer) error
```

Starts the gRPC server on the specified address.

**Parameters:**
- `address` - Address to bind (e.g., ":50051")
- `server` - Service implementation with RegisterServer method

**Returns:**
- `error` - Error if server fails to start or during serving

**Example:**

```go
var server grpc.Server

err := server.Start(":50051", &MyService{})
if err != nil {
    log.Fatal(err)
}
```

#### Stop

```go
func (grpcSrv *Server) Stop() error
```

Gracefully stops the gRPC server.

**Returns:**
- `error` - Always returns nil

**Example:**

```go
err := server.Stop()
if err != nil {
    log.Printf("Error: %v", err)
}
```

## Complete Examples

### Unary RPC

```go
// Server
package main

import (
    "context"
    "log"
    
    "github.com/common-library/go/grpc"
    pb "example.com/calculator/pb"
    grpclib "google.golang.org/grpc"
)

type CalculatorServer struct {
    pb.UnimplementedCalculatorServer
}

func (s *CalculatorServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddReply, error) {
    result := req.A + req.B
    return &pb.AddReply{Result: result}, nil
}

func (s *CalculatorServer) RegisterServer(server *grpclib.Server) {
    pb.RegisterCalculatorServer(server, s)
}

func main() {
    var server grpc.Server
    log.Fatal(server.Start(":50051", &CalculatorServer{}))
}
```

```go
// Client
package main

import (
    "context"
    "log"
    
    "github.com/common-library/go/grpc"
    pb "example.com/calculator/pb"
)

func main() {
    conn, _ := grpc.GetConnection("localhost:50051")
    defer conn.Close()
    
    client := pb.NewCalculatorClient(conn)
    
    response, err := client.Add(context.Background(), &pb.AddRequest{
        A: 10,
        B: 20,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("10 + 20 = %d", response.Result)
}
```

### Server Streaming RPC

```go
// Server
func (s *Server) ListItems(req *pb.ListRequest, stream pb.Service_ListItemsServer) error {
    items := []string{"item1", "item2", "item3", "item4", "item5"}
    
    for _, item := range items {
        if err := stream.Send(&pb.Item{Name: item}); err != nil {
            return err
        }
        time.Sleep(time.Second)
    }
    
    return nil
}
```

```go
// Client
stream, err := client.ListItems(context.Background(), &pb.ListRequest{})
if err != nil {
    log.Fatal(err)
}

for {
    item, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Received: %s", item.Name)
}
```

### Client Streaming RPC

```go
// Server
func (s *Server) UploadData(stream pb.Service_UploadDataServer) error {
    var total int32
    
    for {
        data, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.UploadResponse{
                TotalItems: total,
            })
        }
        if err != nil {
            return err
        }
        
        total++
        log.Printf("Received: %s", data.Content)
    }
}
```

```go
// Client
stream, err := client.UploadData(context.Background())
if err != nil {
    log.Fatal(err)
}

items := []string{"data1", "data2", "data3"}
for _, item := range items {
    err := stream.Send(&pb.Data{Content: item})
    if err != nil {
        log.Fatal(err)
    }
}

response, err := stream.CloseAndRecv()
if err != nil {
    log.Fatal(err)
}

log.Printf("Uploaded %d items", response.TotalItems)
```

### Bidirectional Streaming RPC

```go
// Server
func (s *Server) Chat(stream pb.Service_ChatServer) error {
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        
        reply := &pb.Message{
            Text: "Echo: " + msg.Text,
        }
        
        if err := stream.Send(reply); err != nil {
            return err
        }
    }
}
```

```go
// Client
stream, err := client.Chat(context.Background())
if err != nil {
    log.Fatal(err)
}

// Send messages
go func() {
    messages := []string{"Hello", "How are you?", "Goodbye"}
    for _, msg := range messages {
        stream.Send(&pb.Message{Text: msg})
        time.Sleep(time.Second)
    }
    stream.CloseSend()
}()

// Receive messages
for {
    msg, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Received: %s", msg.Text)
}
```

### Server with Graceful Shutdown

```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/common-library/go/grpc"
)

func main() {
    var server grpc.Server
    
    // Start server in goroutine
    go func() {
        log.Println("Starting gRPC server on :50051")
        if err := server.Start(":50051", &MyService{}); err != nil {
            log.Fatal(err)
        }
    }()
    
    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    
    log.Println("Shutting down server...")
    if err := server.Stop(); err != nil {
        log.Printf("Error stopping server: %v", err)
    }
    
    log.Println("Server stopped gracefully")
}
```

### Multiple Services on One Server

```go
type CombinedServer struct {
    pb.UnimplementedGreeterServer
    pb.UnimplementedCalculatorServer
}

func (s *CombinedServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{Message: "Hello, " + req.Name}, nil
}

func (s *CombinedServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddReply, error) {
    return &pb.AddReply{Result: req.A + req.B}, nil
}

func (s *CombinedServer) RegisterServer(server *grpc.Server) {
    pb.RegisterGreeterServer(server, s)
    pb.RegisterCalculatorServer(server, s)
}

func main() {
    var server grpc.Server
    server.Start(":50051", &CombinedServer{})
}
```

### Client with Connection Pooling

```go
type ClientPool struct {
    connections []*grpc.ClientConn
    index       int
    mutex       sync.Mutex
}

func NewClientPool(address string, size int) (*ClientPool, error) {
    pool := &ClientPool{
        connections: make([]*grpc.ClientConn, size),
    }
    
    for i := 0; i < size; i++ {
        conn, err := grpc.GetConnection(address)
        if err != nil {
            return nil, err
        }
        pool.connections[i] = conn
    }
    
    return pool, nil
}

func (p *ClientPool) GetConnection() *grpc.ClientConn {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    conn := p.connections[p.index]
    p.index = (p.index + 1) % len(p.connections)
    
    return conn
}

func (p *ClientPool) Close() {
    for _, conn := range p.connections {
        conn.Close()
    }
}
```

## Best Practices

### 1. Always Close Client Connections

```go
// Good: Defer close
conn, err := grpc.GetConnection("localhost:50051")
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// Avoid: Forgetting to close
conn, _ := grpc.GetConnection("localhost:50051")
// Connection leak!
```

### 2. Use Context for Timeouts

```go
// Good: Set timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

response, err := client.SayHello(ctx, request)

// Avoid: No timeout
response, err := client.SayHello(context.Background(), request)
// May hang indefinitely
```

### 3. Handle Graceful Shutdown

```go
// Good: Graceful shutdown
go server.Start(":50051", service)

sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
<-sigChan

server.Stop() // Wait for active RPCs

// Avoid: Abrupt termination
// Just exiting without calling Stop()
```

### 4. Implement Error Handling

```go
// Good: Check and handle errors
response, err := client.Call(ctx, request)
if err != nil {
    if stat, ok := status.FromError(err); ok {
        switch stat.Code() {
        case codes.NotFound:
            log.Println("Resource not found")
        case codes.DeadlineExceeded:
            log.Println("Request timeout")
        default:
            log.Printf("RPC error: %v", err)
        }
    }
    return err
}

// Avoid: Ignoring errors
client.Call(ctx, request) // Error ignored
```

### 5. Use Secure Connections in Production

```go
// Development: Insecure (current implementation)
conn, _ := grpc.GetConnection("localhost:50051")

// Production: Use TLS
creds, _ := credentials.NewClientTLSFromFile("cert.pem", "")
conn, _ := grpc.Dial("localhost:50051", 
    grpc.WithTransportCredentials(creds))
```

### 6. Implement Health Checks

```go
import "google.golang.org/grpc/health/grpc_health_v1"

type HealthServer struct {
    grpc_health_v1.UnimplementedHealthServer
}

func (s *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
    return &grpc_health_v1.HealthCheckResponse{
        Status: grpc_health_v1.HealthCheckResponse_SERVING,
    }, nil
}

func (s *HealthServer) RegisterServer(server *grpc.Server) {
    grpc_health_v1.RegisterHealthServer(server, s)
}
```

## Error Handling

### Common gRPC Errors

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// Server-side: Return errors
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    user, err := database.GetUser(req.Id)
    if err == sql.ErrNoRows {
        return nil, status.Errorf(codes.NotFound, "user not found: %d", req.Id)
    }
    if err != nil {
        return nil, status.Errorf(codes.Internal, "database error: %v", err)
    }
    return user, nil
}

// Client-side: Handle errors
response, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 123})
if err != nil {
    stat, ok := status.FromError(err)
    if !ok {
        log.Fatal("Unknown error type")
    }
    
    switch stat.Code() {
    case codes.NotFound:
        log.Println("User not found")
    case codes.InvalidArgument:
        log.Println("Invalid request")
    case codes.Unavailable:
        log.Println("Service unavailable")
    default:
        log.Printf("Error: %v", stat.Message())
    }
}
```

## Performance Tips

1. **Connection Reuse** - Reuse client connections instead of creating new ones
2. **Streaming** - Use streaming RPCs for large data transfers
3. **Compression** - Enable compression for large messages
4. **Keep-Alive** - Configure keep-alive to detect dead connections
5. **Connection Pooling** - Use multiple connections for high throughput

## Testing

### Unit Testing with Mock

```go
func TestGreeterServer(t *testing.T) {
    server := &GreeterServer{}
    
    req := &pb.HelloRequest{Name: "Alice"}
    resp, err := server.SayHello(context.Background(), req)
    
    if err != nil {
        t.Fatalf("SayHello failed: %v", err)
    }
    
    expected := "Hello, Alice"
    if resp.Message != expected {
        t.Errorf("Expected %s, got %s", expected, resp.Message)
    }
}
```

### Integration Testing

```go
func TestClientServer(t *testing.T) {
    // Start server
    var server grpc.Server
    go server.Start(":50052", &GreeterServer{})
    defer server.Stop()
    
    time.Sleep(100 * time.Millisecond) // Wait for server
    
    // Test client
    conn, err := grpc.GetConnection("localhost:50052")
    if err != nil {
        t.Fatal(err)
    }
    defer conn.Close()
    
    client := pb.NewGreeterClient(conn)
    resp, err := client.SayHello(context.Background(), &pb.HelloRequest{
        Name: "Test",
    })
    
    if err != nil {
        t.Fatal(err)
    }
    
    if resp.Message != "Hello, Test" {
        t.Errorf("Unexpected response: %s", resp.Message)
    }
}
```

## Limitations

1. **Insecure Connections Only** - Current implementation only supports insecure connections
2. **No TLS Support** - For production, extend with TLS credentials
3. **Basic Server Management** - No advanced features like reflection or health checks
4. **No Interceptors** - No built-in support for interceptors/middleware

## Migration to Secure Connections

To add TLS support:

```go
import "google.golang.org/grpc/credentials"

// Client
creds, err := credentials.NewClientTLSFromFile("server.crt", "")
conn, err := grpc.Dial("localhost:50051",
    grpc.WithTransportCredentials(creds))

// Server
creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
s := grpc.NewServer(grpc.Creds(creds))
```

## Dependencies

- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/protobuf` - Protocol Buffers

## Further Reading

- [gRPC Official Documentation](https://grpc.io/docs/)
- [gRPC-Go Documentation](https://pkg.go.dev/google.golang.org/grpc)
- [Protocol Buffers Guide](https://protobuf.dev/)
- [gRPC Best Practices](https://grpc.io/docs/guides/performance/)
- [gRPC Error Handling](https://grpc.io/docs/guides/error/)
