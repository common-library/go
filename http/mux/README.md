# Gorilla Mux HTTP Server Package

Gorilla Mux router wrapper package for Go.

## Features

- **Based on Gorilla Mux**: Powerful URL router and dispatcher
- **Easy routing**: Path variables and method restrictions support
- **Path Prefix routing**: Path prefix-based routing
- **Middleware**: Global middleware support
- **Graceful Shutdown**: Timeout-based safe shutdown
- **Thread-safe**: Concurrency safety guaranteed

## Installation

```bash
go get github.com/common-library/go/http/mux
go get github.com/gorilla/mux
```

## When to Choose Gorilla Mux

- ✅ Complex URL pattern matching needed
- ✅ Subrouters and path prefixes needed
- ✅ Regular expression-based routing needed
- ✅ URL builder functionality needed
- ✅ Standard net/http compatibility important

## Usage Examples

### Basic Server

```go
package main

import (
    "net/http"
    "github.com/common-library/go/http/mux"
)

func main() {
    server := &mux.Server{}
    
    // Simple handler
    server.RegisterHandlerFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    // JSON response
    server.RegisterHandlerFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"message":"Hello from Gorilla Mux"}`))
    })
    
    // Start server
    err := server.Start(":8080", func(err error) {
        log.Fatalf("Server error: %v", err)
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Wait for shutdown
    select {}
}
```

### Route Registration Methods

```go
server := &mux.Server{}

// Handler registration
handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("handler"))
})
server.RegisterHandler("/api", handler)

// HandlerFunc registration
server.RegisterHandlerFunc("/users", getUsersHandler)

// HTTP method restriction
server.RegisterHandlerFunc("/users", createUserHandler, http.MethodPost)
server.RegisterHandlerFunc("/users/{id}", updateUserHandler, http.MethodPut)
server.RegisterHandlerFunc("/users/{id}", deleteUserHandler, http.MethodDelete)

// Allow multiple methods
server.RegisterHandlerFunc("/data", dataHandler, http.MethodGet, http.MethodPost)
```

### Path Variables Usage

```go
import "github.com/gorilla/mux"

server := &mux.Server{}

// Path variables
server.RegisterHandlerFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    fmt.Fprintf(w, "User ID: %s", id)
})

// Multiple path variables
server.RegisterHandlerFunc("/posts/{category}/{id}", func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    category := vars["category"]
    id := vars["id"]
    fmt.Fprintf(w, "Category: %s, ID: %s", category, id)
})

// Regular expression constraints
router := server.GetRouter()
router.HandleFunc("/articles/{id:[0-9]+}", articleHandler)
```

### Path Prefix Routing

```go
server := &mux.Server{}

// Serve static files
fileServer := http.FileServer(http.Dir("./static"))
server.RegisterPathPrefixHandler("/static/", http.StripPrefix("/static/", fileServer))

// API version routing
server.RegisterPathPrefixHandlerFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("API v1"))
}, http.MethodGet)

server.RegisterPathPrefixHandlerFunc("/api/v2/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("API v2"))
}, http.MethodGet)

// Admin pages
server.RegisterPathPrefixHandler("/admin/", adminHandler, http.MethodGet, http.MethodPost)
```

### Middleware Usage

```go
server := &mux.Server{}

// Logging middleware
loggingMiddleware := func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
        next.ServeHTTP(w, r)
    })
}

// Authentication middleware
authMiddleware := func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Register global middleware
server.Use(loggingMiddleware)
server.Use(authMiddleware)

// Register routes
server.RegisterHandlerFunc("/protected", protectedHandler)

server.Start(":8080", nil)
```

### Advanced Router Settings

```go
server := &mux.Server{}

// Direct router access
router := server.GetRouter()

// Option settings
router.StrictSlash(true)  // Treat /path and /path/ identically
router.SkipClean(false)   // Enable URL cleanup

// Create subrouter
api := router.PathPrefix("/api").Subrouter()
api.HandleFunc("/users", usersHandler).Methods("GET")
api.HandleFunc("/users", createUserHandler).Methods("POST")

// Host-based routing
router.Host("api.example.com").Handler(apiHandler)
router.Host("www.example.com").Handler(webHandler)

// Scheme restriction
router.Schemes("https").Handler(secureHandler)

// Query matching
router.Queries("key", "value").HandlerFunc(queryHandler)
```

### Graceful Shutdown

```go
import (
    "os"
    "os/signal"
    "syscall"
    "time"
)

server := &mux.Server{}
server.RegisterHandlerFunc("/", homeHandler)

// Start server
server.Start(":8080", func(err error) {
    log.Printf("Server error: %v", err)
})

// Wait for signal
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// Graceful shutdown (10 second timeout)
log.Println("Shutting down server...")
if err := server.Stop(10 * time.Second); err != nil {
    log.Fatal(err)
}
log.Println("Server stopped")
```

### Server Status Check

```go
server := &mux.Server{}

// Before starting server
if !server.IsRunning() {
    server.Start(":8080", nil)
}

// Prevent duplicate start
if server.IsRunning() {
    log.Println("Server already running")
    return
}

// Health check endpoint
server.RegisterHandlerFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    status := "down"
    if server.IsRunning() {
        status = "up"
    }
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status":"%s"}`, status)
})
```

### URL Builder

```go
router := server.GetRouter()

// Named routes
router.HandleFunc("/articles/{category}/{id:[0-9]+}", articleHandler).
    Name("article")

// Generate URL
url, err := router.Get("article").URL("category", "tech", "id", "42")
// Result: /articles/tech/42

// Current route information
router.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
    route := mux.CurrentRoute(r)
    pathTemplate, _ := route.GetPathTemplate()
    fmt.Fprintf(w, "Route: %s", pathTemplate)
})
```

## API Documentation

### Server

#### `RegisterHandler(path string, handler http.Handler, methods ...string)`
Registers an HTTP handler.

- **path**: URL path pattern (supports path variables: `/users/{id}`)
- **handler**: HTTP Handler
- **methods**: Optional HTTP method restrictions

#### `RegisterHandlerFunc(path string, handlerFunc http.HandlerFunc, methods ...string)`
Registers an HTTP handler function.

- **path**: URL path pattern
- **handlerFunc**: HTTP HandlerFunc
- **methods**: Optional HTTP method restrictions

#### `RegisterPathPrefixHandler(prefix string, handler http.Handler, methods ...string)`
Registers a handler for a path prefix.

- **prefix**: Path prefix (e.g., `/api/`, `/static/`)
- **handler**: HTTP Handler
- **methods**: Optional HTTP method restrictions

#### `RegisterPathPrefixHandlerFunc(prefix string, handlerFunc http.HandlerFunc, methods ...string)`
Registers a handler function for a path prefix.

- **prefix**: Path prefix
- **handlerFunc**: HTTP HandlerFunc
- **methods**: Optional HTTP method restrictions

#### `Use(middleware ...mux.MiddlewareFunc)`
Registers global middleware.

- **middleware**: Middleware functions

#### `Start(address string, listenAndServeFailureFunc func(err error)) error`
Starts the server.

- **address**: Server address (e.g., ":8080")
- **listenAndServeFailureFunc**: Error callback (excludes http.ErrServerClosed)
- **return**: Error if already started, otherwise nil

#### `Stop(shutdownTimeout time.Duration) error`
Safely shuts down the server.

- **shutdownTimeout**: Shutdown wait time

#### `IsRunning() bool`
Returns the server running status.

#### `GetRouter() *mux.Router`
Provides direct access to the Gorilla Mux router.

## Framework Comparison

| Feature | Gorilla Mux | Gin | Echo |
|---------|-------------|-----|------|
| **Performance** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Learning Curve** | Easy | Easy | Easy |
| **Standard Compatibility** | Perfect | Medium | Medium |
| **URL Patterns** | Very Powerful | Basic | Basic |
| **Middleware** | Basic | Rich | Rich |
| **Data Binding** | ❌ | ✅ | ✅ |
| **Community** | Active | Very Active | Active |
| **Regular Expressions** | ✅ | ❌ | ❌ |

## Notes

- Gorilla Mux uses the standard `http.Handler` interface
- Access path variables with `mux.Vars(r)`
- `http.ErrServerClosed` is considered normal shutdown and won't trigger the error callback
- Cannot start multiple servers on the same port
- Path prefix must end with a slash (/)

## References

- [Gorilla Mux Official Documentation](https://github.com/gorilla/mux)
- [Go net/http Documentation](https://pkg.go.dev/net/http)
- [Gorilla Toolkit](https://www.gorillatoolkit.org/)

## License

This package follows the project's license.
