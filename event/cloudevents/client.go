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
	"sync"

	cloudeventssdk "github.com/cloudevents/sdk-go/v2"
	cloudeventssdk_client "github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
)

type clientType int

const (
	clientTypeHttp clientType = iota + 1
)

// NewHttp creates and returns an HTTP client for CloudEvents.
//
// Parameters:
//   - address: Target URL for sending events (e.g., "http://localhost:8080")
//   - httpOption: HTTP protocol options from CloudEvents SDK (e.g., WithPort, WithPath)
//   - clientOption: Client configuration options from CloudEvents SDK
//
// Returns:
//   - *client: Configured CloudEvents client
//   - error: Error if client creation or protocol initialization fails
//
// The client can be used for sending events, making request-response calls, or receiving
// events. For receivers, leave address empty and configure port via httpOption.
//
// Example:
//
//	// Sender client
//	client, err := cloudevents.NewHttp("http://localhost:8080", nil, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Receiver client
//	httpOpts := []http.Option{http.WithPort(8080)}
//	receiver, err := cloudevents.NewHttp("", httpOpts, nil)
func NewHttp(address string, httpOption []http.Option, clientOption []cloudeventssdk_client.Option) (*client, error) {
	if protocol, err := cloudeventssdk.NewHTTP(httpOption...); err != nil {
		return nil, err
	} else if clientOfSdk, err := cloudeventssdk.NewClient(protocol, clientOption...); err != nil {
		return nil, err
	} else {
		return &client{clientType: clientTypeHttp, clientOfSdk: clientOfSdk, address: address}, nil
	}
}

// client is a struct that provides client related methods.
type client struct {
	clientType clientType

	clientOfSdk cloudeventssdk_client.Client
	address     string

	wgForReceiver         sync.WaitGroup
	cancelFuncForReceiver context.CancelFunc
}

// Send transmits a CloudEvent to the configured target address.
//
// Parameters:
//   - event: CloudEvent to send
//
// Returns:
//   - Result: Delivery result containing status (ACK/NACK/Undelivered) and HTTP status code
//
// This is a one-way fire-and-forget operation. The result indicates whether the event
// was successfully delivered, but no response event is returned.
//
// Example:
//
//	event := cloudevents.NewEvent()
//	event.SetType("com.example.user.created")
//	event.SetSource("example/users")
//	event.SetData("application/json", map[string]string{"name": "Alice"})
//
//	result := client.Send(event)
//	if result.IsUndelivered() {
//	    log.Printf("Failed to send: %s", result.Error())
//	} else {
//	    log.Println("Event sent successfully")
//	}
func (c *client) Send(event Event) Result {
	return Result{result: c.clientOfSdk.Send(c.getContext(), event)}
}

// Request transmits a CloudEvent and waits for a response event.
//
// Parameters:
//   - event: CloudEvent to send
//
// Returns:
//   - *Event: Response CloudEvent from the receiver, nil if request failed
//   - Result: Delivery result containing status and HTTP status code
//
// This implements the request-response pattern where the client sends an event and
// expects a response event back. The call blocks until a response is received or
// an error occurs.
//
// Example:
//
//	event := cloudevents.NewEvent()
//	event.SetType("com.example.query")
//	event.SetSource("example/client")
//	event.SetData("application/json", map[string]string{"query": "status"})
//
//	response, result := client.Request(event)
//	if result.IsUndelivered() {
//	    log.Printf("Request failed: %s", result.Error())
//	} else {
//	    log.Printf("Response: %v", response.Data())
//	}
func (c *client) Request(event Event) (*Event, Result) {
	responseEvent, result := c.clientOfSdk.Request(c.getContext(), event)

	return responseEvent, Result{result: result}
}

// StartReceiver starts receiving CloudEvents in a background goroutine.
//
// Parameters:
//   - handler: Function called for each received event with context and event data
//   - failureFunc: Function called if the receiver encounters a fatal error
//
// The receiver runs asynchronously until StopReceiver is called. Events are processed
// by the handler function in the order they are received. This method returns immediately
// after starting the receiver goroutine.
//
// Example:
//
//	httpOpts := []http.Option{http.WithPort(8080)}
//	receiver, err := cloudevents.NewHttp("", httpOpts, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	handler := func(ctx context.Context, event cloudevents.Event) {
//	    log.Printf("Received event: %s", event.Type())
//	    // Process event...
//	}
//
//	failureFunc := func(err error) {
//	    log.Printf("Receiver error: %v", err)
//	}
//
//	receiver.StartReceiver(handler, failureFunc)
//	defer receiver.StopReceiver()
func (c *client) StartReceiver(handler func(context.Context, Event), failureFunc func(error)) {
	c.wgForReceiver.Add(1)
	go func() {
		defer c.wgForReceiver.Done()

		ctx, cancel := context.WithCancel(c.getContext())
		c.cancelFuncForReceiver = cancel

		if err := c.clientOfSdk.StartReceiver(ctx, handler); err != nil {
			failureFunc(err)
		}
	}()
}

// StopReceiver gracefully stops the event receiver started by StartReceiver.
//
// This method blocks until the receiver goroutine has fully terminated. It is safe
// to call StopReceiver even if StartReceiver was never called or has already stopped.
//
// Example:
//
//	receiver.StartReceiver(handler, failureFunc)
//	// ... receive events ...
//	receiver.StopReceiver() // Blocks until receiver stops
//	log.Println("Receiver stopped")
func (c *client) StopReceiver() {
	if c.cancelFuncForReceiver != nil {
		c.cancelFuncForReceiver()
	}
	c.wgForReceiver.Wait()
}

func (c *client) getContext() context.Context {
	switch c.clientType {
	case clientTypeHttp:
		return cloudeventssdk.ContextWithTarget(context.Background(), c.address)
	default:
		return nil
	}
}
