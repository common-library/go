// Package cloudevents provides CloudEvents client and server implementations.
//
// This package wraps the official CloudEvents SDK for Go, offering simplified interfaces
// for sending, receiving, and processing CloudEvents over HTTP and other protocols.
//
// Features:
//   - HTTP client for sending and receiving CloudEvents
//   - Request-response pattern support
//   - Asynchronous event receiver with lifecycle management
//   - Result types for event delivery status
//   - CloudEvents v1.0 specification compliance
//
// Example:
//
//	client, _ := cloudevents.NewHttp("http://localhost:8080", nil, nil)
//	event := cloudevents.NewEvent()
//	event.SetType("com.example.event")
//	event.SetSource("example/source")
//	result := client.Send(event)
package cloudevents

import (
	"context"
	"time"

	cloudeventssdk "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/common-library/go/http"
)

// Server is a struct that provides server related methods.
type Server struct {
	server http.Server
}

// Start starts the CloudEvents HTTP server on the specified address.
//
// Parameters:
//   - address: Server address to bind to (e.g., "localhost:8080" or ":8080")
//   - handler: Function to process incoming events and optionally return response events
//   - listenAndServeFailureFunc: Function called if the HTTP server encounters a fatal error
//
// Returns:
//   - error: Error if server initialization fails, nil if started successfully
//
// The handler function receives each incoming CloudEvent and can return a response event
// for request-response patterns. Return nil for the event if no response is needed.
// The server runs asynchronously; this method returns after starting the server.
//
// Example:
//
//	var server cloudevents.Server
//
//	handler := func(event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
//	    log.Printf("Received: %s", event.Type())
//
//	    // Process event...
//
//	    // Return response event
//	    response := cloudevents.NewEvent()
//	    response.SetType("com.example.response")
//	    response.SetSource("example/server")
//	    return &response, cloudevents.NewHTTPResult(200, "OK")
//	}
//
//	failureFunc := func(err error) {
//	    log.Printf("Server error: %v", err)
//	}
//
//	err := server.Start(":8080", handler, failureFunc)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (s *Server) Start(address string, handler func(Event) (*Event, Result), listenAndServeFailureFunc func(error)) error {
	finalHandler := func(requestEvent Event) (*Event, protocol.Result) {
		responseEvent, result := handler(requestEvent)
		return responseEvent, result.result
	}

	if protocol, err := cloudeventssdk.NewHTTP(); err != nil {
		return err
	} else if eventReceiver, err := cloudeventssdk.NewHTTPReceiveHandler(context.Background(), protocol, finalHandler); err != nil {

		return err
	} else {
		s.server.RegisterPathPrefixHandler("/", eventReceiver)

		return s.server.Start(address, listenAndServeFailureFunc)
	}
}

// Stop gracefully shuts down the CloudEvents server.
//
// Parameters:
//   - shutdownTimeout: Maximum duration to wait for active connections to complete
//
// Returns:
//   - error: Error if shutdown fails or times out, nil if successful
//
// The server stops accepting new connections immediately and waits for active
// connections to complete within the timeout period. After the timeout, the server
// forcefully closes remaining connections.
//
// Example:
//
//	server.Start(":8080", handler, failureFunc)
//	// ... server running ...
//
//	// Graceful shutdown with 10 second timeout
//	err := server.Stop(10 * time.Second)
//	if err != nil {
//	    log.Printf("Shutdown error: %v", err)
//	}
//	log.Println("Server stopped")
func (s *Server) Stop(shutdownTimeout time.Duration) error {
	return s.server.Stop(shutdownTimeout)
}
