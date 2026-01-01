// Package http provides simplified HTTP client and server utilities.
//
// This package wraps Go's standard net/http and gorilla/mux packages, offering convenient
// functions for making HTTP requests and managing HTTP servers with routing capabilities.
//
// Features:
//   - Simplified HTTP client with authentication support
//   - HTTP server with routing (powered by gorilla/mux)
//   - Handler and middleware registration
//   - Path prefix routing
//   - Graceful server shutdown
//   - Custom transport configuration
//
// Example:
//
//	// Client
//	resp, _ := http.Request("http://api.example.com", http.MethodGet, nil, "", 10, "", "", nil)
//
//	// Server
//	var server http.Server
//	server.RegisterHandlerFunc("/api", handler)
//	server.Start(":8080", nil)
package mux

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

// Server is a struct that provides server related methods.
type Server struct {
	mutex   sync.Mutex
	running atomic.Bool

	server *http.Server
	router *mux.Router
}

// RegisterHandler registers an HTTP handler for the specified path and method.
//
// Parameters:
//   - method: HTTP method (http.MethodGet, http.MethodPost, etc.)
//   - path: URL path pattern to match (e.g., "/api/users")
//   - handler: HTTP handler to handle requests
//
// Example:
//
//	// Handle GET requests
//	server.RegisterHandler(http.MethodGet, "/api/users", usersHandler)
//
//	// Handle POST requests
//	server.RegisterHandler(http.MethodPost, "/api/items", itemsHandler)
func (s *Server) RegisterHandler(method, path string, handler http.Handler) {
	s.getRouter().Handle(path, handler).Methods(method)
}

// RegisterHandlerAny registers an HTTP handler for all methods on the specified path.
//
// Parameters:
//   - path: URL path pattern to match (e.g., "/api/webhook")
//   - handler: HTTP handler to handle requests
//
// Example:
//
//	// Handle all HTTP methods
//	server.RegisterHandlerAny("/api/webhook", webhookHandler)
func (s *Server) RegisterHandlerAny(path string, handler http.Handler) {
	s.getRouter().Handle(path, handler)
}

// RegisterHandlerFunc registers an HTTP handler function for the specified path and method.
//
// Parameters:
//   - method: HTTP method (http.MethodGet, http.MethodPost, etc.)
//   - path: URL path pattern to match (e.g., "/api/users")
//   - handlerFunc: HTTP handler function to handle requests
//
// This is a convenience method for registering handler functions directly without wrapping
// them in a Handler type.
//
// Example:
//
//	// Handle GET requests
//	server.RegisterHandlerFunc(http.MethodGet, "/api/ping", func(w http.ResponseWriter, r *http.Request) {
//	    w.Write([]byte("pong"))
//	})
//
//	// Handle POST requests
//	server.RegisterHandlerFunc(http.MethodPost, "/api/data", dataHandler)
func (s *Server) RegisterHandlerFunc(method, path string, handlerFunc http.HandlerFunc) {
	s.getRouter().HandleFunc(path, handlerFunc).Methods(method)
}

// RegisterHandlerFuncAny registers an HTTP handler function for all methods on the specified path.
//
// Parameters:
//   - path: URL path pattern to match (e.g., "/api/catch-all")
//   - handlerFunc: HTTP handler function to handle requests
//
// Example:
//
//	// Handle all HTTP methods
//	server.RegisterHandlerFuncAny("/api/catch-all", func(w http.ResponseWriter, r *http.Request) {
//	    w.Write([]byte("method: " + r.Method))
//	})
func (s *Server) RegisterHandlerFuncAny(path string, handlerFunc http.HandlerFunc) {
	s.getRouter().HandleFunc(path, handlerFunc)
}

// RegisterPathPrefixHandler registers an HTTP handler for all paths matching the prefix and method.
//
// Parameters:
//   - method: HTTP method (http.MethodGet, http.MethodPost, etc.)
//   - prefix: URL path prefix to match (e.g., "/static/", "/api/v1/")
//   - handler: HTTP handler to handle requests
//
// This method matches all paths that start with the specified prefix. Useful for serving
// static files or grouping related endpoints.
//
// Example:
//
//	// Serve static files with GET
//	fileServer := http.FileServer(http.Dir("./static"))
//	server.RegisterPathPrefixHandler(http.MethodGet, "/static/", fileServer)
//
//	// API v1 endpoints
//	server.RegisterPathPrefixHandler(http.MethodGet, "/api/v1/", apiV1Handler)
func (s *Server) RegisterPathPrefixHandler(method, prefix string, handler http.Handler) {
	s.getRouter().PathPrefix(prefix).Handler(handler).Methods(method)
}

// RegisterPathPrefixHandlerAny registers an HTTP handler for all methods matching the prefix.
//
// Parameters:
//   - prefix: URL path prefix to match (e.g., "/static/", "/api/")
//   - handler: HTTP handler to handle requests
//
// Example:
//
//	// Serve static files (all methods)
//	fileServer := http.FileServer(http.Dir("./static"))
//	server.RegisterPathPrefixHandlerAny("/static/", fileServer)
func (s *Server) RegisterPathPrefixHandlerAny(prefix string, handler http.Handler) {
	s.getRouter().PathPrefix(prefix).Handler(handler)
}

// RegisterPathPrefixHandlerFunc registers an HTTP handler function for paths matching the prefix and method.
//
// Parameters:
//   - method: HTTP method (http.MethodGet, http.MethodPost, etc.)
//   - prefix: URL path prefix to match (e.g., "/static/", "/api/v1/")
//   - handlerFunc: HTTP handler function to handle requests
//
// This is a convenience method for registering handler functions for path prefixes without
// wrapping them in a Handler type. Matches all paths starting with the prefix.
//
// Example:
//
//	// Handle all /api/ GET requests
//	server.RegisterPathPrefixHandlerFunc(http.MethodGet, "/api/", func(w http.ResponseWriter, r *http.Request) {
//	    w.Header().Set("Content-Type", "application/json")
//	    w.Write([]byte(`{"status": "ok"}`))
//	})
//
//	// Handle /files/ with GET only
//	server.RegisterPathPrefixHandlerFunc(http.MethodGet, "/files/", filesHandler)
func (s *Server) RegisterPathPrefixHandlerFunc(method, prefix string, handlerFunc http.HandlerFunc) {
	s.getRouter().PathPrefix(prefix).HandlerFunc(handlerFunc).Methods(method)
}

// RegisterPathPrefixHandlerFuncAny registers an HTTP handler function for all methods matching the prefix.
//
// Parameters:
//   - prefix: URL path prefix to match (e.g., "/api/", "/files/")
//   - handlerFunc: HTTP handler function to handle requests
//
// Example:
//
//	// Handle all /api/ requests (all methods)
//	server.RegisterPathPrefixHandlerFuncAny("/api/", func(w http.ResponseWriter, r *http.Request) {
//	    w.Header().Set("Content-Type", "application/json")
//	    w.Write([]byte(`{"status": "ok"}`))
//	})
func (s *Server) RegisterPathPrefixHandlerFuncAny(prefix string, handlerFunc http.HandlerFunc) {
	s.getRouter().PathPrefix(prefix).HandlerFunc(handlerFunc)
}

// Use registers global middleware.
//
// Parameters:
//   - middleware: Middleware functions to apply globally
//
// Example:
//
//	server.Use(loggingMiddleware)
//	server.Use(authMiddleware)
func (s *Server) Use(middleware ...mux.MiddlewareFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.getRouter().Use(middleware...)
}

// Start starts the HTTP server on the specified address.
//
// Parameters:
//   - address: Address to bind the server to (e.g., ":8080", "localhost:8080")
//   - listenAndServeFailureFunc: Optional callback function called if server fails to start
//
// Returns:
//   - error: Error if server is already started, nil otherwise
//
// The server starts in a background goroutine, so this method returns immediately.
// Use the Use() method to register middleware before calling Start().
//
// Example:
//
//	server := &mux.Server{}
//	server.Use(loggingMiddleware)
//	server.RegisterHandlerFunc("/api", apiHandler)
//	err := server.Start(":8080", func(err error) {
//	    log.Fatal(err)
//	})
func (s *Server) Start(address string, listenAndServeFailureFunc func(err error)) error {
	if !s.running.CompareAndSwap(false, true) {
		return errors.New("server already started")
	}

	s.mutex.Lock()
	s.server = &http.Server{
		Addr:    address,
		Handler: s.getRouter(),
	}
	server := s.server
	s.mutex.Unlock()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.running.Store(false)
			if listenAndServeFailureFunc != nil {
				listenAndServeFailureFunc(err)
			}
		}
	}()

	return nil
}

// Stop gracefully shuts down the HTTP server.
//
// Parameters:
//   - shutdownTimeout: Maximum duration to wait for active connections to complete
//
// Returns:
//   - error: Error if shutdown fails or times out, nil on successful shutdown
//
// The server stops accepting new connections immediately and waits for active connections
// to complete within the timeout period. After the timeout, remaining connections are
// forcefully closed.
//
// Example:
//
//	server := &mux.Server{}
//	server.Start(":8080", nil)
//
//	// Later, graceful shutdown with 10 second timeout
//	err := server.Stop(10 * time.Second)
//	if err != nil {
//	    log.Printf("Shutdown error: %v", err)
//	}
func (s *Server) Stop(shutdownTimeout time.Duration) error {
	s.mutex.Lock()
	server := s.server
	s.server = nil
	s.mutex.Unlock()

	if server == nil {
		return nil
	}

	defer s.running.Store(false)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return server.Shutdown(ctx)
}

// IsRunning returns whether the server is currently running.
//
// Returns:
//   - bool: true if server is running, false otherwise
//
// Example:
//
//	if server.IsRunning() {
//	    log.Println("Server is running")
//	}
func (s *Server) IsRunning() bool {
	return s.running.Load()
}

// GetRouter returns the gorilla/mux router instance.
//
// Returns:
//   - *mux.Router: The router instance
//
// This method provides direct access to the underlying router for advanced configurations.
//
// Example:
//
//	server := &mux.Server{}
//	router := server.GetRouter()
//	router.StrictSlash(true)
//	router.HandleFunc("/", homeHandler)
func (s *Server) GetRouter() *mux.Router {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.getRouter()
}

func (s *Server) getRouter() *mux.Router {
	if s.router == nil {
		s.router = mux.NewRouter()
	}

	return s.router
}
