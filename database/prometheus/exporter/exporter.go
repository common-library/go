// Package exporter provides a framework for creating Prometheus exporters that expose
// custom metrics through an HTTP endpoint. It simplifies the process of creating,
// registering, and serving Prometheus metrics.
//
// Features:
//   - Easy collector creation with custom metrics
//   - Automatic metric collection and exposition
//   - Built-in HTTP server for /metrics endpoint
//   - Support for multiple collectors
//   - Graceful server shutdown
//
// Example:
//
//	// Define your metric
//	desc := prometheus.NewDesc(
//	    "my_metric",
//	    "My custom metric",
//	    []string{"label"},
//	    nil,
//	)
//
//	// Create collector with metrics
//	collector := exporter.NewCollector([]exporter.Metric{myMetric})
//
//	// Register and start server
//	exporter.RegisterCollector(collector)
//	exporter.Start(":9090", "/metrics", func(err error) {
//	    log.Println(err)
//	})
package exporter

import (
	net_http "net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewCollector creates a new Prometheus collector from a slice of metrics.
//
// A collector implements the prometheus.Collector interface and is responsible
// for gathering metrics and exposing them to Prometheus. Each metric in the
// provided slice should implement the Metric interface.
//
// Parameters:
//   - metrics: Slice of Metric implementations to be collected
//
// Returns:
//   - prometheus.Collector: A collector that can be registered with Prometheus
//
// Example:
//
//	desc := prometheus.NewDesc(
//	    "my_custom_metric",
//	    "Description of my metric",
//	    []string{"label1"},
//	    nil,
//	)
//
//	metrics := []exporter.Metric{
//	    &MyMetric{
//	        desc:      desc,
//	        valueType: prometheus.GaugeValue,
//	        values:    []exporter.Value{{Value: 42.0, LabelValues: []string{"test"}}},
//	    },
//	}
//
//	collector := exporter.NewCollector(metrics)
func NewCollector(metrics []Metric) prometheus.Collector {
	return &collector{metrics: metrics}
}

// RegisterCollector registers one or more collectors with the default Prometheus registry.
//
// Once registered, the collectors' metrics will be included in the exposition
// format when the /metrics endpoint is scraped. If any collector fails to
// register (e.g., duplicate metrics), an error is returned and registration stops.
//
// Parameters:
//   - collectors: One or more prometheus.Collector instances to register
//
// Returns:
//   - error: Error if registration fails (e.g., duplicate collector)
//
// Example:
//
//	collector1 := exporter.NewCollector(metrics1)
//	collector2 := exporter.NewCollector(metrics2)
//
//	err := exporter.RegisterCollector(collector1, collector2)
//	if err != nil {
//	    log.Fatal(err)
//	}
func RegisterCollector(collectors ...prometheus.Collector) error {
	for _, collector := range collectors {
		if err := prometheus.Register(collector); err != nil {
			return err
		}
	}

	return nil
}

// UnRegisterCollector unregisters one or more collectors from the default Prometheus registry.
//
// After unregistration, the collectors' metrics will no longer be included in
// the exposition format. Returns false if any collector was not registered.
//
// Parameters:
//   - collectors: One or more prometheus.Collector instances to unregister
//
// Returns:
//   - bool: true if all collectors were successfully unregistered, false otherwise
//
// Example:
//
//	collector := exporter.NewCollector(metrics)
//	exporter.RegisterCollector(collector)
//
//	// Later, when metrics are no longer needed
//	result := exporter.UnRegisterCollector(collector)
//	if !result {
//	    log.Println("Failed to unregister collector")
//	}
func UnRegisterCollector(collectors ...prometheus.Collector) bool {
	for _, collector := range collectors {
		if !prometheus.Unregister(collector) {
			return false
		}
	}

	return true
}

// Start starts the HTTP server that exposes Prometheus metrics at the specified endpoint.
//
// The server runs in the current goroutine and blocks until it's stopped via the
// Stop function or encounters an error. The metrics handler is automatically
// registered at the specified URL path.
//
// Parameters:
//   - address: Server bind address (e.g., ":9090" or "0.0.0.0:9090")
//   - urlPath: URL path for metrics endpoint (e.g., "/metrics")
//   - listenAndServeFailureFunc: Callback function invoked when server fails to start
//
// Returns:
//   - error: Error if server fails to start or configuration is invalid
//
// Example:
//
//	// Start server in a goroutine
//	go func() {
//	    err := exporter.Start(
//	        ":9090",
//	        "/metrics",
//	        func(err error) {
//	            log.Printf("Server error: %v", err)
//	        },
//	    )
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	}()
//
//	// Wait for server to be ready
//	time.Sleep(100 * time.Millisecond)
func Start(address, urlPath string, listenAndServeFailureFunc func(error)) error {
	server.RegisterHandler(urlPath, promhttp.Handler(), net_http.MethodGet)

	return server.Start(address, listenAndServeFailureFunc)
}

// Stop gracefully shuts down the HTTP server with a timeout.
//
// The server will attempt to complete all active requests within the timeout
// period before shutting down. If the timeout is exceeded, the server is
// forcefully closed.
//
// Parameters:
//   - timeout: Maximum duration to wait for graceful shutdown
//
// Returns:
//   - error: Error if shutdown fails or times out
//
// Example:
//
//	// Graceful shutdown with 30 second timeout
//	err := exporter.Stop(30 * time.Second)
//	if err != nil {
//	    log.Printf("Error during shutdown: %v", err)
//	}
func Stop(timeout time.Duration) error {
	return server.Stop(timeout)
}
