# HTTP

Simplified HTTP client and server utilities for Go.

## Overview

The http package provides convenient wrapper functions around Go's standard `net/http` and `gorilla/mux` packages. It simplifies making HTTP requests and managing HTTP servers with powerful routing capabilities while reducing boilerplate code.

## Features

- **HTTP Client** - Simplified request function with built-in authentication
- **HTTP Server** - Easy server setup with gorilla/mux routing
- **Handler Registration** - Multiple registration methods for different use cases
- **Path Prefix Routing** - Route groups and static file serving
- **Middleware Support** - Chain middleware functions
- **Graceful Shutdown** - Proper server lifecycle management
- **Custom Transport** - Configure HTTP client transport settings

## Installation

```bash
go get -u github.com/common-library/go/http
go get -u github.com/gorilla/mux
```

## Quick Start

### HTTP Client

```go
import "github.com/common-library/go/http"

func main() {
    resp, err := http.Request(
        "https://api.example.com/users",
        http.MethodGet,
        nil,  // headers
        "",   // body
        10,   // timeout in seconds
        "",   // username
        "",   // password
        nil,  // transport
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Status: %d\n", resp.StatusCode)
    fmt.Printf("Body: %s\n", resp.Body)
}
```

### HTTP Server

```go
import (
    "net/http"
    "github.com/common-library/go/http"
)

func main() {
    var server http.Server
    
    // Register handlers
    server.RegisterHandlerFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    server.RegisterHandlerFunc("/api/users", usersHandler, http.MethodGet)
    
    // Start server
    err := server.Start(":8080", func(err error) {
        log.Fatal(err)
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Keep running
    select {}
}
```

## HTTP Client

### Request Function

```go
func Request(
    url string,
    method string,
    header map[string][]string,
    body string,
    timeout time.Duration,
    username string,
    password string,
    transport *http.Transport,
) (Response, error)
```

Performs an HTTP request and returns the complete response.

**Parameters:**
- `url` - Target URL
- `method` - HTTP method (http.MethodGet, http.MethodPost, etc.)
- `header` - HTTP headers as map[string][]string
- `body` - Request body as string
- `timeout` - Request timeout duration (e.g., 10*time.Second)
- `username` - Username for Basic Auth (empty for no auth)
- `password` - Password for Basic Auth (empty for no auth)
- `transport` - Custom HTTP transport (nil for default)

**Returns:**
- `Response` - Response struct with Header, Body, and StatusCode
- `error` - Error if request fails

### Response Type

```go
type Response struct {
    Header     http.Header
    Body       string
    StatusCode int
}
```

### Client Examples

#### Simple GET Request

```go
resp, err := http.Request(
    "https://api.example.com/users",
    http.MethodGet,
    nil,
    "",
    10*time.Second,
    "", "",
    nil,
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %d\n", resp.StatusCode)
fmt.Printf("Body: %s\n", resp.Body)
```

#### POST Request with JSON

```go
headers := map[string][]string{
    "Content-Type": {"application/json"},
}

body := `{
    "name": "Alice",
    "email": "alice@example.com"
}`

resp, err := http.Request(
    "https://api.example.com/users",
    http.MethodPost,
    headers,
    body,
    30*time.Second,
    "", "",
    nil,
)

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created user, status: %d\n", resp.StatusCode)
```

#### Request with Authentication

```go
resp, err := http.Request(
    "https://api.example.com/admin/stats",
    http.MethodGet,
    nil,
    "",
    10*time.Second,
    "admin",      // username
    "password123", // password
    nil,
)

if err != nil {
    log.Fatal(err)
}

if resp.StatusCode == 401 {
    log.Println("Authentication failed")
} else {
    fmt.Println(resp.Body)
}
```

#### Request with Custom Headers

```go
headers := map[string][]string{
    "Authorization": {"Bearer token123"},
    "X-API-Key":     {"secret456"},
    "Accept":        {"application/json"},
}

resp, err := http.Request(
    "https://api.example.com/data",
    http.MethodGet,
    headers,
    "",
    15*time.Second,
    "", "",
    nil,
)
```

#### Request with Custom Transport

```go
import "crypto/tls"

// Skip TLS verification (not recommended for production)
transport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

resp, err := http.Request(
    "https://self-signed.example.com/api",
    http.MethodGet,
    nil,
    "",
    10*time.Second,
    "", "",
    transport,
)
```

## HTTP Server

### Server Type

```go
type Server struct {
    // Contains unexported fields
}
```

### Handler Registration Methods

#### RegisterHandler

```go
func (s *Server) RegisterHandler(path string, handler http.Handler, methods ...string)
```

Registers an HTTP handler for a specific path.

**Example:**

```go
type MyHandler struct{}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello from handler"))
}

var server http.Server

// All methods
server.RegisterHandler("/api", &MyHandler{})

// GET only
server.RegisterHandler("/users", &MyHandler{}, http.MethodGet)

// GET and POST
server.RegisterHandler("/items", &MyHandler{}, http.MethodGet, http.MethodPost)
```

#### RegisterHandlerFunc

```go
func (s *Server) RegisterHandlerFunc(path string, handlerFunc http.HandlerFunc, methods ...string)
```

Registers an HTTP handler function for a specific path.

**Example:**

```go
var server http.Server

// All methods
server.RegisterHandlerFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("pong"))
})

// GET only
server.RegisterHandlerFunc("/users", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`[{"id": 1, "name": "Alice"}]`))
}, http.MethodGet)

// POST only
server.RegisterHandlerFunc("/users", func(w http.ResponseWriter, r *http.Request) {
    // Create user
    w.WriteHeader(http.StatusCreated)
}, http.MethodPost)
```

#### RegisterPathPrefixHandler

```go
func (s *Server) RegisterPathPrefixHandler(prefix string, handler http.Handler, methods ...string)
```

Registers a handler for all paths matching the prefix.

**Example:**

```go
// Serve static files
fileServer := http.FileServer(http.Dir("./static"))
server.RegisterPathPrefixHandler("/static/", http.StripPrefix("/static/", fileServer))

// API v1 endpoints
server.RegisterPathPrefixHandler("/api/v1/", apiV1Handler, http.MethodGet)
```

#### RegisterPathPrefixHandlerFunc

```go
func (s *Server) RegisterPathPrefixHandlerFunc(prefix string, handlerFunc http.HandlerFunc, methods ...string)
```

Registers a handler function for paths matching the prefix.

**Example:**

```go
server.RegisterPathPrefixHandlerFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"status": "ok"}`))
})
```

### Server Lifecycle

#### Start

```go
func (s *Server) Start(address string, listenAndServeFailureFunc func(err error)) error
```

Starts the HTTP server.

**Example:**

```go
var server http.Server

// Simple start
err := server.Start(":8080", func(err error) {
    log.Fatal(err)
})
```

#### Use

```go
func (s *Server) Use(middleware ...mux.MiddlewareFunc)
```

Registers global middleware.

**Example:**

```go
var server http.Server

loggingMiddleware := func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

server.Use(loggingMiddleware)
server.Start(":8080", nil)
```

#### Stop

```go
func (s *Server) Stop(shutdownTimeout time.Duration) error
```

Gracefully shuts down the server.

**Example:**

```go
// Shutdown with 30 second timeout
err := server.Stop(30 * time.Second)
if err != nil {
    log.Printf("Shutdown error: %v", err)
}
```

#### IsRunning

```go
func (s *Server) IsRunning() bool
```

Returns whether the server is currently running.

**Example:**

```go
if server.IsRunning() {
    log.Println("Server is running")
}
```

#### GetRouter

```go
func (s *Server) GetRouter() *mux.Router
```

Returns the gorilla/mux router instance for advanced configuration.

**Example:**

```go
router := server.GetRouter()
router.StrictSlash(true)
router.HandleFunc("/", homeHandler)
```

#### SetRouter

```go
func (s *Server) SetRouter(router *mux.Router)
```

Sets a custom router.

**Example:**

```go
router := mux.NewRouter()
router.StrictSlash(true)
router.HandleFunc("/", homeHandler)

var server http.Server
server.SetRouter(router)
server.Start(":8080", nil)
```

## Complete Examples

### REST API Server

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    
    httplib "github.com/common-library/go/http"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var users = []User{
    {ID: 1, Name: "Alice"},
    {ID: 2, Name: "Bob"},
}

func main() {
    var server httplib.Server
    
    // GET /api/users - List all users
    server.RegisterHandlerFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    }, http.MethodGet)
    
    // POST /api/users - Create user
    server.RegisterHandlerFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        var user User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        user.ID = len(users) + 1
        users = append(users, user)
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
    }, http.MethodPost)
    
    log.Println("Server starting on :8080")
    server.Start(":8080", func(err error) {
        log.Fatal(err)
    })
    
    select {}
}
```

### Static File Server

```go
package main

import (
    "log"
    "net/http"
    
    httplib "github.com/common-library/go/http"
)

func main() {
    var server httplib.Server
    
    // Serve static files from ./public directory
    fileServer := http.FileServer(http.Dir("./public"))
    server.RegisterPathPrefixHandler("/", fileServer)
    
    log.Println("Static file server on :8080")
    server.Start(":8080", nil)
    
    select {}
}
```

### Server with Middleware

```go
package main

import (
    "log"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
    httplib "github.com/common-library/go/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)
        
        next.ServeHTTP(w, r)
        
        log.Printf("Completed in %v", time.Since(start))
    })
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        if apiKey != "secret123" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func main() {
    var server httplib.Server
    
    server.RegisterHandlerFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Authenticated data"))
    })
    
    server.Use(loggingMiddleware)
    server.Use(authMiddleware)
    server.Start(":8080", nil)
    
    select {}
}
```

### Server with Graceful Shutdown

```go
package main

import (
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    
    httplib "github.com/common-library/go/http"
)

func main() {
    var server httplib.Server
    
    server.RegisterHandlerFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    // Start server
    log.Println("Starting server on :8080")
    server.Start(":8080", func(err error) {
        log.Printf("Server error: %v", err)
    })
    
    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    
    log.Println("Shutting down server...")
    if err := server.Stop(30 * time.Second); err != nil {
        log.Printf("Shutdown error: %v", err)
    }
    
    log.Println("Server stopped gracefully")
}
```

### API Client Wrapper

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    
    httplib "github.com/common-library/go/http"
)

type APIClient struct {
    baseURL string
    apiKey  string
}

func NewAPIClient(baseURL, apiKey string) *APIClient {
    return &APIClient{
        baseURL: baseURL,
        apiKey:  apiKey,
    }
}

func (c *APIClient) GetUsers() ([]User, error) {
    headers := map[string][]string{
        "Authorization": {fmt.Sprintf("Bearer %s", c.apiKey)},
        "Accept":        {"application/json"},
    }
    
    resp, err := httplib.Request(
        c.baseURL+"/api/users",
        http.MethodGet,
        headers,
        "",
        30*time.Second,
        "", "",
        nil,
    )
    if err != nil {
        return nil, err
    }
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("API error: %d", resp.StatusCode)
    }
    
    var users []User
    if err := json.Unmarshal([]byte(resp.Body), &users); err != nil {
        return nil, err
    }
    
    return users, nil
}

func (c *APIClient) CreateUser(user User) error {
    headers := map[string][]string{
        "Authorization": {fmt.Sprintf("Bearer %s", c.apiKey)},
        "Content-Type":  {"application/json"},
    }
    
    body, _ := json.Marshal(user)
    
    resp, err := httplib.Request(
        c.baseURL+"/api/users",
        http.MethodPost,
        headers,
        string(body),
        30*time.Second,
        "", "",
        nil,
    )
    if err != nil {
        return err
    }
    
    if resp.StatusCode != 201 {
        return fmt.Errorf("API error: %d - %s", resp.StatusCode, resp.Body)
    }
    
    return nil
}

func main() {
    client := NewAPIClient("https://api.example.com", "secret123")
    
    users, err := client.GetUsers()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, user := range users {
        fmt.Printf("User: %+v\n", user)
    }
}
```

## Best Practices

### 1. Always Set Timeouts

```go
// Good: Set reasonable timeout
resp, err := http.Request(url, http.MethodGet, nil, "", 30*time.Second, "", "", nil)

// Avoid: Very long or no timeout
// May hang indefinitely
```

### 2. Handle Response Status Codes

```go
// Good: Check status code
resp, err := http.Request(url, http.MethodGet, nil, "", 10, "", "", nil)
if err != nil {
    log.Fatal(err)
}

switch resp.StatusCode {
case 200:
    fmt.Println("Success:", resp.Body)
case 404:
    fmt.Println("Not found")
case 500:
    fmt.Println("Server error")
default:
    fmt.Printf("Unexpected status: %d\n", resp.StatusCode)
}

// Avoid: Ignoring status code
// resp, _ := http.Request(...)
// fmt.Println(resp.Body) // May not be what you expect
```

### 3. Use Proper HTTP Methods

```go
// Good: Use correct method constants
http.Request(url, http.MethodGet, ...)
http.Request(url, http.MethodPost, ...)
http.Request(url, http.MethodPut, ...)

// Avoid: String literals
// http.Request(url, "GET", ...)
```

### 4. Set Content-Type Headers

```go
// Good: Set Content-Type for POST/PUT
headers := map[string][]string{
    "Content-Type": {"application/json"},
}

resp, _ := http.Request(url, http.MethodPost, headers, jsonBody, 10*time.Second, "", "", nil)

// Avoid: Missing Content-Type
// Server may not parse body correctly
```

### 5. Implement Graceful Shutdown

```go
// Good: Graceful shutdown
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
<-sigChan

server.Stop(30) // Wait for active connections

// Avoid: Abrupt termination
// os.Exit(0) // May interrupt active requests
```

### 6. Use Middleware for Cross-Cutting Concerns

```go
// Good: Centralize logging, auth, CORS in middleware
server.Use(loggingMiddleware)
server.Use(authMiddleware)
server.Use(corsMiddleware)
server.Start(":8080", nil)

// Avoid: Repeating logic in every handler
// Each handler implements its own logging
```

## Error Handling

### Client Errors

```go
resp, err := http.Request(url, http.MethodGet, nil, "", 10, "", "", nil)
if err != nil {
    // Network error, timeout, DNS failure, etc.
    log.Printf("Request failed: %v", err)
    return
}

// Check HTTP status
if resp.StatusCode >= 400 {
    log.Printf("HTTP error: %d - %s", resp.StatusCode, resp.Body)
}
```

### Server Errors

```go
server.RegisterHandlerFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
    data, err := fetchData()
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        log.Printf("Error fetching data: %v", err)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
})
```

## Performance Tips

1. **Connection Reuse** - HTTP client reuses connections automatically
2. **Timeouts** - Set appropriate timeouts to prevent hanging requests
3. **Custom Transport** - Use transport for connection pooling configuration
4. **Middleware Ordering** - Order middleware from general to specific
5. **Graceful Shutdown** - Allow active requests to complete

## Testing

### Unit Testing HTTP Handlers

```go
func TestUserHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
    rec := httptest.NewRecorder()
    
    userHandler(rec, req)
    
    if rec.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", rec.Code)
    }
}
```

### Integration Testing

```go
func TestServer(t *testing.T) {
    var server httplib.Server
    
    server.RegisterHandlerFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("test response"))
    })
    
    server.Start(":8081", nil)
    defer server.Stop(5)
    
    time.Sleep(100 * time.Millisecond) // Wait for server
    
    resp, err := httplib.Request(
        "http://localhost:8081/test",
        http.MethodGet,
        nil, "", 5, "", "", nil,
    )
    
    if err != nil {
        t.Fatal(err)
    }
    
    if resp.StatusCode != 200 {
        t.Errorf("Expected 200, got %d", resp.StatusCode)
    }
    
    if resp.Body != "test response" {
        t.Errorf("Unexpected body: %s", resp.Body)
    }
}
```

## Dependencies

- `net/http` - Go standard library
- `github.com/gorilla/mux` - HTTP router and dispatcher

## Further Reading

- [net/http documentation](https://pkg.go.dev/net/http)
- [gorilla/mux documentation](https://github.com/gorilla/mux)
- [HTTP Protocol](https://developer.mozilla.org/en-US/docs/Web/HTTP)
- [RESTful API Design](https://restfulapi.net/)
