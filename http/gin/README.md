# Gin HTTP Server Package

Gin web framework wrapper package for Go.

## Features

- **Based on Gin**: High-performance web framework
- **Easy routing**: GET, POST, PUT, DELETE, PATCH, Any methods
- **Group routing**: URL path grouping support
- **Middleware**: Global and route-specific middleware
- **Graceful Shutdown**: Timeout-based safe shutdown
- **Thread-safe**: Concurrency safety guaranteed

## Installation

```bash
go get github.com/common-library/go/http/gin
go get github.com/gin-gonic/gin
```

## When to Choose Gin

- ✅ Ultra-high performance API servers needed
- ✅ REST APIs with heavy JSON processing
- ✅ Preference for simple API structure
- ✅ Fast development speed important
- ✅ Active community support needed

## Usage Examples

### Basic Server

```go
package main

import (
    "net/http"
    "github.com/common-library/go/http/gin"
    "github.com/gin-gonic/gin"
)

func main() {
    server := &gin.Server{}
    
    // Simple handler
    server.RegisterHandler(http.MethodGet, "/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, World!")
    })
    
    // JSON response
    server.RegisterHandler(http.MethodGet, "/api/hello", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Hello from Gin",
        })
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
server := &gin.Server{}

// Basic registration
server.RegisterHandler(http.MethodGet, "/users", getUsersHandler)
server.RegisterHandler(http.MethodPost, "/users", createUserHandler)
server.RegisterHandler(http.MethodPut, "/users/:id", updateUserHandler)
server.RegisterHandler(http.MethodDelete, "/users/:id", deleteUserHandler)
server.RegisterHandler(http.MethodPatch, "/users/:id", patchUserHandler)

// Register multiple methods at once
server.RegisterHandlerMethods([]string{http.MethodGet, http.MethodPost}, "/multi", multiHandler)

// Allow all methods
server.RegisterHandlerAny("/any", anyHandler)

// Register with middleware
server.RegisterHandler(http.MethodGet, "/admin", adminHandler, authMiddleware, logMiddleware)
```

### Group Routing

```go
server := &gin.Server{}

// Create API group
api := server.Group("/api")

// Routes within group
api.GET("/users", func(c *gin.Context) {
    c.JSON(http.StatusOK, users)
})

api.POST("/users", func(c *gin.Context) {
    // Create user logic
    c.JSON(http.StatusCreated, newUser)
})

// Nested groups
v1 := api.Group("/v1")
v1.GET("/info", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"version": "1.0"})
})

// Group-specific middleware
admin := api.Group("/admin", adminAuthMiddleware)
admin.GET("/stats", getStatsHandler)
```

### Middleware Usage

```go
import "github.com/gin-gonic/gin"

server := &gin.Server{}

// Global middleware
server.Use(gin.Logger())
server.Use(gin.Recovery())

// Custom middleware
authMiddleware := func(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
        return
    }
    c.Next()
}

server.RegisterHandler(http.MethodGet, "/protected", protectedHandler, authMiddleware)

// CORS middleware
corsMiddleware := func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
    c.Next()
}
server.Use(corsMiddleware)
```

### Request Processing Examples

```go
// Path parameters
server.RegisterHandler(http.MethodGet, "/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.String(http.StatusOK, "User ID: "+id)
})

// Query parameters
server.RegisterHandler(http.MethodGet, "/search", func(c *gin.Context) {
    query := c.Query("q")
    page := c.DefaultQuery("page", "1")
    c.JSON(http.StatusOK, gin.H{
        "query": query,
        "page":  page,
    })
})

// JSON binding
type User struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

server.RegisterHandler(http.MethodPost, "/users", func(c *gin.Context) {
    var u User
    if err := c.ShouldBindJSON(&u); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, u)
})

// Form data
server.RegisterHandler(http.MethodPost, "/form", func(c *gin.Context) {
    name := c.PostForm("name")
    email := c.DefaultPostForm("email", "default@example.com")
    c.JSON(http.StatusOK, gin.H{
        "name":  name,
        "email": email,
    })
})
```

### Graceful Shutdown

```go
import (
    "os"
    "os/signal"
    "syscall"
    "time"
)

server := &gin.Server{}
server.RegisterHandler(http.MethodGet, "/", homeHandler)

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

### Direct Gin Engine Access

```go
server := &gin.Server{}

// Get Gin engine
engine := server.GetEngine()

// Advanced settings
engine.MaxMultipartMemory = 8 << 20 // 8 MiB

// Serve static files
engine.Static("/static", "./public")
engine.StaticFile("/favicon.ico", "./resources/favicon.ico")

// HTML templates
engine.LoadHTMLGlob("templates/*")
engine.GET("/index", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{
        "title": "Home",
    })
})

// File upload
engine.POST("/upload", func(c *gin.Context) {
    file, _ := c.FormFile("file")
    c.SaveUploadedFile(file, "./uploads/"+file.Filename)
    c.String(http.StatusOK, "Uploaded: "+file.Filename)
})
```

### Server Status Check

```go
server := &gin.Server{}

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
server.RegisterHandler(http.MethodGet, "/health", func(c *gin.Context) {
    status := "down"
    if server.IsRunning() {
        status = "up"
    }
    c.JSON(http.StatusOK, gin.H{"status": status})
})
```

## API Documentation

### Server

#### `RegisterHandler(method, path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc)`
Registers an HTTP handler.

- **method**: HTTP method (http.MethodGet, http.MethodPost, etc.)
- **path**: URL path pattern
- **handler**: Gin handler function
- **middleware**: Optional middleware

#### `RegisterHandlerMethods(methods []string, path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc)`
Registers a handler for multiple HTTP methods.

- **methods**: Array of HTTP methods
- **path**: URL path pattern
- **handler**: Gin handler function
- **middleware**: Optional middleware

#### `RegisterHandlerAny(path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc)`
Registers a handler for all HTTP methods.

- **path**: URL path pattern
- **handler**: Gin handler function
- **middleware**: Optional middleware

#### `Use(middleware ...gin.HandlerFunc)`
Registers global middleware.

#### `Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup`
Creates a router group.

- **relativePath**: Group path prefix
- **handlers**: Group middleware
- **return**: Gin RouterGroup instance

#### `Start(address string, listenAndServeFailureFunc func(err error)) error`
Starts the server.

- **address**: Server address (e.g., ":8080")
- **listenAndServeFailureFunc**: Error callback (excludes http.ErrServerClosed)

#### `Stop(shutdownTimeout time.Duration) error`
Safely shuts down the server.

- **shutdownTimeout**: Shutdown wait time

#### `IsRunning() bool`
Returns the server running status.

#### `GetEngine() *gin.Engine`
Provides direct access to the Gin engine.

## Framework Comparison

| Feature | Gin | Echo | Gorilla Mux |
|---------|-----|------|-------------|
| **Performance** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Learning Curve** | Easy | Easy | Medium |
| **Middleware** | Rich | Rich | Basic |
| **Data Binding** | ✅ | ✅ | ❌ |
| **Validation** | ✅ | ✅ | ❌ |
| **JSON Performance** | Best | Very Good | Average |
| **Community** | Very Active | Active | Active |
| **GitHub Stars** | 78k+ | 29k+ | 20k+ |

## Notes

- Gin handlers have no return value (`void`)
- Control middleware chain with `c.Next()` or `c.Abort()`
- `http.ErrServerClosed` is considered normal shutdown and won't trigger the error callback
- Cannot start multiple servers on the same port
- For production environments, recommend setting `gin.SetMode(gin.ReleaseMode)`

## References

- [Gin Official Documentation](https://gin-gonic.com/)
- [Gin GitHub](https://github.com/gin-gonic/gin)
- [Gin Examples](https://github.com/gin-gonic/examples)

## License

This package follows the project's license.
