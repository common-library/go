# Echo HTTP Server Package

Echo web framework wrapper package for Go.

## Features

- **Based on Echo v4**: High-performance web framework
- **Easy routing**: GET, POST, PUT, DELETE, PATCH methods
- **Group routing**: URL path grouping support
- **Middleware**: Global and route-specific middleware
- **Graceful Shutdown**: Timeout-based safe shutdown
- **Thread-safe**: Concurrency safety guaranteed

## Installation

```bash
go get github.com/common-library/go/http/echo
go get github.com/labstack/echo/v4
```

## When to Choose Echo

- ✅ API servers where high performance is critical
- ✅ When middleware chains are needed
- ✅ Data binding and validation features required
- ✅ WebSocket support needed
- ✅ Preference for concise APIs

## Usage Examples

### Basic Server

```go
package main

import (
    "net/http"
    "github.com/common-library/go/http/echo"
    echo_lib "github.com/labstack/echo/v4"
)

func main() {
    server := &echo.Server{}
    
    // Simple handler
    server.RegisterHandler(echo_lib.GET, "/", func(c echo_lib.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    
    // JSON response
    server.RegisterHandler(echo_lib.GET, "/api/hello", func(c echo_lib.Context) error {
        return c.JSON(http.StatusOK, map[string]string{
            "message": "Hello from Echo",
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
server := &echo.Server{}

// Basic registration
server.RegisterHandler(echo_lib.GET, "/users", getUsersHandler)
server.RegisterHandler(echo_lib.POST, "/users", createUserHandler)
server.RegisterHandler(echo_lib.PUT, "/users/:id", updateUserHandler)
server.RegisterHandler(echo_lib.DELETE, "/users/:id", deleteUserHandler)
server.RegisterHandler(echo_lib.PATCH, "/users/:id", patchUserHandler)

// Register with middleware
server.RegisterHandler(echo_lib.GET, "/admin", adminHandler, authMiddleware, logMiddleware)
```

### Group Routing

```go
server := &echo.Server{}

// Create API group
api := server.Group("/api")

// Routes within group
api.GET("/users", func(c echo_lib.Context) error {
    return c.JSON(http.StatusOK, users)
})

api.POST("/users", func(c echo_lib.Context) error {
    // Create user logic
    return c.JSON(http.StatusCreated, newUser)
})

// Nested groups
v1 := api.Group("/v1")
v1.GET("/info", func(c echo_lib.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"version": "1.0"})
})

// Group-specific middleware
admin := api.Group("/admin", adminAuthMiddleware)
admin.GET("/stats", getStatsHandler)
```

### Middleware Usage

```go
import "github.com/labstack/echo/v4/middleware"

server := &echo.Server{}

// Global middleware
server.Use(middleware.Logger())
server.Use(middleware.Recover())
server.Use(middleware.CORS())

// Route-specific middleware
authMiddleware := func(next echo_lib.HandlerFunc) echo_lib.HandlerFunc {
    return func(c echo_lib.Context) error {
        token := c.Request().Header.Get("Authorization")
        if token == "" {
            return echo_lib.NewHTTPError(http.StatusUnauthorized, "missing token")
        }
        return next(c)
    }
}

server.RegisterHandler(echo_lib.GET, "/protected", protectedHandler, authMiddleware)
```

### Request Processing Examples

```go
// Path parameters
server.RegisterHandler(echo_lib.GET, "/users/:id", func(c echo_lib.Context) error {
    id := c.Param("id")
    return c.String(http.StatusOK, "User ID: "+id)
})

// Query parameters
server.RegisterHandler(echo_lib.GET, "/search", func(c echo_lib.Context) error {
    query := c.QueryParam("q")
    page := c.QueryParam("page")
    return c.JSON(http.StatusOK, map[string]string{
        "query": query,
        "page":  page,
    })
})

// JSON binding
type User struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

server.RegisterHandler(echo_lib.POST, "/users", func(c echo_lib.Context) error {
    u := new(User)
    if err := c.Bind(u); err != nil {
        return err
    }
    if err := c.Validate(u); err != nil {
        return err
    }
    return c.JSON(http.StatusCreated, u)
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

server := &echo.Server{}
server.RegisterHandler(echo_lib.GET, "/", homeHandler)

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

### Direct Echo Instance Access

```go
server := &echo.Server{}

// Get Echo instance
e := server.GetEcho()

// Advanced settings
e.HideBanner = true
e.Debug = false
e.Server.ReadTimeout = 30 * time.Second
e.Server.WriteTimeout = 30 * time.Second

// Serve static files
e.Static("/static", "public")

// File upload
e.POST("/upload", uploadHandler)
```

## API Documentation

### Server

#### `RegisterHandler(method, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc)`
Registers an HTTP handler.

- **method**: HTTP method (echo.GET, echo.POST, etc.)
- **path**: URL path pattern
- **handler**: Echo handler function
- **middleware**: Optional middleware

#### `Use(middleware ...echo.MiddlewareFunc)`
Registers global middleware.

#### `Group(prefix string, middleware ...echo.MiddlewareFunc) *echo.Group`
Creates a router group.

- **prefix**: Group path prefix
- **middleware**: Group middleware
- **return**: Echo Group instance

#### `Start(address string, listenAndServeFailureFunc func(err error)) error`
Starts the server.

- **address**: Server address (e.g., ":8080")
- **listenAndServeFailureFunc**: Error callback (excludes http.ErrServerClosed)

#### `Stop(shutdownTimeout time.Duration) error`
Safely shuts down the server.

- **shutdownTimeout**: Shutdown wait time

#### `IsRunning() bool`
Returns the server running status.

#### `GetEcho() *echo.Echo`
Provides direct access to the Echo instance.

## Framework Comparison

| Feature | Echo | Gin | Gorilla Mux |
|---------|------|-----|-------------|
| **Performance** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Learning Curve** | Easy | Easy | Medium |
| **Middleware** | Rich | Rich | Basic |
| **Data Binding** | ✅ | ✅ | ❌ |
| **Validation** | ✅ | ✅ | ❌ |
| **WebSocket** | ✅ | ✅ | ❌ |
| **Community** | Active | Very Active | Active |

## Notes

- Echo handlers must return `error`
- Middleware order matters (Logger → Recover sequence recommended)
- `http.ErrServerClosed` is considered normal shutdown and won't trigger the error callback
- Cannot start multiple servers on the same port

## References

- [Echo Official Documentation](https://echo.labstack.com/)
- [Echo GitHub](https://github.com/labstack/echo)
- [Echo Guide](https://echo.labstack.com/guide/)

## License

This package follows the project's license.
