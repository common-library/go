// Package long_polling provides HTTP long polling server and client implementations.
//
// This package enables real-time communication using HTTP long polling patterns,
// allowing clients to receive server-side events with minimal latency without
// WebSocket connections. It wraps golongpoll for easy server setup and client communication.
//
// # Features
//
//   - Long polling server with configurable timeouts
//   - Event subscription and publishing
//   - File-based persistence for event durability
//   - Custom middleware support for authentication and validation
//   - Client helpers for subscription and publishing
//
// # Basic Server Example
//
//	server := &long_polling.Server{}
//	err := server.Start(long_polling.ServerInfo{
//	    Address: ":8080",
//	    TimeoutSeconds: 120,
//	    SubscriptionURI: "/events",
//	    PublishURI: "/publish",
//	}, long_polling.FilePersistorInfo{Use: false}, nil)
package long_polling

import (
	"net/http"
	"time"

	"github.com/common-library/go/http/mux"
	"github.com/jcuga/golongpoll"
)

// ServerInfo is server information.
type ServerInfo struct {
	Address        string
	TimeoutSeconds int

	SubscriptionURI                string
	HandlerToRunBeforeSubscription func(w http.ResponseWriter, r *http.Request) bool

	PublishURI                string
	HandlerToRunBeforePublish func(w http.ResponseWriter, r *http.Request) bool
}

// FilePersistorInfo is file persistor information.
type FilePersistorInfo struct {
	Use                     bool
	FileName                string
	WriteBufferSize         int
	WriteFlushPeriodSeconds int
}

// Server is a struct that provides server related methods.
type Server struct {
	server mux.Server

	longpollManager *golongpoll.LongpollManager
}

// Start initializes and starts the long polling server.
//
// This method creates a long polling manager, sets up subscription and publish
// handlers, and starts the HTTP server. It supports optional file persistence
// for event durability across server restarts.
//
// # Parameters
//
//   - serverInfo: Server configuration including address, timeout, URIs, and middleware
//   - filePersistorInfo: File persistence configuration for event durability
//   - listenAndServeFailureFunc: Optional callback for listen and serve failures
//
// # Returns
//
//   - error: Error if server initialization or start fails, nil on success
//
// # Behavior
//
// The server creates two endpoints:
//   - Subscription endpoint (GET): Clients subscribe for events on specific categories
//   - Publish endpoint (POST): Server or authorized clients publish events
//
// Custom handlers can be provided to run before subscription or publish operations,
// enabling authentication, validation, or logging. If a handler returns false,
// the request is rejected.
//
// # Examples
//
// Basic server:
//
//	server := &long_polling.Server{}
//	err := server.Start(long_polling.ServerInfo{
//	    Address: ":8080",
//	    TimeoutSeconds: 120,
//	    SubscriptionURI: "/events",
//	    PublishURI: "/publish",
//	}, long_polling.FilePersistorInfo{Use: false}, nil)
func (s *Server) Start(serverInfo ServerInfo, filePersistorInfo FilePersistorInfo, listenAndServeFailureFunc func(err error)) error {
	option := golongpoll.Options{
		LoggingEnabled:            false,
		MaxLongpollTimeoutSeconds: serverInfo.TimeoutSeconds,
		//MaxEventBufferSize: 250,
		//EventTimeToLiveSeconds:,
		//DeleteEventAfterFirstRetrieval:,
	}

	if filePersistorInfo.Use {
		filePersistor, err := golongpoll.NewFilePersistor(filePersistorInfo.FileName, filePersistorInfo.WriteBufferSize, filePersistorInfo.WriteFlushPeriodSeconds)
		if err != nil {
			return err
		}

		option.AddOn = filePersistor
	}

	longpollManager, err := golongpoll.StartLongpoll(option)
	if err != nil {
		return err
	}
	s.longpollManager = longpollManager

	// Get router and configure handlers
	router := s.server.GetRouter()

	router.Use(func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			nextHandler.ServeHTTP(responseWriter, request)
		})
	})

	subscriptionHandler := func() func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			if serverInfo.HandlerToRunBeforeSubscription != nil &&
				!serverInfo.HandlerToRunBeforeSubscription(w, r) {
				return
			}

			s.longpollManager.SubscriptionHandler(w, r)
		}
	}
	router.HandleFunc(serverInfo.SubscriptionURI, subscriptionHandler()).Methods(http.MethodGet)

	publishHandler := func() func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			if serverInfo.HandlerToRunBeforePublish != nil &&
				!serverInfo.HandlerToRunBeforePublish(w, r) {
				return
			}

			s.longpollManager.PublishHandler(w, r)
		}
	}
	router.HandleFunc(serverInfo.PublishURI, publishHandler()).Methods(http.MethodPost)

	return s.server.Start(serverInfo.Address, listenAndServeFailureFunc)
}

// Stop gracefully shuts down the long polling server.
//
// This method stops the HTTP server with a timeout for existing connections
// and shuts down the long polling manager, ensuring all events are flushed
// if file persistence is enabled.
//
// # Parameters
//
//   - shutdownTimeout: Maximum duration to wait for active connections to close
//
// # Returns
//
//   - error: Error if shutdown fails, nil on successful shutdown
//
// # Behavior
//
// The shutdown process:
//  1. Stops accepting new connections
//  2. Waits up to shutdownTimeout for active connections to close
//  3. Shuts down the long polling manager (flushes pending events)
//  4. Closes all resources
//
// If the timeout is reached before all connections close, the server
// forcibly terminates remaining connections.
//
// # Examples
//
// Graceful shutdown with 10 second timeout:
//
//	err := server.Stop(10 * time.Second)
//	if err != nil {
//	    log.Printf("Shutdown error: %v", err)
//	}
func (s *Server) Stop(shutdownTimeout time.Duration) error {
	if s.longpollManager != nil {
		defer s.longpollManager.Shutdown()
	}

	err := s.server.Stop(shutdownTimeout)
	if err != nil {
		return err
	}

	return nil
}
