// Package echo provides a simplified wrapper around the Echo web framework.
//
// This package offers convenient methods for creating and managing Echo HTTP servers
// with support for routing, middleware, and graceful shutdown.
//
// Features:
//   - Simplified server setup with Echo v4
//   - Handler and middleware registration
//   - Group routing support
//   - Graceful shutdown with timeout
//   - Thread-safe operations
//
// Example:
//
//	server := &echo.Server{}
//	server.RegisterHandler(echo.GET, "/hello", func(c echo.Context) error {
//	    return c.String(200, "Hello, World!")
//	})
//	server.Start(":8080", nil)
package echo

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
)

// Server is a struct that provides Echo HTTP server related methods.
type Server struct {
	mutex   sync.Mutex
	running atomic.Bool

	e *echo.Echo
}

// RegisterHandler registers an HTTP handler with optional middleware.
//
// Parameters:
//   - method: HTTP method (echo.GET, echo.POST, etc.)
//   - path: URL path pattern
//   - handler: Echo handler function
//   - middleware: Optional middleware functions
//
// Example:
//
//	server.RegisterHandler(echo.GET, "/users/:id", getUserHandler)
func (s *Server) RegisterHandler(method string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.getEcho().Add(method, path, handler, middleware...)
}

// Use registers global middleware.
//
// Parameters:
//   - middleware: Middleware functions to apply globally
//
// Example:
//
//	server.Use(middleware.Logger(), middleware.Recover())
func (s *Server) Use(middleware ...echo.MiddlewareFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.getEcho().Use(middleware...)
}

// Group creates a new router group with optional middleware.
//
// Parameters:
//   - prefix: Path prefix for the group
//   - middleware: Optional middleware functions for the group
//
// Returns:
//   - *echo.Group: Echo router group
//
// Example:
//
//	api := server.Group("/api/v1")
//	api.GET("/users", getUsersHandler)
func (s *Server) Group(prefix string, middleware ...echo.MiddlewareFunc) *echo.Group {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.getEcho().Group(prefix, middleware...)
}

// Start starts the Echo HTTP server.
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

	// Initialize Echo instance if not already done, under mutex to protect s.e
	s.mutex.Lock()
	e := s.getEcho()
	s.mutex.Unlock()

	go func() {
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			if listenAndServeFailureFunc != nil {
				listenAndServeFailureFunc(err)
			}
		}
		s.running.Store(false)
	}()

	return nil
}

// Stop gracefully shuts down the Echo server.
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

	if s.e == nil || !s.running.Load() {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := s.e.Shutdown(ctx)
	s.running.Store(false)
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

// GetEcho returns the underlying Echo instance for advanced usage.
//
// Returns:
//   - *echo.Echo: Echo instance
//
// Use this method when you need direct access to Echo's API for advanced configuration.
//
// Example:
//
//	e := server.GetEcho()
//	e.HideBanner = true
//	e.Debug = true
func (s *Server) GetEcho() *echo.Echo {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.getEcho()
}

func (s *Server) getEcho() *echo.Echo {
	if s.e == nil {
		s.e = echo.New()
	}
	return s.e
}
