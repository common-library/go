package exporter_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/common-library/go/database/prometheus/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func waitForServerReady(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("server not ready within %v", timeout)
}

func waitForServerShutdown(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		_, err := http.Get(url)
		if err != nil {
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("server still running after %v", timeout)
}

func waitForPrometheusReady(endpoint string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get("http://" + endpoint + "/-/ready")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("prometheus not ready within %v", timeout)
}

type TestExporterMetric struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	values    []exporter.Value
}

func (tm *TestExporterMetric) GetDesc() *prometheus.Desc {
	return tm.desc
}

func (tm *TestExporterMetric) GetValueType() prometheus.ValueType {
	return tm.valueType
}

func (tm *TestExporterMetric) GetValues() []exporter.Value {
	return tm.values
}

func TestExporterNewCollector(t *testing.T) {
	desc := prometheus.NewDesc(
		"test_exporter_metric",
		"A test exporter metric",
		[]string{"label1"},
		prometheus.Labels{"const": "value"},
	)

	metrics := []exporter.Metric{
		&TestExporterMetric{
			desc:      desc,
			valueType: prometheus.CounterValue,
			values: []exporter.Value{
				{Value: 42.0, LabelValues: []string{"test"}},
			},
		},
	}

	collector := exporter.NewCollector(metrics)
	assert.NotNil(t, collector)

	var _ prometheus.Collector = collector
}

func TestExporterRegisterCollector(t *testing.T) {
	registry := prometheus.NewRegistry()

	desc := prometheus.NewDesc(
		"test_register_exporter_metric",
		"A test metric for registration",
		[]string{"label1"},
		nil,
	)

	metrics := []exporter.Metric{
		&TestExporterMetric{
			desc:      desc,
			valueType: prometheus.CounterValue,
			values: []exporter.Value{
				{Value: 1.0, LabelValues: []string{"test"}},
			},
		},
	}

	collector := exporter.NewCollector(metrics)

	err := registry.Register(collector)
	assert.NoError(t, err)

	err = registry.Register(collector)
	assert.Error(t, err)

	registry.Unregister(collector)
}

func TestExporterUnRegisterCollector(t *testing.T) {
	registry := prometheus.NewRegistry()

	desc := prometheus.NewDesc(
		"test_unregister_exporter_metric",
		"A test metric for unregistration",
		[]string{"label1"},
		nil,
	)

	metrics := []exporter.Metric{
		&TestExporterMetric{
			desc:      desc,
			valueType: prometheus.CounterValue,
			values: []exporter.Value{
				{Value: 1.0, LabelValues: []string{"test"}},
			},
		},
	}

	collector := exporter.NewCollector(metrics)

	err := registry.Register(collector)
	assert.NoError(t, err)

	result := registry.Unregister(collector)
	assert.True(t, result)

	result = registry.Unregister(collector)
	assert.False(t, result)
}

func TestExporterStartAndStopServer(t *testing.T) {
	address := ":18080"
	urlPath := "/test-metrics"

	go func() {
		_ = exporter.Start(address, urlPath, func(err error) {
			if err != nil && !strings.Contains(err.Error(), "Server closed") {
				t.Logf("Server error: %v", err)
			}
		})
	}()

	serverURL := "http://localhost:18080" + urlPath
	err := waitForServerReady(serverURL, 5*time.Second)
	require.NoError(t, err, "Server should start successfully")

	resp, err := http.Get(serverURL)
	if err == nil {
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	err = exporter.Stop(5 * time.Second)
	assert.NoError(t, err)

	err = waitForServerShutdown(serverURL, 2*time.Second)
	assert.NoError(t, err, "Server should shutdown successfully")
}

func TestExporterMetricsEndpointWithPrometheus(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "prom/prometheus:v3.6.0",
		ExposedPorts: []string{"9090/tcp"},
		WaitingFor:   wait.ForListeningPort("9090/tcp"),
		Cmd:          []string{"--config.file=/etc/prometheus/prometheus.yml", "--storage.tsdb.path=/prometheus", "--web.console.libraries=/etc/prometheus/console_libraries", "--web.console.templates=/etc/prometheus/consoles", "--web.enable-lifecycle"},
	}

	prometheusContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, prometheusContainer.Terminate(ctx))
	}()

	endpoint, err := prometheusContainer.Endpoint(ctx, "")
	require.NoError(t, err)

	err = waitForPrometheusReady(endpoint, 10*time.Second)
	require.NoError(t, err, "Prometheus should be ready")

	resp, err := http.Get("http://" + endpoint + "/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Get("http://" + endpoint + "/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestExporterMetricsCollection(t *testing.T) {
	registry := prometheus.NewRegistry()

	counterDesc := prometheus.NewDesc(
		"test_counter_total",
		"A test counter",
		[]string{"method", "status"},
		nil,
	)

	gaugeDesc := prometheus.NewDesc(
		"test_gauge_current",
		"A test gauge",
		[]string{"instance"},
		nil,
	)

	metrics := []exporter.Metric{
		&TestExporterMetric{
			desc:      counterDesc,
			valueType: prometheus.CounterValue,
			values: []exporter.Value{
				{Value: 100.0, LabelValues: []string{"GET", "200"}},
				{Value: 50.0, LabelValues: []string{"POST", "201"}},
				{Value: 10.0, LabelValues: []string{"GET", "404"}},
			},
		},
		&TestExporterMetric{
			desc:      gaugeDesc,
			valueType: prometheus.GaugeValue,
			values: []exporter.Value{
				{Value: 85.5, LabelValues: []string{"server1"}},
				{Value: 92.3, LabelValues: []string{"server2"}},
			},
		},
	}

	collector := exporter.NewCollector(metrics)
	err := registry.Register(collector)
	require.NoError(t, err)

	metricFamilies, err := registry.Gather()
	require.NoError(t, err)

	assert.Len(t, metricFamilies, 2)

	var counterFamily *dto.MetricFamily
	var gaugeFamily *dto.MetricFamily

	for _, mf := range metricFamilies {
		switch *mf.Name {
		case "test_counter_total":
			counterFamily = mf
		case "test_gauge_current":
			gaugeFamily = mf
		}
	}

	require.NotNil(t, counterFamily)
	require.NotNil(t, gaugeFamily)

	assert.Len(t, counterFamily.Metric, 3)
	assert.Equal(t, dto.MetricType_COUNTER, *counterFamily.Type)

	assert.Len(t, gaugeFamily.Metric, 2)
	assert.Equal(t, dto.MetricType_GAUGE, *gaugeFamily.Type)

	registry.Unregister(collector)
}

func TestExporterMetricsHandler(t *testing.T) {
	registry := prometheus.NewRegistry()

	desc := prometheus.NewDesc(
		"test_http_requests_total",
		"Total HTTP requests",
		[]string{"method"},
		nil,
	)

	metrics := []exporter.Metric{
		&TestExporterMetric{
			desc:      desc,
			valueType: prometheus.CounterValue,
			values: []exporter.Value{
				{Value: 123.0, LabelValues: []string{"GET"}},
				{Value: 456.0, LabelValues: []string{"POST"}},
			},
		},
	}

	collector := exporter.NewCollector(metrics)
	err := registry.Register(collector)
	require.NoError(t, err)

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	server := &http.Server{
		Addr:    ":18081",
		Handler: handler,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	serverURL := "http://localhost:18081"
	err = waitForServerReady(serverURL, 5*time.Second)
	require.NoError(t, err, "Server should start successfully")

	resp, err := http.Get(serverURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Contains(t, resp.Header.Get("Content-Type"), "text/plain")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	assert.NoError(t, err)

	registry.Unregister(collector)
}

func TestExporterErrorHandling(t *testing.T) {
	err := exporter.Start("invalid-address", "/metrics", func(err error) {
		t.Logf("Expected error: %v", err)
	})
	_ = err

	err = exporter.Stop(1 * time.Second)

	_ = err
}

func BenchmarkExporterCollector(b *testing.B) {
	desc := prometheus.NewDesc(
		"benchmark_exporter_metric",
		"A benchmark metric",
		[]string{"label1", "label2"},
		nil,
	)

	var values []exporter.Value
	for i := 0; i < 1000; i++ {
		values = append(values, exporter.Value{
			Value:       float64(i),
			LabelValues: []string{fmt.Sprintf("label_%d", i), fmt.Sprintf("value_%d", i%10)},
		})
	}

	metrics := []exporter.Metric{
		&TestExporterMetric{
			desc:      desc,
			valueType: prometheus.CounterValue,
			values:    values,
		},
	}

	collector := exporter.NewCollector(metrics)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metricChan := make(chan prometheus.Metric, len(values))
		collector.Collect(metricChan)
		close(metricChan)

		for range metricChan {
		}
	}
}
