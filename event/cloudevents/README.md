# CloudEvents

CloudEvents client and server implementation for event-driven architectures.

## Overview

The cloudevents package provides a simplified wrapper around the official CloudEvents SDK for Go. It enables sending, receiving, and processing CloudEvents over HTTP with support for both one-way messaging and request-response patterns.

CloudEvents is a specification for describing event data in a common way, promoting interoperability across services, platforms, and systems.

## Features

- **HTTP Client & Server** - Send and receive CloudEvents over HTTP
- **Request-Response Pattern** - Synchronous request-response event flow
- **Asynchronous Receiver** - Background event receiver with lifecycle management
- **Result Types** - Detailed delivery status (ACK/NACK/Undelivered)
- **CloudEvents v1.0 Compliant** - Full specification compliance
- **Flexible Configuration** - Customizable HTTP options and client settings

## Installation

```bash
go get -u github.com/common-library/go/event/cloudevents
go get -u github.com/cloudevents/sdk-go/v2
```

## Quick Start

### Sending Events

```go
import "github.com/common-library/go/event/cloudevents"

func main() {
    // Create client
    client, err := cloudevents.NewHttp("http://localhost:8080", nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create event
    event := cloudevents.NewEvent()
    event.SetType("com.example.user.created")
    event.SetSource("example/users")
    event.SetData("application/json", map[string]string{
        "user_id": "123",
        "name":    "Alice",
    })
    
    // Send event
    result := client.Send(event)
    if result.IsUndelivered() {
        log.Printf("Failed to send: %s", result.Error())
    } else {
        log.Println("Event sent successfully")
    }
}
```

### Receiving Events

```go
import (
    "context"
    "github.com/common-library/go/event/cloudevents"
    "github.com/cloudevents/sdk-go/v2/protocol/http"
)

func main() {
    // Create receiver with port configuration
    httpOpts := []http.Option{http.WithPort(8080)}
    receiver, err := cloudevents.NewHttp("", httpOpts, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Define event handler
    handler := func(ctx context.Context, event cloudevents.Event) {
        log.Printf("Received event: %s", event.Type())
        log.Printf("Data: %v", event.Data())
    }
    
    // Define failure handler
    failureFunc := func(err error) {
        log.Printf("Receiver error: %v", err)
    }
    
    // Start receiver
    receiver.StartReceiver(handler, failureFunc)
    defer receiver.StopReceiver()
    
    // Wait for shutdown signal
    // ...
}
```

### Request-Response Pattern

```go
// Client side - send request and wait for response
event := cloudevents.NewEvent()
event.SetType("com.example.query")
event.SetSource("example/client")
event.SetData("application/json", map[string]string{"query": "status"})

response, result := client.Request(event)
if result.IsUndelivered() {
    log.Printf("Request failed: %s", result.Error())
} else {
    log.Printf("Response: %v", response.Data())
}
```

### Running a Server

```go
import "github.com/common-library/go/event/cloudevents"

func main() {
    var server cloudevents.Server
    
    // Define request handler
    handler := func(event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
        log.Printf("Received: %s", event.Type())
        
        // Process event...
        
        // Return response event
        response := cloudevents.NewEvent()
        response.SetType("com.example.response")
        response.SetSource("example/server")
        response.SetData("application/json", map[string]string{
            "status": "processed",
        })
        
        return &response, cloudevents.NewHTTPResult(200, "OK")
    }
    
    failureFunc := func(err error) {
        log.Printf("Server error: %v", err)
    }
    
    // Start server
    err := server.Start(":8080", handler, failureFunc)
    if err != nil {
        log.Fatal(err)
    }
    
    // ... wait for shutdown signal ...
    
    // Graceful shutdown
    server.Stop(10 * time.Second)
}
```

## CloudEvents Structure

A CloudEvent has the following required attributes:

```go
event := cloudevents.NewEvent()

// Required attributes
event.SetID("abc-123")                    // Unique identifier
event.SetType("com.example.object.action") // Event type
event.SetSource("example/source")          // Event source

// Optional attributes
event.SetSubject("user/123")              // Subject within source
event.SetTime(time.Now())                 // Timestamp
event.SetDataContentType("application/json") // Content type

// Event data
event.SetData("application/json", myData)

// Access attributes
id := event.ID()
eventType := event.Type()
source := event.Source()
data := event.Data()
```

## Result Handling

### Result Types

```go
result := client.Send(event)

// Check delivery status
if result.IsACK() {
    // Event acknowledged (success)
    log.Println("Event delivered and acknowledged")
}

if result.IsNACK() {
    // Event not acknowledged (rejected)
    log.Printf("Event rejected: %s", result.Error())
}

if result.IsUndelivered() {
    // Event could not be delivered (network/connection error)
    log.Printf("Delivery failed: %s", result.Error())
}

// Get HTTP status code (for HTTP transport)
statusCode, err := result.GetHttpStatusCode()
if err == nil {
    log.Printf("HTTP Status: %d", statusCode)
}
```

### Creating Results

```go
// Generic result
result := cloudevents.NewResult("event processed")

// HTTP result with status code
result = cloudevents.NewHTTPResult(200, "OK")
result = cloudevents.NewHTTPResult(400, "invalid event type: %s", eventType)
result = cloudevents.NewHTTPResult(500, "processing error: %v", err)
```

## Complete Examples

### Microservice Event Producer

```go
package main

import (
    "log"
    "time"
    
    "github.com/common-library/go/event/cloudevents"
)

func main() {
    client, err := cloudevents.NewHttp("http://event-gateway:8080", nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Publish user creation event
    event := cloudevents.NewEvent()
    event.SetType("com.example.user.created")
    event.SetSource("users-service")
    event.SetData("application/json", map[string]interface{}{
        "user_id":   "12345",
        "email":     "alice@example.com",
        "created_at": time.Now(),
    })
    
    result := client.Send(event)
    if result.IsUndelivered() {
        log.Fatal(result.Error())
    }
    
    statusCode, _ := result.GetHttpStatusCode()
    log.Printf("Event sent with status: %d", statusCode)
}
```

### Event Consumer Service

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/common-library/go/event/cloudevents"
    "github.com/cloudevents/sdk-go/v2/protocol/http"
)

func main() {
    httpOpts := []http.Option{http.WithPort(8080)}
    receiver, err := cloudevents.NewHttp("", httpOpts, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    handler := func(ctx context.Context, event cloudevents.Event) {
        switch event.Type() {
        case "com.example.user.created":
            handleUserCreated(event)
        case "com.example.user.updated":
            handleUserUpdated(event)
        default:
            log.Printf("Unknown event type: %s", event.Type())
        }
    }
    
    failureFunc := func(err error) {
        log.Fatalf("Receiver error: %v", err)
    }
    
    receiver.StartReceiver(handler, failureFunc)
    defer receiver.StopReceiver()
    
    log.Println("Event receiver started on :8080")
    
    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    
    log.Println("Shutting down...")
}

func handleUserCreated(event cloudevents.Event) {
    log.Printf("User created: %v", event.Data())
    // Process event...
}

func handleUserUpdated(event cloudevents.Event) {
    log.Printf("User updated: %v", event.Data())
    // Process event...
}
```

### Request-Response Service

```go
package main

import (
    "log"
    "time"
    
    "github.com/common-library/go/event/cloudevents"
)

func main() {
    var server cloudevents.Server
    
    handler := func(event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
        log.Printf("Request: %s", event.Type())
        
        switch event.Type() {
        case "com.example.query.status":
            // Process query and create response
            response := cloudevents.NewEvent()
            response.SetType("com.example.response.status")
            response.SetSource("status-service")
            response.SetData("application/json", map[string]string{
                "status": "healthy",
                "uptime": "24h",
            })
            return &response, cloudevents.NewHTTPResult(200, "OK")
            
        default:
            return nil, cloudevents.NewHTTPResult(400, "unknown query type")
        }
    }
    
    failureFunc := func(err error) {
        log.Printf("Server error: %v", err)
    }
    
    err := server.Start(":8080", handler, failureFunc)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Request-response server started on :8080")
    
    // Keep running...
    select {}
}
```

## API Reference

### Client Creation

#### `NewHttp(address string, httpOption []http.Option, clientOption []client.Option) (*client, error)`

Creates an HTTP CloudEvents client.

**Parameters:**
- `address` - Target URL for sending events (e.g., "http://localhost:8080")
- `httpOption` - HTTP protocol options (WithPort, WithPath, etc.)
- `clientOption` - Client configuration options

**Returns:** Client instance and error

### Client Methods

#### `Send(event Event) Result`

Sends a CloudEvent (one-way).

#### `Request(event Event) (*Event, Result)`

Sends a CloudEvent and waits for response.

#### `StartReceiver(handler func(context.Context, Event), failureFunc func(error))`

Starts receiving events asynchronously.

#### `StopReceiver()`

Stops the event receiver gracefully.

### Server Methods

#### `Start(address string, handler func(Event) (*Event, Result), listenAndServeFailureFunc func(error)) error`

Starts the CloudEvents HTTP server.

#### `Stop(shutdownTimeout time.Duration) error`

Stops the server gracefully.

### Result Methods

#### `NewResult(format string, arguments ...any) Result`

Creates a generic result.

#### `NewHTTPResult(statusCode int, format string, arguments ...any) Result`

Creates an HTTP result with status code.

#### `IsACK() bool`

Checks if event was acknowledged.

#### `IsNACK() bool`

Checks if event was not acknowledged.

#### `IsUndelivered() bool`

Checks if event could not be delivered.

#### `GetHttpStatusCode() (int, error)`

Gets HTTP status code from result.

#### `Error() string`

Gets error message.

### Event Methods

#### `NewEvent() Event`

Creates a new CloudEvent.

**Common Event Methods:**
- `SetID(string)` / `ID() string`
- `SetType(string)` / `Type() string`
- `SetSource(string)` / `Source() string`
- `SetSubject(string)` / `Subject() string`
- `SetTime(time.Time)` / `Time() time.Time`
- `SetData(contentType string, data interface{}) error` / `Data() interface{}`
- `SetDataContentType(string)` / `DataContentType() string`

## Best Practices

### 1. Event Type Naming

Use reverse DNS notation for event types:

```go
// Good
event.SetType("com.example.users.created")
event.SetType("com.example.orders.shipped")

// Avoid
event.SetType("user_created")
event.SetType("shipped")
```

### 2. Source Identification

Use clear, consistent source identifiers:

```go
// Good
event.SetSource("users-service")
event.SetSource("example.com/payment-gateway")

// Avoid
event.SetSource("server1")
event.SetSource("app")
```

### 3. Error Handling

Always check result status:

```go
// Good
result := client.Send(event)
if result.IsUndelivered() {
    // Retry or log error
    log.Printf("Send failed: %s", result.Error())
    return
}

// Avoid
client.Send(event) // Ignoring result
```

### 4. Resource Cleanup

Stop receivers when shutting down:

```go
receiver.StartReceiver(handler, failureFunc)
defer receiver.StopReceiver() // Ensure cleanup

// Or with signal handling
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
<-sigChan
receiver.StopReceiver()
```

### 5. Structured Data

Use structured data formats:

```go
// Good - structured JSON
event.SetData("application/json", map[string]interface{}{
    "user_id": 123,
    "action":  "login",
})

// Avoid - unstructured strings
event.SetData("text/plain", "user 123 logged in")
```

### 6. Server Graceful Shutdown

Use appropriate timeout for graceful shutdown:

```go
// Give active connections time to complete
server.Stop(30 * time.Second)
```

## Event Patterns

### Fire-and-Forget

```go
result := client.Send(event)
// No response expected
```

### Request-Response

```go
response, result := client.Request(event)
// Synchronous response
```

### Publish-Subscribe

```go
// Publisher
client.Send(event)

// Multiple subscribers
receiver1.StartReceiver(handler1, failureFunc)
receiver2.StartReceiver(handler2, failureFunc)
```

## Error Handling

### Network Errors

```go
result := client.Send(event)
if result.IsUndelivered() {
    log.Printf("Network error: %s", result.Error())
    // Connection refused, timeout, etc.
}
```

### Application Errors

```go
result := client.Send(event)
if result.IsNACK() {
    statusCode, _ := result.GetHttpStatusCode()
    log.Printf("Application rejected: %d - %s", statusCode, result.Error())
    // 4xx or 5xx status codes
}
```

### Validation Errors

```go
event := cloudevents.NewEvent()
// Missing required fields
result := client.Send(event)
// Will fail validation
```

## Testing

### Unit Testing with Mock Events

```go
func TestEventHandler(t *testing.T) {
    event := cloudevents.NewEvent()
    event.SetType("com.example.test")
    event.SetSource("test")
    event.SetData("application/json", map[string]string{"key": "value"})
    
    // Test handler
    response, result := handler(event)
    
    if !result.IsACK() {
        t.Errorf("Expected ACK, got: %s", result.Error())
    }
}
```

### Integration Testing

```go
func TestClientServer(t *testing.T) {
    // Start server
    var server cloudevents.Server
    handler := func(event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
        return nil, cloudevents.NewHTTPResult(200, "OK")
    }
    server.Start(":8081", handler, func(err error) { t.Fatal(err) })
    defer server.Stop(5 * time.Second)
    
    // Test client
    client, _ := cloudevents.NewHttp("http://localhost:8081", nil, nil)
    event := cloudevents.NewEvent()
    event.SetType("test")
    event.SetSource("test")
    
    result := client.Send(event)
    if result.IsUndelivered() {
        t.Fatal(result.Error())
    }
}
```

## Performance Considerations

1. **Connection Pooling** - HTTP client reuses connections
2. **Async Processing** - Use StartReceiver for non-blocking event receipt
3. **Batch Events** - Group related events when possible
4. **Timeout Configuration** - Set appropriate timeouts for network operations
5. **Resource Cleanup** - Always stop receivers and servers

## Dependencies

- `github.com/cloudevents/sdk-go/v2` - Official CloudEvents SDK
- `github.com/common-library/go/http` - HTTP server utilities

## CloudEvents Specification

This package implements CloudEvents v1.0 specification:
- [CloudEvents Spec](https://github.com/cloudevents/spec/blob/v1.0/spec.md)
- [HTTP Protocol Binding](https://github.com/cloudevents/spec/blob/v1.0/http-protocol-binding.md)

## Further Reading

- [CloudEvents Primer](https://github.com/cloudevents/spec/blob/v1.0/primer.md)
- [CloudEvents SDK Go](https://github.com/cloudevents/sdk-go)
- [Event-Driven Architecture](https://martinfowler.com/articles/201701-event-driven.html)
