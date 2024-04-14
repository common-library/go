// Package exporter provides prometheus exporter implementations.
package exporter

import (
	net_http "net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// New is creates a Collector.
//
// ex) collector01 := exporter.New([]exporter.Metric{...})
func New(metrics []Metric) prometheus.Collector {
	return &collector{metrics: metrics}
}

// Register registers the Collectors.
//
// ex) err := exporter.Register(collector01)
func Register(collectors ...prometheus.Collector) error {
	for _, collector := range collectors {
		if err := prometheus.Register(collector); err != nil {
			return err
		}
	}

	return nil
}

// UnRegister unregister the Collectors.
//
// ex) result := exporter.UnRegister(collector01)
func UnRegister(collectors ...prometheus.Collector) bool {
	for _, collector := range collectors {
		if prometheus.Unregister(collector) == false {
			return false
		}
	}

	return true
}

// Start is start the server.
//
// ex) err := exporter.Start(":10000", "metrics", func(err error) { klog.ErrorS(err, "") })
func Start(address, urlPath string, listenAndServeFailureFunc func(error)) error {
	server.RegisterHandlerFunc(urlPath, net_http.MethodGet, promhttp.Handler().ServeHTTP)

	return server.Start(address, listenAndServeFailureFunc)
}

// Stop is stop the server.
//
// ex) err := exporter.Stop(60)
func Stop(timeout time.Duration) error {
	return server.Stop(timeout)
}
