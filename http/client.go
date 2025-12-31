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
	"io"
	net_http "net/http"
	"strings"
	"time"
)

// Response is response information.
type Response struct {
	Header     net_http.Header
	Body       string
	StatusCode int
}

// Request performs an HTTP request and returns the response.
//
// Parameters:
//   - url: Target URL for the HTTP request
//   - method: HTTP method (http.MethodGet, http.MethodPost, etc.)
//   - header: HTTP headers as a map of header names to value slices
//   - body: Request body as a string
//   - timeout: Request timeout duration (e.g., 10*time.Second)
//   - username: Username for HTTP Basic Authentication (empty string for no auth)
//   - password: Password for HTTP Basic Authentication (empty string for no auth)
//   - transport: Custom HTTP transport configuration (nil for default)
//
// Returns:
//   - Response: HTTP response containing headers, body, and status code
//   - error: Error if request fails, nil on success
//
// The function creates an HTTP client with the specified timeout and transport settings,
// performs the request, and returns the complete response. The response body is read
// entirely into memory.
//
// Example:
//
//	// Simple GET request
//	resp, err := http.Request(
//	    "https://api.example.com/users",
//	    http.MethodGet,
//	    nil,              // no headers
//	    "",               // no body
//	    10*time.Second,   // 10 second timeout
//	    "",               // no username
//	    "",               // no password
//	    nil,              // default transport
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Status: %d\n", resp.StatusCode)
//	fmt.Printf("Body: %s\n", resp.Body)
//
//	// POST request with headers and authentication
//	headers := map[string][]string{
//	    "Content-Type": {"application/json"},
//	    "X-API-Key":    {"secret123"},
//	}
//	body := `{"name": "Alice", "email": "alice@example.com"}`
//	resp, err = http.Request(
//	    "https://api.example.com/users",
//	    http.MethodPost,
//	    headers,
//	    body,
//	    30*time.Second,
//	    "admin",
//	    "password123",
//	    nil,
//	)
func Request(url, method string, header map[string][]string, body string, timeout time.Duration, username, password string, transport *net_http.Transport) (Response, error) {
	if request, err := getRequest(url, method, header, body, username, password); err != nil {
		return Response{}, err
	} else {
		return getResponse(request, timeout, transport)
	}
}

func getRequest(url, method string, header map[string][]string, body string, username, password string) (*net_http.Request, error) {
	if request, err := net_http.NewRequest(method, url, strings.NewReader(body)); err != nil {
		return nil, err
	} else {
		if username != "" && password != "" {
			request.SetBasicAuth(username, password)
		}

		for key, array := range header {
			for _, value := range array {
				request.Header.Add(key, value)
			}
		}

		return request, nil
	}
}

func getResponse(request *net_http.Request, timeout time.Duration, transport *net_http.Transport) (Response, error) {
	if transport == nil {
		transport = &net_http.Transport{}
	}

	client := net_http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	if response, err := client.Do(request); err != nil {
		return Response{}, err
	} else {
		defer response.Body.Close()

		if responseBody, err := io.ReadAll(response.Body); err != nil {
			return Response{}, err
		} else {
			return Response{Header: response.Header, Body: string(responseBody), StatusCode: response.StatusCode}, nil
		}
	}
}
