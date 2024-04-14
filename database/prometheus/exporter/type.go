// Package exporter provides prometheus exporter implementations.
package exporter

import (
	"github.com/common-library/go/http"
	"github.com/prometheus/client_golang/prometheus"
)

var server http.Server

// Metric is an interface that provides the information to be collected.
type Metric interface {
	GetDesc() *prometheus.Desc
	GetValueType() prometheus.ValueType
	GetValues() []Value
}

// value is a struct that provides the value to collect.
type Value struct {
	Value       float64
	LabelValues []string
}

type collector struct {
	metrics []Metric
}

func (this *collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range this.metrics {
		ch <- metric.GetDesc()
	}
}

func (this *collector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range this.metrics {
		for _, value := range metric.GetValues() {
			ch <- prometheus.MustNewConstMetric(metric.GetDesc(), metric.GetValueType(), value.Value, value.LabelValues...)
		}
	}
}
