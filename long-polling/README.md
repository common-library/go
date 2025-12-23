# Long Polling

HTTP long polling server and client implementation for real-time event delivery.

## Overview

The long-polling package provides a production-ready long polling solution built on top of golongpoll. It enables real-time communication between servers and clients using HTTP without requiring WebSocket connections, making it ideal for environments where WebSockets are blocked or unavailable.

## Features

- **Long Polling Server** - Event-driven server with configurable timeouts
- **Event Categories** - Organize events into logical categories
- **File Persistence** - Optional file-based event storage for durability
- **Custom Middleware** - Authentication and validation hooks
- **Client Helpers** - Simple subscription and publishing functions
- **Graceful Shutdown** - Proper server cleanup with timeout control

## Installation

```bash
go get -u github.com/common-library/go/long-polling
```

## Quick Start

### Server

```go
server := &long_polling.Server{}
err := server.Start(long_polling.ServerInfo{
    Address: ":8080",
    TimeoutSeconds: 120,
    SubscriptionURI: "/events",
    PublishURI: "/publish",
}, long_polling.FilePersistorInfo{Use: false}, nil)
```

### Client

```go
// Subscribe
response, err := long_polling.Subscription(
    "http://localhost:8080/events",
    nil,
    long_polling.SubscriptionRequest{
        Category: "notifications",
        TimeoutSeconds: 60,
    },
    "", "", nil,
)

// Publish
_, err = long_polling.Publish(
    "http://localhost:8080/publish",
    10 * time.Second,
    nil,
    long_polling.PublishRequest{
        Category: "notifications",
        Data: `{"message": "Hello"}`,
    },
    "", "", nil,
)
```

## API Reference

### Server Types

#### ServerInfo

```go
type ServerInfo struct {
    Address        string
    TimeoutSeconds int
    SubscriptionURI                string
    HandlerToRunBeforeSubscription func(w http.ResponseWriter, r *http.Request) bool
    PublishURI                string
    HandlerToRunBeforePublish func(w http.ResponseWriter, r *http.Request) bool
}
```

Server configuration parameters.

#### FilePersistorInfo

```go
type FilePersistorInfo struct {
    Use                     bool
    FileName                string
    WriteBufferSize         int
    WriteFlushPeriodSeconds int
}
```

File persistence configuration for event durability.

### Server Methods

#### Start

```go
func (s *Server) Start(serverInfo ServerInfo, filePersistorInfo FilePersistorInfo, 
    listenAndServeFailureFunc func(err error)) error
```

Starts the long polling server.

#### Stop

```go
func (s *Server) Stop(shutdownTimeout time.Duration) error
```

Gracefully shuts down the server.

### Client Types

#### SubscriptionRequest

```go
type SubscriptionRequest struct {
    Category       string
    TimeoutSeconds int
    SinceTime      int64
    LastID         string
}
```

Subscription parameters.

#### SubscriptionResponse

```go
type SubscriptionResponse struct {
    Header     http.Header
    StatusCode int
    Events     []Event
}
```

Subscription response with events.

#### PublishRequest

```go
type PublishRequest struct {
    Category string
    Data     string
}
```

Event publishing parameters.

### Client Functions

#### Subscription

```go
func Subscription(url string, header map[string][]string, request SubscriptionRequest,
    username, password string, transport *http.Transport) (SubscriptionResponse, error)
```

Subscribes to server events.

#### Publish

```go
func Publish(url string, timeout time.Duration, header map[string][]string,
    publishRequest PublishRequest, username, password string, 
    transport *http.Transport) (http.Response, error)
```

Publishes an event to the server.

## Complete Examples

### Basic Server

```go
package main

import (
    "log"
    "time"
    "github.com/common-library/go/long-polling"
)

func main() {
    server := &long_polling.Server{}
    
    err := server.Start(long_polling.ServerInfo{
        Address:         ":8080",
        TimeoutSeconds:  120,
        SubscriptionURI: "/events",
        PublishURI:      "/publish",
    }, long_polling.FilePersistorInfo{
        Use: false,
    }, func(err error) {
        log.Fatalf("Server error: %v", err)
    })
    
    if err != nil {
        log.Fatalf("Failed to start: %v", err)
    }
    
    // Wait for shutdown signal
    // ...
    
    server.Stop(10 * time.Second)
}
```

### Server with Authentication

```go
package main

import (
    "log"
    "net/http"
    "github.com/common-library/go/long-polling"
)

func authenticateSubscription(w http.ResponseWriter, r *http.Request) bool {
    token := r.Header.Get("Authorization")
    if token != "Bearer secret-token" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return false
    }
    return true
}

func authenticatePublish(w http.ResponseWriter, r *http.Request) bool {
    token := r.Header.Get("Authorization")
    if token != "Bearer admin-token" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return false
    }
    return true
}

func main() {
    server := &long_polling.Server{}
    
    err := server.Start(long_polling.ServerInfo{
        Address:                        ":8080",
        TimeoutSeconds:                 120,
        SubscriptionURI:                "/events",
        HandlerToRunBeforeSubscription: authenticateSubscription,
        PublishURI:                     "/publish",
        HandlerToRunBeforePublish:      authenticatePublish,
    }, long_polling.FilePersistorInfo{
        Use: false,
    }, nil)
    
    if err != nil {
        log.Fatal(err)
    }
}
```

### Server with File Persistence

```go
package main

import (
    "log"
    "github.com/common-library/go/long-polling"
)

func main() {
    server := &long_polling.Server{}
    
    err := server.Start(long_polling.ServerInfo{
        Address:         ":8080",
        TimeoutSeconds:  120,
        SubscriptionURI: "/events",
        PublishURI:      "/publish",
    }, long_polling.FilePersistorInfo{
        Use:                     true,
        FileName:                "/var/lib/longpoll/events.db",
        WriteBufferSize:         1000,
        WriteFlushPeriodSeconds: 5,
    }, nil)
    
    if err != nil {
        log.Fatal(err)
    }
}
```

### Basic Client Subscription

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/long-polling"
)

func main() {
    for {
        response, err := long_polling.Subscription(
            "http://localhost:8080/events",
            nil,
            long_polling.SubscriptionRequest{
                Category:       "notifications",
                TimeoutSeconds: 60,
            },
            "", "", nil,
        )
        
        if err != nil {
            log.Printf("Subscription error: %v", err)
            continue
        }
        
        if response.StatusCode == 200 {
            for _, event := range response.Events {
                fmt.Printf("Event: %s - %s\n", event.Category, event.Data)
            }
        }
    }
}
```

### Client with Authentication

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/common-library/go/long-polling"
)

func main() {
    headers := map[string][]string{
        "Authorization": {"Bearer secret-token"},
    }
    
    for {
        response, err := long_polling.Subscription(
            "http://localhost:8080/events",
            headers,
            long_polling.SubscriptionRequest{
                Category:       "notifications",
                TimeoutSeconds: 60,
            },
            "", "", nil,
        )
        
        if err != nil {
            log.Printf("Error: %v", err)
            continue
        }
        
        if response.StatusCode == http.StatusUnauthorized {
            log.Fatal("Authentication failed")
        }
        
        for _, event := range response.Events {
            fmt.Printf("Event: %s\n", event.Data)
        }
    }
}
```

### Client Publishing Events

```go
package main

import (
    "log"
    "time"
    "github.com/common-library/go/long-polling"
)

func main() {
    response, err := long_polling.Publish(
        "http://localhost:8080/publish",
        10 * time.Second,
        nil,
        long_polling.PublishRequest{
            Category: "notifications",
            Data:     `{"type": "alert", "message": "System maintenance"}`,
        },
        "", "", nil,
    )
    
    if err != nil {
        log.Fatalf("Publish error: %v", err)
    }
    
    if response.StatusCode == 200 {
        log.Println("Event published successfully")
    }
}
```

### Client with Incremental Updates

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/long-polling"
)

func main() {
    var lastEventID string
    
    for {
        request := long_polling.SubscriptionRequest{
            Category:       "updates",
            TimeoutSeconds: 60,
        }
        
        if lastEventID != "" {
            request.LastID = lastEventID
        }
        
        response, err := long_polling.Subscription(
            "http://localhost:8080/events",
            nil,
            request,
            "", "", nil,
        )
        
        if err != nil {
            log.Printf("Error: %v", err)
            continue
        }
        
        for _, event := range response.Events {
            fmt.Printf("Update: %s\n", event.Data)
            lastEventID = event.ID
        }
    }
}
```

### Multi-Category Subscription

```go
package main

import (
    "fmt"
    "log"
    "sync"
    "github.com/common-library/go/long-polling"
)

func subscribe(category string, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for {
        response, err := long_polling.Subscription(
            "http://localhost:8080/events",
            nil,
            long_polling.SubscriptionRequest{
                Category:       category,
                TimeoutSeconds: 60,
            },
            "", "", nil,
        )
        
        if err != nil {
            log.Printf("[%s] Error: %v", category, err)
            continue
        }
        
        for _, event := range response.Events {
            fmt.Printf("[%s] Event: %s\n", category, event.Data)
        }
    }
}

func main() {
    var wg sync.WaitGroup
    
    categories := []string{"notifications", "alerts", "updates"}
    
    for _, category := range categories {
        wg.Add(1)
        go subscribe(category, &wg)
    }
    
    wg.Wait()
}
```

### Complete Chat Application

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "time"
    "github.com/common-library/go/long-polling"
)

func receiveMessages() {
    for {
        response, err := long_polling.Subscription(
            "http://localhost:8080/events",
            nil,
            long_polling.SubscriptionRequest{
                Category:       "chat",
                TimeoutSeconds: 60,
            },
            "", "", nil,
        )
        
        if err != nil {
            continue
        }
        
        for _, event := range response.Events {
            fmt.Printf("Message: %s\n", event.Data)
        }
    }
}

func sendMessages() {
    scanner := bufio.NewScanner(os.Stdin)
    
    for {
        fmt.Print("Enter message: ")
        if !scanner.Scan() {
            break
        }
        
        message := scanner.Text()
        
        _, err := long_polling.Publish(
            "http://localhost:8080/publish",
            10 * time.Second,
            nil,
            long_polling.PublishRequest{
                Category: "chat",
                Data:     message,
            },
            "", "", nil,
        )
        
        if err != nil {
            log.Printf("Failed to send: %v", err)
        }
    }
}

func main() {
    go receiveMessages()
    sendMessages()
}
```

## Best Practices

### 1. Set Appropriate Timeouts

```go
// Good: Reasonable timeout for user interactions
ServerInfo{
    TimeoutSeconds: 60,  // 1 minute
}

// Avoid: Too short (excessive requests)
ServerInfo{
    TimeoutSeconds: 5,   // Too short
}

// Avoid: Too long (resource waste)
ServerInfo{
    TimeoutSeconds: 600, // 10 minutes
}
```

### 2. Handle Subscription Loops

```go
// Good: Continue on errors
for {
    response, err := long_polling.Subscription(...)
    if err != nil {
        time.Sleep(1 * time.Second)
        continue
    }
    processEvents(response.Events)
}

// Avoid: Exit on first error
response, err := long_polling.Subscription(...)
if err != nil {
    return err // Don't give up
}
```

### 3. Use Categories Wisely

```go
// Good: Specific categories
PublishRequest{
    Category: "user.1234.notifications",
    Data: "...",
}

// Avoid: Generic categories
PublishRequest{
    Category: "events", // Too broad
    Data: "...",
}
```

### 4. Implement Graceful Shutdown

```go
// Good: Graceful shutdown
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
<-sigChan

server.Stop(10 * time.Second)

// Avoid: Abrupt termination
os.Exit(0) // No cleanup
```

### 5. Enable Persistence for Critical Events

```go
// Good: Persist important events
FilePersistorInfo{
    Use: true,
    FileName: "/var/lib/app/events.db",
    WriteBufferSize: 1000,
    WriteFlushPeriodSeconds: 5,
}

// Consider: No persistence for transient events
FilePersistorInfo{
    Use: false, // OK for temporary notifications
}
```

## Performance Tips

1. **Connection Pooling** - Use custom transport with connection pooling for multiple requests
2. **Timeout Configuration** - Balance between latency and resource usage
3. **File Persistence** - Use reasonable buffer size and flush period
4. **Event Categories** - Use specific categories to reduce client filtering
5. **Graceful Shutdown** - Always call Stop() to flush pending events

## Testing

```go
func TestLongPolling(t *testing.T) {
    server := &long_polling.Server{}
    
    go server.Start(long_polling.ServerInfo{
        Address:         ":8081",
        TimeoutSeconds:  10,
        SubscriptionURI: "/events",
        PublishURI:      "/publish",
    }, long_polling.FilePersistorInfo{Use: false}, nil)
    
    time.Sleep(100 * time.Millisecond)
    defer server.Stop(1 * time.Second)
    
    // Publish event
    _, err := long_polling.Publish(
        "http://localhost:8081/publish",
        5 * time.Second,
        nil,
        long_polling.PublishRequest{
            Category: "test",
            Data:     "hello",
        },
        "", "", nil,
    )
    
    if err != nil {
        t.Fatalf("Publish failed: %v", err)
    }
    
    // Subscribe
    response, err := long_polling.Subscription(
        "http://localhost:8081/events",
        nil,
        long_polling.SubscriptionRequest{
            Category:       "test",
            TimeoutSeconds: 5,
        },
        "", "", nil,
    )
    
    if err != nil {
        t.Fatalf("Subscription failed: %v", err)
    }
    
    if len(response.Events) == 0 {
        t.Error("No events received")
    }
}
```

## Dependencies

- `github.com/jcuga/golongpoll` - Core long polling implementation
- `github.com/gorilla/mux` - HTTP router
- `github.com/common-library/go/http` - HTTP utilities
- `github.com/common-library/go/json` - JSON utilities
- `github.com/google/go-querystring` - Query string encoding

## Further Reading

- [HTTP Long Polling](https://javascript.info/long-polling)
- [golongpoll Documentation](https://github.com/jcuga/golongpoll)
- [Real-time Web Technologies](https://www.ably.io/blog/websockets-vs-long-polling)
