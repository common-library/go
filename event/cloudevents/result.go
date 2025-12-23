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
	"errors"

	cloudeventssdk "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	cloudeventssdk_http "github.com/cloudevents/sdk-go/v2/protocol/http"
)

// NewResult creates and returns a generic CloudEvents result.
//
// Parameters:
//   - format: Printf-style format string describing the result
//   - arguments: Optional format arguments
//
// Returns:
//   - Result: CloudEvents result with the formatted message
//
// Use this for protocol-agnostic results. For HTTP-specific results with status codes,
// use NewHTTPResult instead.
//
// Example:
//
//	// Success result
//	result := cloudevents.NewResult("event processed successfully")
//
//	// Error result with details
//	result = cloudevents.NewResult("validation failed: %s", validationError)
func NewResult(format string, arguments ...any) Result {
	return Result{result: cloudeventssdk.NewResult(format, arguments)}
}

// NewHTTPResult creates and returns an HTTP-specific CloudEvents result with status code.
//
// Parameters:
//   - statusCode: HTTP status code (e.g., 200, 400, 500)
//   - format: Printf-style format string describing the result
//   - arguments: Optional format arguments
//
// Returns:
//   - Result: CloudEvents result with HTTP status code and formatted message
//
// This result type includes HTTP-specific information and can be retrieved using
// GetHttpStatusCode method.
//
// Example:
//
//	// Success with 200 OK
//	result := cloudevents.NewHTTPResult(200, "event accepted")
//
//	// Client error with 400 Bad Request
//	result = cloudevents.NewHTTPResult(400, "invalid event type: %s", eventType)
//
//	// Server error with 500 Internal Server Error
//	result = cloudevents.NewHTTPResult(500, "processing error: %v", err)
func NewHTTPResult(statusCode int, format string, arguments ...any) Result {
	return Result{result: cloudeventssdk.NewHTTPResult(statusCode, format, arguments)}
}

// Result is the result of event delivery.
type Result struct {
	result protocol.Result
}

// IsACK returns whether the recipient acknowledged the event.
//
// Returns:
//   - bool: true if event was acknowledged (ACK), false otherwise
//
// An ACK indicates successful delivery and processing. For HTTP, this typically
// corresponds to 2xx status codes.
//
// Example:
//
//	result := client.Send(event)
//	if result.IsACK() {
//	    log.Println("Event acknowledged")
//	}
func (r *Result) IsACK() bool {
	return cloudeventssdk.IsACK(r.result)
}

// IsNACK returns whether the recipient did not acknowledge the event.
//
// Returns:
//   - bool: true if event was not acknowledged (NACK), false otherwise
//
// A NACK indicates the event was delivered but rejected or failed processing.
// For HTTP, this typically corresponds to 4xx or 5xx status codes.
//
// Example:
//
//	result := client.Send(event)
//	if result.IsNACK() {
//	    log.Printf("Event rejected: %s", result.Error())
//	}
func (r *Result) IsNACK() bool {
	return cloudeventssdk.IsNACK(r.result)
}

// IsUndelivered returns whether the event could not be delivered.
//
// Returns:
//   - bool: true if event was undelivered, false if delivery was attempted
//
// Undelivered indicates a network error, connection failure, or other issue that
// prevented the event from reaching the recipient. This is different from NACK
// where the event was delivered but rejected.
//
// Example:
//
//	result := client.Send(event)
//	if result.IsUndelivered() {
//	    log.Printf("Delivery failed: %s", result.Error())
//	    // Retry or handle connection error
//	}
func (r *Result) IsUndelivered() bool {
	return cloudeventssdk.IsUndelivered(r.result)
}

// GetHttpStatusCode extracts the HTTP status code from an HTTP result.
//
// Returns:
//   - int: HTTP status code (e.g., 200, 404, 500), or -1 if not an HTTP result
//   - error: Error if result is not an HTTP result, nil otherwise
//
// This method only works with results created by NewHTTPResult or returned from
// HTTP-based event operations. Returns an error for non-HTTP results.
//
// Example:
//
//	result := client.Send(event)
//	statusCode, err := result.GetHttpStatusCode()
//	if err != nil {
//	    log.Println("Not an HTTP result")
//	} else {
//	    log.Printf("HTTP Status: %d", statusCode)
//	    if statusCode >= 400 {
//	        log.Println("HTTP error occurred")
//	    }
//	}
func (r *Result) GetHttpStatusCode() (int, error) {
	httpResult := new(cloudeventssdk_http.Result)

	if !cloudeventssdk.ResultAs(r.result, &httpResult) {
		return -1, errors.New("match failed")
	} else {
		return httpResult.StatusCode, nil
	}
}

// Error returns the error message string from the result.
//
// Returns:
//   - string: Error message, or empty string if no error
//
// This implements the error interface, allowing Result to be used as an error type.
// The message includes details about why an event was NACK'd or undelivered.
//
// Example:
//
//	result := client.Send(event)
//	if !result.IsACK() {
//	    log.Printf("Event failed: %s", result.Error())
//	}
func (r *Result) Error() string {
	return r.result.Error()
}
