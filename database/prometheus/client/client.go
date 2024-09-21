// Package client provides prometheus client implementations.
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

// NewClient creates a client.
//
// ex) c, err := client.NewClient("http://:9090")
func NewClient(address string) (*client, error) {
	if prometheusClient, err := api.NewClient(api.Config{Address: address}); err != nil {
		return nil, err
	} else {
		return &client{prometheusClient: prometheusClient}, nil
	}
}

// NewClientWithBasicAuth creates a client with basic authentication.
//
// ex) c, err := client.NewClientWithBasicAuth("http://:9090", "username", "password")
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

// NewClientWithBearerToken creates a client that performs bearer token authentication.
//
// ex) c, err := client.NewClientWithBearerToken("http://:9090", "token")
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

// Queries perform PromQL queries over a given period of time.
//
// ex)
//
//	c, err := client.NewClient(address)
//	value, warnings, err := c.Query("up", time.Now(), 10*time.Second)
func (this *client) Query(query string, when time.Time, timeout time.Duration) (model.Value, v1.Warnings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return v1.NewAPI(this.prometheusClient).Query(ctx, query, when, v1.WithTimeout(timeout))
}

// QueryRange perform PromQL queries over a given period of range.
//
// ex)
//
//	c, err := client.NewClient(address)
//	value, warnings, err := c.QueryRange("rate(process_cpu_seconds_total[5m])", r, 10*time.Second)
func (this *client) QueryRange(query string, r v1.Range, timeout time.Duration) (model.Value, v1.Warnings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return v1.NewAPI(this.prometheusClient).QueryRange(ctx, query, r)
}
