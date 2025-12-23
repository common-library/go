// Package client provides a wrapper around the official Prometheus Go client library
// for querying Prometheus metrics with support for basic authentication, bearer token
// authentication, and PromQL queries.
//
// Features:
//   - Simple client creation with various authentication methods
//   - Instant vector queries with Query()
//   - Range vector queries with QueryRange()
//   - Configurable timeout for all operations
//   - Support for basic auth and bearer token authentication
//
// Example:
//
//	client, err := client.NewClient("http://localhost:9090")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	value, warnings, err := client.Query("up", time.Now(), 10*time.Second)
//	if err != nil {
//	    log.Fatal(err)
//	}
package client

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
)

type Range = v1.Range

// NewClient creates a new Prometheus client with the specified server address.
//
// Parameters:
//   - address: Prometheus server URL (e.g., "http://localhost:9090")
//
// Returns:
//   - *client: Configured Prometheus client
//   - error: Error if client creation fails
//
// Example:
//
//	c, err := client.NewClient("http://localhost:9090")
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewClient(address string) (*client, error) {
	if prometheusClient, err := api.NewClient(api.Config{Address: address}); err != nil {
		return nil, err
	} else {
		return &client{prometheusClient: prometheusClient}, nil
	}
}

// NewClientWithBasicAuth creates a new Prometheus client with HTTP basic authentication.
//
// Parameters:
//   - address: Prometheus server URL (e.g., "http://localhost:9090")
//   - username: Basic auth username
//   - password: Basic auth password
//
// Returns:
//   - *client: Configured Prometheus client with basic auth
//   - error: Error if client creation fails
//
// Example:
//
//	c, err := client.NewClientWithBasicAuth(
//	    "http://localhost:9090",
//	    "admin",
//	    "secret",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewClientWithBasicAuth(address, username string, password string) (*client, error) {
	config := api.Config{
		Address:      address,
		RoundTripper: config.NewBasicAuthRoundTripper(config.NewInlineSecret(username), config.NewInlineSecret(password), api.DefaultRoundTripper),
	}

	if prometheusClient, err := api.NewClient(config); err != nil {
		return nil, err
	} else {
		return &client{prometheusClient: prometheusClient}, nil
	}
}

// NewClientWithBearerToken creates a new Prometheus client with bearer token authentication.
//
// Parameters:
//   - address: Prometheus server URL (e.g., "http://localhost:9090")
//   - token: Bearer token for authentication
//
// Returns:
//   - *client: Configured Prometheus client with bearer token auth
//   - error: Error if client creation fails
//
// Example:
//
//	c, err := client.NewClientWithBearerToken(
//	    "http://localhost:9090",
//	    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewClientWithBearerToken(address string, token string) (*client, error) {
	config := api.Config{
		Address:      address,
		RoundTripper: config.NewAuthorizationCredentialsRoundTripper("Bearer", config.NewInlineSecret(token), api.DefaultRoundTripper),
	}

	if prometheusClient, err := api.NewClient(config); err != nil {
		return nil, err
	} else {
		return &client{prometheusClient: prometheusClient}, nil
	}
}

// client is a struct that provides client related methods.
type client struct {
	prometheusClient api.Client
}

// Query executes an instant PromQL query at a specific point in time.
//
// This method performs a query evaluation at a single timestamp, useful for
// retrieving the current state of metrics or evaluating expressions at a
// specific moment.
//
// Parameters:
//   - query: PromQL query string (e.g., "up", "rate(http_requests_total[5m])")
//   - when: Timestamp for query evaluation
//   - timeout: Maximum duration to wait for query completion
//
// Returns:
//   - model.Value: Query result (can be scalar, vector, matrix, or string)
//   - v1.Warnings: Warnings returned by Prometheus
//   - error: Error if query fails or times out
//
// Example:
//
//	c, err := client.NewClient("http://localhost:9090")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	value, warnings, err := c.Query("up", time.Now(), 10*time.Second)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, warning := range warnings {
//	    log.Println("Warning:", warning)
//	}
func (c *client) Query(query string, when time.Time, timeout time.Duration) (model.Value, v1.Warnings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return v1.NewAPI(c.prometheusClient).Query(ctx, query, when, v1.WithTimeout(timeout))
}

// QueryRange executes a PromQL query over a time range with a specified step interval.
//
// This method performs a range query evaluation, useful for retrieving time series
// data over a period of time. The query is evaluated at regular intervals (steps)
// between the start and end times.
//
// Parameters:
//   - query: PromQL query string (e.g., "rate(http_requests_total[5m])")
//   - r: Time range with start, end, and step interval
//   - timeout: Maximum duration to wait for query completion
//
// Returns:
//   - model.Value: Query result as a matrix with time series data
//   - v1.Warnings: Warnings returned by Prometheus
//   - error: Error if query fails or times out
//
// Example:
//
//	c, err := client.NewClient("http://localhost:9090")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	now := time.Now()
//	r := client.Range{
//	    Start: now.Add(-1 * time.Hour),
//	    End:   now,
//	    Step:  time.Minute,
//	}
//
//	value, warnings, err := c.QueryRange(
//	    "rate(process_cpu_seconds_total[5m])",
//	    r,
//	    10*time.Second,
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *client) QueryRange(query string, r v1.Range, timeout time.Duration) (model.Value, v1.Warnings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return v1.NewAPI(c.prometheusClient).QueryRange(ctx, query, r)
}
