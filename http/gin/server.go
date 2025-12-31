// Package gin provides a simplified wrapper around the Gin web framework.
//
// This package offers convenient methods for creating and managing Gin HTTP servers
// with support for routing, middleware, and graceful shutdown.
//
// Features:
//   - Simplified server setup with Gin
//   - Handler and middleware registration
//   - Group routing support
//   - Graceful shutdown with timeout
//   - Thread-safe operations
//
// Example:
//
//	server := &gin.Server{}
//	server.RegisterHandler(http.MethodGet, "/hello", func(c *gin.Context) {
//	    c.String(200, "Hello, World!")
//	})
//	server.Start(":8080", nil)
package gin

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// Server is a struct that provides Gin HTTP server related methods.
type Server struct {
	mutex   sync.Mutex
	running atomic.Bool

	server *http.Server
	engine *gin.Engine
}

// RegisterHandler registers an HTTP handler with optional middleware.
//
// Parameters:
//   - method: HTTP method (http.MethodGet, http.MethodPost, etc.)
//   - relativePath: URL path pattern
//   - handlers: Gin handler functions (handler + optional middleware)
//
// Example:
//
//	server.RegisterHandler(http.MethodGet, "/users/:id", getUserHandler)
func (s *Server) RegisterHandler(method, relativePath string, handlers ...gin.HandlerFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.getEngine().Handle(method, relativePath, handlers...)
}

// RegisterHandlerAny registers handlers for all HTTP methods.
//
// Parameters:
//   - relativePath: URL path pattern
//   - handlers: Gin handler functions
//
// Example:
//
//	server.RegisterHandlerAny("/ping", pingHandler)
func (s *Server) RegisterHandlerAny(relativePath string, handlers ...gin.HandlerFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.getEngine().Any(relativePath, handlers...)
}

// Use registers global middleware.
//
// Parameters:
//   - middleware: Middleware functions to apply globally
//
// Example:
//
//	server.Use(gin.Logger(), gin.Recovery())
func (s *Server) Use(middleware ...gin.HandlerFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.getEngine().Use(middleware...)
}

// Group creates a new router group.
//
// Parameters:
//   - relativePath: Path prefix for the group
//   - handlers: Optional middleware for the group
//
// Returns:
//   - *gin.RouterGroup: Gin router group
//
// Example:
//
//	api := server.Group("/api/v1")
//	api.GET("/users", getUsersHandler)
func (s *Server) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.getEngine().Group(relativePath, handlers...)
}

// Start starts the Gin HTTP server.
//
// Parameters:
//   - address: Server address (e.g., ":8080", "localhost:3000")
//   - listenAndServeFailureFunc: Optional callback for server errors (excluding graceful shutdown)
//
// Returns:
//   - error: Error if server is already running or fails to start
//
// The server runs in a background goroutine. Use Stop() for graceful shutdown.
//
// Example:
//
//	err := server.Start(":8080", func(err error) {
//	    log.Printf("Server error: %v", err)
//	})
func (s *Server) Start(address string, listenAndServeFailureFunc func(err error)) error {
	// Atomically transition from "not running" to "running"
	if !s.running.CompareAndSwap(false, true) {
		return errors.New("server already started")
	}

	s.mutex.Lock()
	s.server = &http.Server{
		Addr:    address,
		Handler: s.getEngine(),
	}
	s.mutex.Unlock()

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if listenAndServeFailureFunc != nil {
				listenAndServeFailureFunc(err)
			}
		}
		s.running.Store(false)
	}()

	return nil
}

// Stop gracefully shuts down the Gin server.
//
// Parameters:
//   - shutdownTimeout: Maximum time to wait for shutdown
//
// Returns:
//   - error: Error if shutdown fails
//
// Example:
//
//	err := server.Stop(10 * time.Second)
func (s *Server) Stop(shutdownTimeout time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.server == nil || !s.running.Load() {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	s.running.Store(false)
	err := s.server.Shutdown(ctx)
	s.server = nil
	return err
}

// IsRunning returns whether the server is currently running.
//
// Returns:
//   - bool: true if server is running, false otherwise
//
// Example:
//
//	if server.IsRunning() {
//	    fmt.Println("Server is active")
//	}
func (s *Server) IsRunning() bool {
	return s.running.Load()
}

// GetEngine returns the underlying Gin engine for advanced usage.
//
// Returns:
//   - *gin.Engine: Gin engine instance
//
// Use this method when you need direct access to Gin's API for advanced configuration.
//
// Example:
//
//	engine := server.GetEngine()
//	engine.SetTrustedProxies([]string{"127.0.0.1"})
func (s *Server) GetEngine() *gin.Engine {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.getEngine()
}

func (s *Server) getEngine() *gin.Engine {
	if s.engine == nil {
		s.engine = gin.Default()
	}

	return s.engine
}
