# SLog

Structured logging with asynchronous output and flexible configuration based on Go's log/slog.

## Overview

The slog package provides enhanced structured logging capabilities with features like asynchronous log processing, multiple output destinations, automatic file rotation, and caller information tracking. It wraps Go's standard log/slog package with additional convenience features.

## Features

- **Structured Logging** - Key-value pair logging in JSON format
- **Multiple Log Levels** - Trace, Debug, Info, Warn, Error, Fatal
- **Asynchronous Processing** - Queue-based non-blocking log writes
- **Flexible Output** - stdout, stderr, or file output
- **Daily Rotation** - Automatic log file rotation by date
- **Caller Tracking** - Optional file/line/function information
- **Thread-Safe** - Concurrent logging from multiple goroutines
- **Flush Control** - Ensure all logs are written before shutdown

## Installation

```bash
go get -u github.com/common-library/go/log/slog
```

## Quick Start

```go
import "github.com/common-library/go/log/slog"

var logger slog.Log

func main() {
    defer logger.Flush()
    
    logger.SetLevel(slog.LevelInfo)
    logger.SetOutputToFile("app", "log", true)
    
    logger.Info("Server started", "port", 8080)
    logger.Error("Connection failed", "error", err.Error())
}
```

## API Reference

### Log Levels

```go
const (
    LevelTrace = Level(-8)  // Most verbose
    LevelDebug = Level(-4)  // Debug information
    LevelInfo  = Level(0)   // General information
    LevelWarn  = Level(4)   // Warning messages
    LevelError = Level(8)   // Error messages
    LevelFatal = Level(12)  // Critical errors
)
```

### Log Type

```go
type Log struct {
    // Unexported fields
}
```

### Logging Methods

#### Trace

```go
func (l *Log) Trace(message string, arguments ...any)
```

Logs at trace level. Arguments are key-value pairs.

#### Debug

```go
func (l *Log) Debug(message string, arguments ...any)
```

Logs at debug level.

#### Info

```go
func (l *Log) Info(message string, arguments ...any)
```

Logs at info level.

#### Warn

```go
func (l *Log) Warn(message string, arguments ...any)
```

Logs at warning level.

#### Error

```go
func (l *Log) Error(message string, arguments ...any)
```

Logs at error level.

#### Fatal

```go
func (l *Log) Fatal(message string, arguments ...any)
```

Logs at fatal level (does NOT terminate application).

### Configuration Methods

#### SetLevel

```go
func (l *Log) SetLevel(level Level)
```

Sets minimum log level threshold.

#### SetOutputToStdout

```go
func (l *Log) SetOutputToStdout()
```

Redirects logs to standard output.

#### SetOutputToStderr

```go
func (l *Log) SetOutputToStderr()
```

Redirects logs to standard error.

#### SetOutputToFile

```go
func (l *Log) SetOutputToFile(fileName, fileExtensionName string, addDate bool)
```

Redirects logs to file with optional daily rotation.

#### SetWithCallerInfo

```go
func (l *Log) SetWithCallerInfo(withCallerInfo bool)
```

Enables/disables caller information in logs.

#### GetLevel

```go
func (l *Log) GetLevel() Level
```

Returns current log level.

#### Flush

```go
func (l *Log) Flush()
```

Blocks until all queued logs are written.

## Complete Examples

### Basic Logging

```go
package main

import (
    "github.com/common-library/go/log/slog"
)

var logger slog.Log

func main() {
    defer logger.Flush()
    
    logger.SetLevel(slog.LevelInfo)
    logger.SetOutputToStdout()
    
    logger.Info("Application started")
    logger.Info("Processing request", "method", "GET", "path", "/api/users")
    logger.Info("Application stopped")
}
```

### File Logging with Rotation

```go
package main

import (
    "github.com/common-library/go/log/slog"
)

var logger slog.Log

func main() {
    defer logger.Flush()
    
    // Creates files like: app_20231218.log, app_20231219.log, etc.
    logger.SetOutputToFile("app", "log", true)
    logger.SetLevel(slog.LevelDebug)
    
    logger.Debug("Debug information", "value", 42)
    logger.Info("Server started", "port", 8080)
    logger.Warn("High memory usage", "percent", 85)
}
```

### Error Handling

```go
package main

import (
    "errors"
    "github.com/common-library/go/log/slog"
)

var logger slog.Log

func processData() error {
    defer logger.Flush()
    
    logger.SetLevel(slog.LevelInfo)
    logger.SetOutputToFile("errors", "log", true)
    
    if err := validateInput(); err != nil {
        logger.Error("Validation failed",
            "error", err.Error(),
            "function", "validateInput",
        )
        return err
    }
    
    logger.Info("Processing completed successfully")
    return nil
}

func validateInput() error {
    return errors.New("invalid input format")
}

func main() {
    processData()
}
```

### With Caller Information

```go
package main

import (
    "github.com/common-library/go/log/slog"
)

var logger slog.Log

func processRequest(id int) {
    logger.Info("Processing started", "requestID", id)
    // Includes: {"CallerInfo":{"File":"main.go","Line":15,"Function":"main.processRequest"},...}
}

func main() {
    defer logger.Flush()
    
    logger.SetLevel(slog.LevelDebug)
    logger.SetOutputToStdout()
    logger.SetWithCallerInfo(true) // Enable caller info
    
    processRequest(123)
}
```

### Multiple Log Levels

```go
package main

import (
    "github.com/common-library/go/log/slog"
    "time"
)

var logger slog.Log

func main() {
    defer logger.Flush()
    
    logger.SetLevel(slog.LevelTrace)
    logger.SetOutputToFile("detailed", "log", true)
    
    logger.Trace("Entering function", "function", "main")
    logger.Debug("Configuration loaded", "config", "app.yaml")
    logger.Info("Server starting")
    logger.Warn("Cache nearly full", "usage", "90%")
    logger.Error("Failed to connect", "service", "database")
    logger.Fatal("Critical system error", "component", "core")
}
```

### Dynamic Level Change

```go
package main

import (
    "os"
    "github.com/common-library/go/log/slog"
)

var logger slog.Log

func main() {
    defer logger.Flush()
    
    // Set level based on environment
    if os.Getenv("DEBUG") == "true" {
        logger.SetLevel(slog.LevelDebug)
        logger.SetWithCallerInfo(true)
    } else {
        logger.SetLevel(slog.LevelInfo)
        logger.SetWithCallerInfo(false)
    }
    
    logger.SetOutputToStdout()
    
    logger.Debug("Debug mode enabled") // Only if DEBUG=true
    logger.Info("Application started")
}
```

### HTTP Server Logging

```go
package main

import (
    "net/http"
    "time"
    "github.com/common-library/go/log/slog"
)

var logger slog.Log

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        logger.Info("Request started",
            "method", r.Method,
            "path", r.URL.Path,
            "remote", r.RemoteAddr,
        )
        
        next.ServeHTTP(w, r)
        
        logger.Info("Request completed",
            "method", r.Method,
            "path", r.URL.Path,
            "duration", time.Since(start).Milliseconds(),
        )
    })
}

func main() {
    defer logger.Flush()
    
    logger.SetLevel(slog.LevelInfo)
    logger.SetOutputToFile("access", "log", true)
    
    http.Handle("/", loggingMiddleware(http.HandlerFunc(handler)))
    http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello"))
}
```

### Structured Application Logging

```go
package main

import (
    "github.com/common-library/go/log/slog"
)

var (
    appLogger   slog.Log
    auditLogger slog.Log
)

func main() {
    defer appLogger.Flush()
    defer auditLogger.Flush()
    
    // Application logs
    appLogger.SetLevel(slog.LevelInfo)
    appLogger.SetOutputToFile("app", "log", true)
    
    // Audit logs (separate file)
    auditLogger.SetLevel(slog.LevelInfo)
    auditLogger.SetOutputToFile("audit", "log", true)
    
    appLogger.Info("Application started")
    
    // Log user action to audit log
    auditLogger.Info("User action",
        "userID", 123,
        "action", "login",
        "ip", "192.168.1.1",
    )
    
    appLogger.Info("Processing request")
    
    auditLogger.Info("Data access",
        "userID", 123,
        "resource", "/api/users/456",
        "operation", "READ",
    )
}
```

## Best Practices

### 1. Always Flush Before Exit

```go
// Good: Ensures all logs are written
func main() {
    defer logger.Flush()
    // ... application code ...
}

// Avoid: May lose recent logs
func main() {
    // ... application code ...
    // No flush
}
```

### 2. Use Appropriate Log Levels

```go
// Good: Match severity to level
logger.Trace("Variable value", "x", x)     // Detailed debugging
logger.Debug("Function called", "name", f) // Development info
logger.Info("Server started", "port", 80)  // Normal operation
logger.Warn("Retry attempt", "count", 3)   // Potential issue
logger.Error("Failed to save", "err", err) // Operation failed
logger.Fatal("Config missing")             // Critical error

// Avoid: Wrong levels
logger.Error("Server started") // Not an error
logger.Info("Critical failure") // Should be Fatal
```

### 3. Use Structured Key-Value Pairs

```go
// Good: Structured data
logger.Info("User logged in",
    "userID", 123,
    "username", "alice",
    "ip", remoteAddr,
)

// Avoid: Unstructured messages
logger.Info(fmt.Sprintf("User %d (%s) logged in from %s", 123, "alice", remoteAddr))
```

### 4. Enable Caller Info Selectively

```go
// Good: Only in development
if os.Getenv("ENV") == "development" {
    logger.SetWithCallerInfo(true)
}

// Avoid: Always enabled in production (performance impact)
logger.SetWithCallerInfo(true)
```

### 5. Use Daily Rotation for Production

```go
// Good: Automatic rotation
logger.SetOutputToFile("app", "log", true)
// Creates: app_20231218.log, app_20231219.log, etc.

// Avoid: Single large file
logger.SetOutputToFile("app", "log", false)
// Creates: app.log (grows indefinitely)
```

## Performance Tips

1. **Asynchronous by Design** - Logs are queued and written asynchronously
2. **Flush Strategically** - Only flush when necessary (shutdown, critical points)
3. **Disable Caller Info in Production** - Adds overhead for getting stack information
4. **Use Appropriate Levels** - Set level to Info or higher in production
5. **Structured Logging** - More efficient than string formatting

## Testing

```go
func TestLogging(t *testing.T) {
    var logger slog.Log
    
    // Use file output for test verification
    logger.SetOutputToFile("test", "log", false)
    logger.SetLevel(slog.LevelDebug)
    
    logger.Info("Test message", "key", "value")
    logger.Flush()
    
    // Read and verify log file
    data, err := os.ReadFile("test.log")
    if err != nil {
        t.Fatal(err)
    }
    
    if !strings.Contains(string(data), "Test message") {
        t.Error("Log message not found")
    }
    
    // Cleanup
    os.Remove("test.log")
}
```

## Dependencies

- `log/slog` - Go standard library structured logging
- `github.com/common-library/go/collection` - Queue implementation
- `github.com/common-library/go/lock` - Mutex implementation
- `github.com/common-library/go/utility` - Caller info utilities

## Further Reading

- [Go slog package](https://pkg.go.dev/log/slog)
- [Structured Logging in Go](https://go.dev/blog/slog)
