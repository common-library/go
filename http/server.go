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
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server is a struct that provides server related methods.
type Server struct {
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

// Start starts the HTTP server on the specified address.
//
// Parameters:
//   - address: Address to bind the server to (e.g., ":8080", "localhost:8080")
//   - listenAndServeFailureFunc: Optional callback function called if server fails to start
//   - middlewareFunc: Optional middleware functions to apply to all requests
//
// Returns:
//   - error: Always returns nil (errors are reported via listenAndServeFailureFunc)
//
// The server starts in a background goroutine, so this method returns immediately. If no
// middleware is provided, a pass-through middleware is added. Multiple middleware functions
// are executed in the order provided.
//
// Example:
//
//	var server http.Server
//	server.RegisterHandlerFunc("/api", apiHandler)
//
//	// Start without middleware
//	err := server.Start(":8080", func(err error) {
//	    log.Fatal(err)
//	})
//
//	// Start with logging middleware
//	loggingMiddleware := func(next http.Handler) http.Handler {
//	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	        log.Printf("%s %s", r.Method, r.URL.Path)
//	        next.ServeHTTP(w, r)
//	    })
//	}
//	err = server.Start(":8080", nil, loggingMiddleware)
func (s *Server) Start(address string, listenAndServeFailureFunc func(err error), middlewareFunc ...mux.MiddlewareFunc) error {
	if middlewareFunc != nil {
		s.getRouter().Use(middlewareFunc...)
	} else {
		s.getRouter().Use(func(nextHandler http.Handler) http.Handler {
			return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
				nextHandler.ServeHTTP(responseWriter, request)
			})
		})
	}

	s.server = &http.Server{
		Addr:    address,
		Handler: s.getRouter()}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed && listenAndServeFailureFunc != nil {
			listenAndServeFailureFunc(err)
		}
	}()

	return nil
}

// Stop gracefully shuts down the HTTP server.
//
// Parameters:
//   - shutdownTimeout: Maximum duration in seconds to wait for active connections to complete
//
// Returns:
//   - error: Error if shutdown fails or times out, nil on successful shutdown
//
// The server stops accepting new connections immediately and waits for active connections
// to complete within the timeout period. After the timeout, remaining connections are
// forcefully closed. The router is reset to nil after shutdown.
//
// Example:
//
//	var server http.Server
//	server.Start(":8080", nil)
//
//	// Later, graceful shutdown with 30 second timeout
//	err := server.Stop(30)
//	if err != nil {
//	    log.Printf("Shutdown error: %v", err)
//	}
//	log.Println("Server stopped gracefully")
func (s *Server) Stop(shutdownTimeout time.Duration) error {
	s.SetRouter(nil)

	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// SetRouter sets a custom router for the server.
//
// Parameters:
//   - router: gorilla/mux router instance, or nil to reset
//
// This method allows using a pre-configured router instead of the default one. Setting
// the router to nil will reset it, and a new router will be created on the next operation.
// This is useful for advanced routing configurations or testing.
//
// Example:
//
//	// Use custom router
//	router := mux.NewRouter()
//	router.StrictSlash(true)
//	router.HandleFunc("/", homeHandler)
//
//	var server http.Server
//	server.SetRouter(router)
//	server.Start(":8080", nil)
//
//	// Reset router
//	server.SetRouter(nil)
func (s *Server) SetRouter(router *mux.Router) {
	s.router = router
}

func (s *Server) getRouter() *mux.Router {
	if s.router == nil {
		s.router = mux.NewRouter()
	}

	return s.router
}
