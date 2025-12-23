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
// # Basic Client Example
//
//	response, err := long_polling.Subscription(
//	    "http://localhost:8080/events",
//	    nil,
//	    long_polling.SubscriptionRequest{
//	        Category: "notifications",
//	        TimeoutSeconds: 60,
//	    },
//	    "", "", nil,
//	)
package long_polling

import (
	"fmt"
	net_http "net/http"
	net_url "net/url"
	"time"

	"github.com/common-library/go/http"
	"github.com/common-library/go/json"
	"github.com/google/go-querystring/query"
)

// SubscriptionRequest is subscription request information.
type SubscriptionRequest struct {
	Category       string `url:"category"`
	TimeoutSeconds int    `url:"timeout"`
	SinceTime      int64  `url:"since_time,omitempty"`
	LastID         string `url:"last_id,omitempty"`
}

// SubscriptionResponse is subscription response information.
type SubscriptionResponse struct {
	Header     net_http.Header
	StatusCode int
	Events     []struct {
		Timestamp int64  `json:"timestamp"`
		Category  string `json:"category"`
		ID        string `json:"id"`
		Data      string `json:"data"`
	} `json:"events"`
}

// PublishRequest is publish request information.
type PublishRequest struct {
	Category string `json:"category"`
	Data     string `json:"data"`
}

// Subscription subscribes to server events using HTTP long polling.
//
// This function sends a long polling request to the server and waits for
// events in the specified category. The request blocks until an event occurs
// or the timeout expires.
//
// # Parameters
//
//   - url: Server subscription endpoint URL (e.g., "http://localhost:8080/events")
//   - header: Optional HTTP headers (e.g., custom headers or authorization)
//   - request: Subscription parameters (category, timeout, since_time, last_id)
//   - username: Optional HTTP basic authentication username
//   - password: Optional HTTP basic authentication password
//   - transport: Optional custom HTTP transport for connection pooling or proxies
//
// # Returns
//
//   - SubscriptionResponse: Response containing events and metadata
//   - error: Error if request fails or response parsing fails, nil on success
//
// # Behavior
//
// The subscription request includes:
//   - category: Event category to subscribe to
//   - timeout: Maximum seconds to wait for events
//   - since_time: Optional timestamp to retrieve events since (Unix milliseconds)
//   - last_id: Optional last event ID to retrieve events after
//
// The response includes:
//   - Header: HTTP response headers
//   - StatusCode: HTTP status code (200 for events, 204 for timeout)
//   - Events: Array of events (empty if timeout)
//
// # Examples
//
// Basic subscription:
//
//	response, err := long_polling.Subscription(
//	    "http://localhost:8080/events",
//	    nil,
//	    long_polling.SubscriptionRequest{
//	        Category: "notifications",
//	        TimeoutSeconds: 60,
//	    },
//	    "", "", nil,
//	)
func Subscription(url string, header map[string][]string, request SubscriptionRequest, username, password string, transport *net_http.Transport) (SubscriptionResponse, error) {
	u, err := net_url.Parse(url)
	if err != nil {
		return SubscriptionResponse{}, err
	}

	values, err := query.Values(request)
	if err != nil {
		return SubscriptionResponse{}, err
	}
	u.RawQuery = values.Encode()

	response, err := http.Request(fmt.Sprintf("%v", u), net_http.MethodGet, header, "", time.Duration(request.TimeoutSeconds)*time.Second, username, password, transport)
	if err != nil {
		return SubscriptionResponse{}, err
	}

	subscriptionResponse := SubscriptionResponse{Header: response.Header, StatusCode: response.StatusCode}

	if response.StatusCode == net_http.StatusOK {
		if result, err := json.ConvertFromString[SubscriptionResponse](response.Body); err != nil {
			return SubscriptionResponse{}, err
		} else {
			subscriptionResponse.Events = result.Events
		}
	}

	return subscriptionResponse, nil
}

// Publish publishes an event to the long polling server.
//
// This function sends an event to the server, which distributes it to all
// subscribed clients in the specified category.
//
// # Parameters
//
//   - url: Server publish endpoint URL (e.g., "http://localhost:8080/publish")
//   - timeout: Maximum duration to wait for publish operation
//   - header: Optional HTTP headers (e.g., custom headers or authorization)
//   - publishRequest: Event data (category and data payload)
//   - username: Optional HTTP basic authentication username
//   - password: Optional HTTP basic authentication password
//   - transport: Optional custom HTTP transport for connection pooling or proxies
//
// # Returns
//
//   - http.Response: HTTP response with headers, status code, and body
//   - error: Error if request fails, nil on successful publish
//
// # Behavior
//
// The publish request:
//   - Sends event to server via POST request
//   - Server assigns unique event ID and timestamp
//   - Event is queued for subscribed clients
//   - Returns immediately (non-blocking)
//
// # Examples
//
// Publish notification:
//
//	response, err := long_polling.Publish(
//	    "http://localhost:8080/publish",
//	    10 * time.Second,
//	    nil,
//	    long_polling.PublishRequest{
//	        Category: "notifications",
//	        Data: `{"message": "Hello"}`
//	    },
//	    "", "", nil,
//	)
func Publish(url string, timeout time.Duration, header map[string][]string, publishRequest PublishRequest, username, password string, transport *net_http.Transport) (http.Response, error) {
	u, err := net_url.Parse(url)
	if err != nil {
		return http.Response{}, err
	}

	body, err := json.ToString(publishRequest)
	if err != nil {
		return http.Response{}, err
	}

	return http.Request(fmt.Sprintf("%v", u), net_http.MethodPost, header, body, timeout, username, password, transport)
}
