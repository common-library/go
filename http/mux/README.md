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

## Mux-Specific Features

### Path Prefix Routing

Unlike Echo and Gin frameworks, this Gorilla Mux wrapper provides dedicated **PathPrefix** methods:

- **Why PathPrefix?**: Gorilla Mux's `PathPrefix()` is optimized for serving static files, API versioning, and mounting sub-applications
- **Echo/Gin Alternative**: These frameworks use `Group()` for similar purposes, but Mux's PathPrefix is more direct
- **Performance**: PathPrefix matching is more efficient for prefix-based routing than pattern matching

**Use Cases**:
- **Static File Serving**: `/static/`, `/assets/`, `/public/`
- **API Versioning**: `/api/v1/`, `/api/v2/`
- **Swagger/Docs**: `/swagger/`, `/docs/`
- **Sub-applications**: `/admin/`, `/dashboard/`

**Example Comparison**:
```go
// Mux (Direct PathPrefix)
server.RegisterPathPrefixHandlerAny("/static/", http.FileServer(http.Dir("./static")))

// Echo/Gin (Using Group)
// Echo: g := e.Group("/static"); g.Use(middleware.Static("./static"))
// Gin:  r.Static("/static", "./static")
```

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

// Handler registration with specific method (echo/gin style)
handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("handler"))
})
server.RegisterHandler(http.MethodGet, "/api", handler)

// HandlerFunc registration with specific method
server.RegisterHandlerFunc(http.MethodGet, "/users", getUsersHandler)
server.RegisterHandlerFunc(http.MethodPost, "/users", createUserHandler)
server.RegisterHandlerFunc(http.MethodPut, "/users/{id}", updateUserHandler)
server.RegisterHandlerFunc(http.MethodDelete, "/users/{id}", deleteUserHandler)

// Register for all HTTP methods
server.RegisterHandlerAny("/webhook", webhookHandler)
server.RegisterHandlerFuncAny("/catch-all", catchAllHandler)
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

Path prefix routing is a **Mux-specific feature** for efficiently handling requests based on URL path prefixes.

```go
server := &mux.Server{}

// Serve static files (all methods)
// Common use case: CSS, JS, images, downloads
fileServer := http.FileServer(http.Dir("./static"))
server.RegisterPathPrefixHandlerAny("/static/", http.StripPrefix("/static/", fileServer))

// API version routing with specific method
// Useful for API versioning without duplicating handlers
server.RegisterPathPrefixHandlerFunc(http.MethodGet, "/api/v1/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("API v1"))
})

server.RegisterPathPrefixHandlerFunc(http.MethodGet, "/api/v2/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("API v2"))
})

// Admin pages with GET method
// All /admin/* routes handled by one handler
server.RegisterPathPrefixHandler(http.MethodGet, "/admin/", adminHandler)

// Swagger/Documentation (real-world example)
server.RegisterPathPrefixHandlerAny("/swagger/", httpSwagger.WrapHandler)
```

**Key Differences from Echo/Gin:**
- Echo/Gin use `Group()` for path grouping and middleware scoping
- Mux's PathPrefix is optimized for prefix matching, not grouping
- PathPrefix matches **any** path starting with the prefix
- Ideal for file servers, API versioning, and catch-all routes

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
server.RegisterHandlerFuncAny("/", homeHandler)

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
server.RegisterHandlerFuncAny("/health", func(w http.ResponseWriter, r *http.Request) {
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

#### `RegisterHandler(method, path string, handler http.Handler)`
Registers an HTTP handler for a specific method.

- **method**: HTTP method (http.MethodGet, http.MethodPost, etc.)
- **path**: URL path pattern (supports path variables: `/users/{id}`)
- **handler**: HTTP Handler

#### `RegisterHandlerAny(path string, handler http.Handler)`
Registers an HTTP handler for all methods.

- **path**: URL path pattern
- **handler**: HTTP Handler

#### `RegisterHandlerFunc(method, path string, handlerFunc http.HandlerFunc)`
Registers an HTTP handler function for a specific method.

- **method**: HTTP method (http.MethodGet, http.MethodPost, etc.)
- **path**: URL path pattern
- **handlerFunc**: HTTP HandlerFunc

#### `RegisterHandlerFuncAny(path string, handlerFunc http.HandlerFunc)`
Registers an HTTP handler function for all methods.

- **path**: URL path pattern
- **handlerFunc**: HTTP HandlerFunc

#### `RegisterPathPrefixHandler(method, prefix string, handler http.Handler)`
**[Mux-Specific]** Registers a handler for a path prefix with specific method.

This is a unique feature of Gorilla Mux not available in Echo/Gin wrappers. Use this for serving static files, API versioning, or mounting sub-applications.

- **method**: HTTP method (http.MethodGet, http.MethodPost, etc.)
- **prefix**: Path prefix (e.g., `/api/`, `/static/`)
- **handler**: HTTP Handler

**Example**: `server.RegisterPathPrefixHandler(http.MethodGet, "/static/", fileServer)`

#### `RegisterPathPrefixHandlerAny(prefix string, handler http.Handler)`
**[Mux-Specific]** Registers a handler for a path prefix for all methods.

- **prefix**: Path prefix
- **handler**: HTTP Handler

**Common Use Cases**:
- Static file serving: `/static/`, `/assets/`
- Swagger UI: `/swagger/`, `/docs/`
- CloudEvents receiver: `/` (root)

**Example**: `server.RegisterPathPrefixHandlerAny("/swagger/", httpSwagger.WrapHandler)`

#### `RegisterPathPrefixHandlerFunc(method, prefix string, handlerFunc http.HandlerFunc)`
**[Mux-Specific]** Registers a handler function for a path prefix with specific method.

- **method**: HTTP method (http.MethodGet, http.MethodPost, etc.)
- **prefix**: Path prefix
- **handlerFunc**: HTTP HandlerFunc

**Example**: `server.RegisterPathPrefixHandlerFunc(http.MethodGet, "/api/v1/", v1Handler)`

#### `RegisterPathPrefixHandlerFuncAny(prefix string, handlerFunc http.HandlerFunc)`
**[Mux-Specific]** Registers a handler function for a path prefix for all methods.

- **prefix**: Path prefix
- **handlerFunc**: HTTP HandlerFunc

**Example**: `server.RegisterPathPrefixHandlerFuncAny("/api/", apiHandler)`

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
| **Path Prefix** | ✅ Dedicated | ❌ Use Group | ❌ Use Group |
| **Subrouters** | ✅ | ✅ Group | ✅ Group |

### When to Use Each Framework

**Choose Gorilla Mux when:**
- You need standard `net/http` compatibility
- Complex URL patterns with regex are required
- Path prefix routing is a key requirement
- You prefer explicit over implicit
- Building API gateways or reverse proxies

**Choose Gin when:**
- Maximum performance is critical
- You need built-in data binding/validation
- JSON-heavy API development
- Middleware ecosystem is important

**Choose Echo when:**
- Balance between performance and features
- Middleware-centric architecture
- WebSocket support needed
- Template rendering required

## Notes

- Gorilla Mux uses the standard `http.Handler` interface
- Access path variables with `mux.Vars(r)`
- `http.ErrServerClosed` is considered normal shutdown and won't trigger the error callback
- Cannot start multiple servers on the same port
- **Path prefix must end with a slash (/)** - This is critical for correct routing
- PathPrefix methods are Mux-specific features not available in Echo/Gin wrappers
- For route grouping with shared middleware, use `GetRouter().PathPrefix().Subrouter()`

## References

- [Gorilla Mux Official Documentation](https://github.com/gorilla/mux)
- [Go net/http Documentation](https://pkg.go.dev/net/http)
- [Gorilla Toolkit](https://www.gorillatoolkit.org/)

## License

This package follows the project's license.
