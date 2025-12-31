// Package exporter provides a framework for creating Prometheus exporters that expose
// custom metrics through an HTTP endpoint.
package exporter

import (
	"github.com/common-library/go/http/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var server mux.Server

// Metric is an interface that must be implemented by custom metrics to be collected
// by a Prometheus collector.
//
// Implementations should provide:
//   - Metric description (name, help text, labels)
//   - Value type (counter, gauge, histogram, etc.)
//   - Current metric values with their label values
//
// Example implementation:
//
//	type MyMetric struct {
//	    desc      *prometheus.Desc
//	    valueType prometheus.ValueType
//	}
//
//	func (m *MyMetric) GetDesc() *prometheus.Desc {
//	    return m.desc
//	}
//
//	func (m *MyMetric) GetValueType() prometheus.ValueType {
//	    return m.valueType
//	}
//
//	func (m *MyMetric) GetValues() []Value {
//	    // Fetch current metric values from your data source
//	    return []Value{
//	        {Value: 42.0, LabelValues: []string{"label1_value"}},
//	    }
//	}
type Metric interface {
	// GetDesc returns the metric descriptor containing name, help text, and labels
	GetDesc() *prometheus.Desc

	// GetValueType returns the Prometheus value type (Counter, Gauge, Histogram, etc.)
	GetValueType() prometheus.ValueType

	// GetValues returns the current metric values with their corresponding label values
	GetValues() []Value
}

// Value represents a single metric value with its associated label values.
//
// This struct is used to represent individual data points for a metric.
// Each value must have label values that correspond to the labels defined
// in the metric's descriptor.
//
// Fields:
//   - Value: The numeric value of the metric
//   - LabelValues: Slice of label values in the same order as labels in the descriptor
//
// Example:
//
//	// For a metric with labels ["instance", "job"]
//	value := Value{
//	    Value:       42.0,
//	    LabelValues: []string{"server1", "api"},
//	}
type Value struct {
	// Value is the numeric metric value
	Value float64

	// LabelValues contains the label values for this metric instance,
	// in the same order as defined in the prometheus.Desc
	LabelValues []string
}

type collector struct {
	metrics []Metric
}

func (t *collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range t.metrics {
		ch <- metric.GetDesc()
	}
}

func (t *collector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range t.metrics {
		for _, value := range metric.GetValues() {
			ch <- prometheus.MustNewConstMetric(metric.GetDesc(), metric.GetValueType(), value.Value, value.LabelValues...)
		}
	}
}
