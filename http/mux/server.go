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

// RegisterHandler registers an HTTP handler for the specified path.
//
// Parameters:
//   - path: URL path pattern to match (e.g., "/api/users")
//   - handler: HTTP handler to handle requests
//   - methods: Optional HTTP methods to restrict (e.g., http.MethodGet, http.MethodPost)
//
// If no methods are specified, the handler accepts all HTTP methods. Multiple methods
// can be provided to restrict the handler to specific request types.
//
// Example:
//
//	// Handle all methods
//	server.RegisterHandler("/api/health", healthHandler)
//
//	// Handle only GET requests
//	server.RegisterHandler("/api/users", usersHandler, http.MethodGet)
//
//	// Handle GET and POST requests
//	server.RegisterHandler("/api/items", itemsHandler, http.MethodGet, http.MethodPost)
func (s *Server) RegisterHandler(path string, handler http.Handler, methods ...string) {
	if len(methods) == 0 {
		s.getRouter().Handle(path, handler)
	} else {
		s.getRouter().Handle(path, handler).Methods(methods...)
	}
}

// RegisterHandlerFunc registers an HTTP handler function for the specified path.
//
// Parameters:
//   - path: URL path pattern to match (e.g., "/api/users")
//   - handlerFunc: HTTP handler function to handle requests
//   - methods: Optional HTTP methods to restrict (e.g., http.MethodGet, http.MethodPost)
//
// This is a convenience method for registering handler functions directly without wrapping
// them in a Handler type. If no methods are specified, accepts all HTTP methods.
//
// Example:
//
//	// Handle all methods
//	server.RegisterHandlerFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
//	    w.Write([]byte("pong"))
//	})
//
//	// Handle only GET requests
//	server.RegisterHandlerFunc("/api/users", getUsersHandler, http.MethodGet)
//
//	// Handle GET and POST requests
//	server.RegisterHandlerFunc("/api/data", dataHandler, http.MethodGet, http.MethodPost)
func (s *Server) RegisterHandlerFunc(path string, handlerFunc http.HandlerFunc, methods ...string) {
	if len(methods) == 0 {
		s.getRouter().HandleFunc(path, handlerFunc)
	} else {
		s.getRouter().HandleFunc(path, handlerFunc).Methods(methods...)
	}
}

// RegisterPathPrefixHandler registers an HTTP handler for all paths matching the prefix.
//
// Parameters:
//   - prefix: URL path prefix to match (e.g., "/static/", "/api/v1/")
//   - handler: HTTP handler to handle requests
//   - methods: Optional HTTP methods to restrict (e.g., http.MethodGet, http.MethodPost)
//
// This method matches all paths that start with the specified prefix. Useful for serving
// static files or grouping related endpoints. If no methods are specified, accepts all
// HTTP methods.
//
// Example:
//
//	// Serve static files (all methods)
//	fileServer := http.FileServer(http.Dir("./static"))
//	server.RegisterPathPrefixHandler("/static/", fileServer)
//
//	// API v1 endpoints (GET only)
//	server.RegisterPathPrefixHandler("/api/v1/", apiV1Handler, http.MethodGet)
//
//	// Admin endpoints (GET and POST)
//	server.RegisterPathPrefixHandler("/admin/", adminHandler, http.MethodGet, http.MethodPost)
func (s *Server) RegisterPathPrefixHandler(prefix string, handler http.Handler, methods ...string) {
	if len(methods) == 0 {
		s.getRouter().PathPrefix(prefix).Handler(handler)
	} else {
		s.getRouter().PathPrefix(prefix).Handler(handler).Methods(methods...)
	}
}

// RegisterPathPrefixHandlerFunc registers an HTTP handler function for paths matching the prefix.
//
// Parameters:
//   - prefix: URL path prefix to match (e.g., "/static/", "/api/v1/")
//   - handlerFunc: HTTP handler function to handle requests
//   - methods: Optional HTTP methods to restrict (e.g., http.MethodGet, http.MethodPost)
//
// This is a convenience method for registering handler functions for path prefixes without
// wrapping them in a Handler type. Matches all paths starting with the prefix. If no methods
// are specified, accepts all HTTP methods.
//
// Example:
//
//	// Handle all /api/ requests
//	server.RegisterPathPrefixHandlerFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
//	    w.Header().Set("Content-Type", "application/json")
//	    w.Write([]byte(`{"status": "ok"}`))
//	})
//
//	// Handle /files/ with GET only
//	server.RegisterPathPrefixHandlerFunc("/files/", filesHandler, http.MethodGet)
func (s *Server) RegisterPathPrefixHandlerFunc(prefix string, handlerFunc http.HandlerFunc, methods ...string) {
	if len(methods) == 0 {
		s.getRouter().PathPrefix(prefix).HandlerFunc(handlerFunc)
	} else {
		s.getRouter().PathPrefix(prefix).HandlerFunc(handlerFunc).Methods(methods...)
	}
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
