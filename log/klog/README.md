# KLog

Kubernetes-style logging wrapper with optional caller information tracking.

## Overview

The klog package provides wrapper functions for k8s.io/klog/v2, the logging library used throughout Kubernetes. It adds optional caller information tracking to help identify log sources in large applications.

## Features

- **Kubernetes Standard** - Uses k8s.io/klog/v2 underneath
- **Multiple Formats** - Print, Printf, Println, and Structured logging
- **Log Levels** - Info, Error, and Fatal levels
- **Caller Tracking** - Optional file/line/function information
- **Structured Logging** - InfoS and ErrorS for key-value pairs
- **Thread-Safe** - Concurrent logging from multiple goroutines
- **Fatal Handling** - Fatal logs terminate application with os.Exit(255)

## Installation

```bash
go get -u github.com/common-library/go/log/klog
go get -u k8s.io/klog/v2
```

## Quick Start

```go
import "github.com/common-library/go/log/klog"

func main() {
    defer klog.Flush()
    
    klog.Info("Application started")
    klog.InfoS("Request processed", "method", "GET", "path", "/api/users")
    klog.Infof("Server listening on port %d", 8080)
}
```

## API Reference

### Info Logging

#### Info

```go
func Info(arguments ...any)
```

Logs informational messages (like fmt.Print).

#### InfoS

```go
func InfoS(message string, keysAndValues ...any)
```

Logs structured messages with key-value pairs.

#### Infof

```go
func Infof(format string, arguments ...any)
```

Logs formatted messages (like fmt.Printf).

#### Infoln

```go
func Infoln(arguments ...any)
```

Logs messages with newline (like fmt.Println).

### Error Logging

#### Error

```go
func Error(arguments ...any)
```

Logs error messages.

#### ErrorS

```go
func ErrorS(err error, message string, keysAndValues ...any)
```

Logs structured error messages with explicit error parameter.

#### Errorf

```go
func Errorf(format string, arguments ...any)
```

Logs formatted error messages.

#### Errorln

```go
func Errorln(arguments ...any)
```

Logs error messages with newline.

### Fatal Logging

#### Fatal

```go
func Fatal(arguments ...any)
```

Logs message and calls os.Exit(255).

#### Fatalf

```go
func Fatalf(format string, arguments ...any)
```

Logs formatted message and calls os.Exit(255).

#### Fatalln

```go
func Fatalln(arguments ...any)
```

Logs message with newline and calls os.Exit(255).

### Utility Functions

#### Flush

```go
func Flush()
```

Flushes all pending log I/O.

#### SetWithCallerInfo

```go
func SetWithCallerInfo(with bool)
```

Enables/disables caller information in logs.

## Complete Examples

### Basic Logging

```go
package main

import "github.com/common-library/go/log/klog"

func main() {
    defer klog.Flush()
    
    klog.Info("Application starting")
    klog.Info("Configuration loaded")
    klog.Info("Server ready")
}
```

### Structured Logging

```go
package main

import "github.com/common-library/go/log/klog"

func main() {
    defer klog.Flush()
    
    klog.InfoS("Server started",
        "port", 8080,
        "environment", "production",
        "version", "1.0.0",
    )
    
    klog.InfoS("Request processed",
        "method", "GET",
        "path", "/api/users",
        "duration", 45,
        "status", 200,
    )
}
```

### Formatted Logging

```go
package main

import "github.com/common-library/go/log/klog"

func main() {
    defer klog.Flush()
    
    port := 8080
    klog.Infof("Server listening on port %d", port)
    
    requestCount := 1000
    avgDuration := 45.5
    klog.Infof("Processed %d requests, avg duration: %.2fms", requestCount, avgDuration)
}
```

### Error Logging

```go
package main

import (
    "errors"
    "github.com/common-library/go/log/klog"
)

func main() {
    defer klog.Flush()
    
    err := connectDatabase()
    if err != nil {
        klog.Error("Database connection failed:", err)
        klog.ErrorS(err, "Failed to connect",
            "host", "localhost",
            "port", 5432,
            "database", "myapp",
        )
    }
}

func connectDatabase() error {
    return errors.New("connection timeout")
}
```

### With Caller Information

```go
package main

import "github.com/common-library/go/log/klog"

func processRequest(id int) {
    klog.Info("Processing request", id)
    // Output: [callerInfo:{File:"main.go" Line:8 Function:"main.processRequest"}] Processing request 123
}

func main() {
    defer klog.Flush()
    
    klog.SetWithCallerInfo(true) // Enable caller tracking
    
    processRequest(123)
}
```

### HTTP Server Logging

```go
package main

import (
    "net/http"
    "time"
    "github.com/common-library/go/log/klog"
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        klog.InfoS("Request started",
            "method", r.Method,
            "path", r.URL.Path,
            "remote", r.RemoteAddr,
        )
        
        next.ServeHTTP(w, r)
        
        klog.InfoS("Request completed",
            "method", r.Method,
            "path", r.URL.Path,
            "duration", time.Since(start).Milliseconds(),
        )
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello"))
}

func main() {
    defer klog.Flush()
    
    klog.Info("Starting HTTP server on :8080")
    
    http.Handle("/", loggingMiddleware(http.HandlerFunc(handler)))
    
    if err := http.ListenAndServe(":8080", nil); err != nil {
        klog.Fatal("Server failed:", err)
    }
}
```

### Kubernetes Controller Logging

```go
package main

import (
    "time"
    "github.com/common-library/go/log/klog"
)

type Controller struct {
    name string
}

func (c *Controller) Run() {
    klog.InfoS("Controller starting", "controller", c.name)
    
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            c.reconcile()
        }
    }
}

func (c *Controller) reconcile() {
    klog.InfoS("Reconciliation started", "controller", c.name)
    
    // Reconciliation logic...
    
    klog.InfoS("Reconciliation completed",
        "controller", c.name,
        "resourcesProcessed", 10,
    )
}

func main() {
    defer klog.Flush()
    
    controller := &Controller{name: "pod-controller"}
    controller.Run()
}
```

### Fatal Error Handling

```go
package main

import (
    "os"
    "github.com/common-library/go/log/klog"
)

func loadConfig() error {
    configFile := "config.yaml"
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        return err
    }
    return nil
}

func main() {
    defer klog.Flush()
    
    if err := loadConfig(); err != nil {
        klog.Fatalf("Configuration file required: %v", err)
        // Application terminates here with os.Exit(255)
    }
    
    klog.Info("Application started successfully")
}
```

### Environment-Based Configuration

```go
package main

import (
    "os"
    "github.com/common-library/go/log/klog"
)

func main() {
    defer klog.Flush()
    
    // Enable caller info in development
    if os.Getenv("ENV") == "development" {
        klog.SetWithCallerInfo(true)
        klog.Info("Development mode: caller info enabled")
    } else {
        klog.SetWithCallerInfo(false)
        klog.Info("Production mode: caller info disabled")
    }
    
    klog.InfoS("Application started",
        "environment", os.Getenv("ENV"),
        "version", "1.0.0",
    )
}
```

### Multi-Level Logging

```go
package main

import (
    "errors"
    "github.com/common-library/go/log/klog"
)

func main() {
    defer klog.Flush()
    
    // Informational
    klog.Info("Application initializing")
    klog.InfoS("Configuration loaded", "source", "config.yaml")
    
    // Warnings (use Error level, klog doesn't have Warn)
    klog.Error("Cache miss, using default value")
    
    // Errors
    err := processData()
    if err != nil {
        klog.Error("Processing failed:", err)
        klog.ErrorS(err, "Data processing error",
            "component", "processor",
            "operation", "transform",
        )
    }
    
    // Fatal (terminates application)
    if criticalError() {
        klog.Fatal("Critical system error, cannot continue")
    }
}

func processData() error {
    return errors.New("invalid data format")
}

func criticalError() bool {
    return false
}
```

## Best Practices

### 1. Always Flush Before Exit

```go
// Good: Ensures logs are written
func main() {
    defer klog.Flush()
    // ... code ...
}

// Avoid: May lose logs
func main() {
    // ... code ...
    // No flush
}
```

### 2. Use Structured Logging

```go
// Good: Structured data
klog.InfoS("Request processed",
    "method", "GET",
    "path", "/api/users",
    "duration", 45,
)

// Avoid: Concatenated strings
klog.Info("Request processed: GET /api/users duration=45ms")
```

### 3. Use ErrorS for Errors

```go
// Good: Explicit error parameter
klog.ErrorS(err, "Database query failed",
    "query", sql,
    "table", "users",
)

// Acceptable but less structured
klog.Error("Database query failed:", err)
```

### 4. Enable Caller Info Selectively

```go
// Good: Only in development
if os.Getenv("DEBUG") == "true" {
    klog.SetWithCallerInfo(true)
}

// Avoid: Always enabled (performance impact)
klog.SetWithCallerInfo(true)
```

### 5. Use Fatal Appropriately

```go
// Good: Truly unrecoverable errors
if configFile == nil {
    klog.Fatal("Configuration file required")
}

// Avoid: Recoverable errors
if err := connect(); err != nil {
    klog.Fatal(err) // Should use Error instead
}
```

## Comparison with slog

| Feature | klog | slog |
|---------|------|------|
| Backend | k8s.io/klog/v2 | log/slog |
| Async | No | Yes |
| Output Dest | Flags-based | Configurable |
| Rotation | External | Built-in (daily) |
| Fatal | Exits app | No exit |
| Use Case | Kubernetes apps | General apps |

## Performance Tips

1. **Minimize Caller Info** - Only enable in development
2. **Use Structured Logging** - InfoS/ErrorS more efficient than formatted strings
3. **Flush Strategically** - Only at shutdown or critical points
4. **Avoid Fatal in Libraries** - Use Error instead

## Testing

```go
func TestLogging(t *testing.T) {
    // Redirect klog output for testing
    flag.Set("logtostderr", "false")
    flag.Set("log_file", "test.log")
    
    klog.Info("Test message")
    klog.Flush()
    
    // Verify log file
    data, err := os.ReadFile("test.log")
    if err != nil {
        t.Fatal(err)
    }
    
    if !strings.Contains(string(data), "Test message") {
        t.Error("Log not found")
    }
    
    // Cleanup
    os.Remove("test.log")
}
```

## Dependencies

- `k8s.io/klog/v2` - Kubernetes logging library
- `github.com/common-library/go/utility` - Caller info utilities

## Further Reading

- [klog Documentation](https://github.com/kubernetes/klog)
- [Kubernetes Logging Conventions](https://kubernetes.io/docs/concepts/cluster-administration/logging/)
