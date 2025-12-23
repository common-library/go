# Prometheus

Client and exporter utilities for Prometheus monitoring system.

## Overview

The prometheus package provides two main components for working with Prometheus:

1. **Client**: Query Prometheus servers using PromQL with support for authentication
2. **Exporter**: Create custom exporters to expose application metrics to Prometheus

Both components are built on top of the official Prometheus Go client library, providing a simplified and more ergonomic API for common use cases.

## Package Structure

```
prometheus/
├── client/          # Prometheus query client
│   ├── client.go    # Client implementation
│   └── client_test.go
└── exporter/        # Custom metrics exporter
    ├── exporter.go  # Server and registration
    ├── type.go      # Metric interfaces
    ├── exporter_test.go
    └── type_test.go
```

## Installation

```bash
go get github.com/common-library/go/database/prometheus
```

## Client

### Overview

The client package provides a wrapper around the Prometheus HTTP API for executing PromQL queries. It supports:

- Instant queries (single point in time)
- Range queries (time series over a period)
- Basic authentication
- Bearer token authentication
- Configurable timeouts

### Quick Start

```go
import (
    "log"
    "time"
    
    "github.com/common-library/go/database/prometheus/client"
)

func main() {
    // Create client
    c, err := client.NewClient("http://localhost:9090")
    if err != nil {
        log.Fatal(err)
    }
    
    // Execute instant query
    value, warnings, err := c.Query("up", time.Now(), 10*time.Second)
    if err != nil {
        log.Fatal(err)
    }
    
    // Handle warnings
    for _, warning := range warnings {
        log.Println("Warning:", warning)
    }
    
    log.Printf("Result: %v", value)
}
```

### Authentication

#### Basic Authentication

```go
client, err := client.NewClientWithBasicAuth(
    "http://localhost:9090",
    "admin",
    "secretpassword",
)
if err != nil {
    log.Fatal(err)
}
```

#### Bearer Token

```go
client, err := client.NewClientWithBearerToken(
    "http://localhost:9090",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
)
if err != nil {
    log.Fatal(err)
}
```

### Query Operations

#### Instant Query

Execute a PromQL query at a specific point in time:

```go
// Query current state
value, warnings, err := c.Query(
    "up",                    // PromQL query
    time.Now(),              // Evaluation time
    10*time.Second,          // Timeout
)

// Query with complex expression
value, warnings, err = c.Query(
    "rate(http_requests_total[5m])",
    time.Now(),
    10*time.Second,
)

// Query with aggregation
value, warnings, err = c.Query(
    "sum(rate(http_requests_total[5m])) by (job)",
    time.Now(),
    10*time.Second,
)
```

#### Range Query

Execute a PromQL query over a time range:

```go
now := time.Now()
r := client.Range{
    Start: now.Add(-1 * time.Hour),  // Start time
    End:   now,                       // End time
    Step:  time.Minute,               // Evaluation step
}

value, warnings, err := c.QueryRange(
    "rate(process_cpu_seconds_total[5m])",
    r,
    10*time.Second,
)
```

### Result Processing

```go
import (
    "github.com/prometheus/common/model"
)

value, warnings, err := c.Query("up", time.Now(), 10*time.Second)
if err != nil {
    log.Fatal(err)
}

switch v := value.(type) {
case model.Vector:
    // Instant vector
    for _, sample := range v {
        log.Printf("Metric: %s, Value: %f, Time: %v",
            sample.Metric, sample.Value, sample.Timestamp)
    }
    
case model.Matrix:
    // Range vector
    for _, stream := range v {
        log.Printf("Metric: %s", stream.Metric)
        for _, value := range stream.Values {
            log.Printf("  Value: %f, Time: %v", value.Value, value.Timestamp)
        }
    }
    
case *model.Scalar:
    // Scalar value
    log.Printf("Scalar: %f", v.Value)
    
case *model.String:
    // String value
    log.Printf("String: %s", v.Value)
}
```

### Error Handling

```go
value, warnings, err := c.Query("invalid_query{", time.Now(), 10*time.Second)
if err != nil {
    // Handle query errors
    log.Printf("Query failed: %v", err)
    return
}

// Check warnings
if len(warnings) > 0 {
    for _, warning := range warnings {
        log.Printf("Warning: %s", warning)
    }
}
```

### Best Practices

1. **Set Appropriate Timeouts**: Adjust timeout based on query complexity
2. **Handle Warnings**: Always check and log warnings
3. **Reuse Clients**: Create client once and reuse for multiple queries
4. **Error Handling**: Always check errors before using results
5. **Query Optimization**: Use appropriate time ranges and aggregations

## Exporter

### Overview

The exporter package provides a framework for creating custom Prometheus exporters. It simplifies:

- Metric definition and collection
- Collector registration
- HTTP server setup for /metrics endpoint
- Graceful shutdown

### Quick Start

```go
import (
    "log"
    "time"
    
    "github.com/common-library/go/database/prometheus/exporter"
    "github.com/prometheus/client_golang/prometheus"
)

// 1. Implement the Metric interface
type MyMetric struct {
    desc      *prometheus.Desc
    valueType prometheus.ValueType
}

func (m *MyMetric) GetDesc() *prometheus.Desc {
    return m.desc
}

func (m *MyMetric) GetValueType() prometheus.ValueType {
    return m.valueType
}

func (m *MyMetric) GetValues() []exporter.Value {
    // Fetch current values from your application
    return []exporter.Value{
        {
            Value:       42.0,
            LabelValues: []string{"instance1"},
        },
    }
}

func main() {
    // 2. Create metric descriptor
    desc := prometheus.NewDesc(
        "my_application_requests_total",
        "Total number of requests",
        []string{"instance"},
        prometheus.Labels{"version": "1.0"},
    )
    
    // 3. Create metric instance
    metric := &MyMetric{
        desc:      desc,
        valueType: prometheus.CounterValue,
    }
    
    // 4. Create collector
    collector := exporter.NewCollector([]exporter.Metric{metric})
    
    // 5. Register collector
    err := exporter.RegisterCollector(collector)
    if err != nil {
        log.Fatal(err)
    }
    
    // 6. Start HTTP server
    go func() {
        err := exporter.Start(
            ":9090",
            "/metrics",
            func(err error) {
                log.Printf("Server error: %v", err)
            },
        )
        if err != nil {
            log.Fatal(err)
        }
    }()
    
    // Keep running
    select {}
}
```

### Metric Types

#### Counter

```go
desc := prometheus.NewDesc(
    "requests_total",
    "Total number of requests",
    []string{"method", "status"},
    nil,
)

type RequestCounter struct {
    desc *prometheus.Desc
}

func (rc *RequestCounter) GetDesc() *prometheus.Desc {
    return rc.desc
}

func (rc *RequestCounter) GetValueType() prometheus.ValueType {
    return prometheus.CounterValue
}

func (rc *RequestCounter) GetValues() []exporter.Value {
    return []exporter.Value{
        {Value: 100.0, LabelValues: []string{"GET", "200"}},
        {Value: 50.0, LabelValues: []string{"POST", "200"}},
        {Value: 5.0, LabelValues: []string{"GET", "404"}},
    }
}
```

#### Gauge

```go
desc := prometheus.NewDesc(
    "active_connections",
    "Number of active connections",
    []string{"service"},
    nil,
)

type ConnectionGauge struct {
    desc *prometheus.Desc
}

func (cg *ConnectionGauge) GetDesc() *prometheus.Desc {
    return cg.desc
}

func (cg *ConnectionGauge) GetValueType() prometheus.ValueType {
    return prometheus.GaugeValue
}

func (cg *ConnectionGauge) GetValues() []exporter.Value {
    // Get current connection count from your application
    activeCount := getCurrentConnectionCount()
    
    return []exporter.Value{
        {Value: float64(activeCount), LabelValues: []string{"api"}},
    }
}
```

#### Histogram

```go
desc := prometheus.NewDesc(
    "request_duration_seconds",
    "Request duration in seconds",
    []string{"endpoint"},
    nil,
)

type DurationHistogram struct {
    desc *prometheus.Desc
}

func (dh *DurationHistogram) GetDesc() *prometheus.Desc {
    return dh.desc
}

func (dh *DurationHistogram) GetValueType() prometheus.ValueType {
    return prometheus.UntypedValue
}

func (dh *DurationHistogram) GetValues() []exporter.Value {
    // Collect duration samples
    durations := collectRequestDurations()
    
    values := make([]exporter.Value, 0, len(durations))
    for endpoint, duration := range durations {
        values = append(values, exporter.Value{
            Value:       duration,
            LabelValues: []string{endpoint},
        })
    }
    return values
}
```

### Multiple Metrics

```go
// Define multiple metrics
requestMetric := &RequestCounter{desc: requestDesc}
connectionMetric := &ConnectionGauge{desc: connectionDesc}
durationMetric := &DurationHistogram{desc: durationDesc}

// Create collector with all metrics
collector := exporter.NewCollector([]exporter.Metric{
    requestMetric,
    connectionMetric,
    durationMetric,
})

// Register
err := exporter.RegisterCollector(collector)
```

### Multiple Collectors

```go
// Create separate collectors for different subsystems
apiCollector := exporter.NewCollector(apiMetrics)
dbCollector := exporter.NewCollector(dbMetrics)
cacheCollector := exporter.NewCollector(cacheMetrics)

// Register all collectors
err := exporter.RegisterCollector(
    apiCollector,
    dbCollector,
    cacheCollector,
)
if err != nil {
    log.Fatal(err)
}
```

### Server Management

#### Start Server

```go
// Start in goroutine
go func() {
    err := exporter.Start(
        ":9090",        // Bind address
        "/metrics",     // URL path
        func(err error) {
            if err != nil && !strings.Contains(err.Error(), "Server closed") {
                log.Printf("Server error: %v", err)
            }
        },
    )
    if err != nil {
        log.Fatal(err)
    }
}()

// Wait for server to be ready
time.Sleep(100 * time.Millisecond)
```

#### Graceful Shutdown

```go
// Setup signal handling
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

// Wait for signal
<-sigChan

// Graceful shutdown with timeout
log.Println("Shutting down...")
err := exporter.Stop(30 * time.Second)
if err != nil {
    log.Printf("Shutdown error: %v", err)
}
```

### Dynamic Metrics

```go
type DynamicMetric struct {
    desc      *prometheus.Desc
    valueType prometheus.ValueType
    dataSource DataSource  // Your data source
}

func (dm *DynamicMetric) GetValues() []exporter.Value {
    // Fetch fresh values on each scrape
    currentData := dm.dataSource.FetchLatest()
    
    values := make([]exporter.Value, 0, len(currentData))
    for key, value := range currentData {
        values = append(values, exporter.Value{
            Value:       value,
            LabelValues: []string{key},
        })
    }
    return values
}
```

### Labels and Descriptors

```go
// Metric with variable labels
desc := prometheus.NewDesc(
    "api_request_duration_seconds",
    "API request duration in seconds",
    []string{"method", "endpoint", "status"},  // Variable labels
    prometheus.Labels{                          // Constant labels
        "service":  "api",
        "version":  "1.0",
        "environment": "production",
    },
)

// Create values with matching label count
values := []exporter.Value{
    {
        Value:       0.123,
        LabelValues: []string{"GET", "/users", "200"},
    },
    {
        Value:       0.456,
        LabelValues: []string{"POST", "/users", "201"},
    },
}
```

### Best Practices

1. **Metric Naming**: Follow Prometheus naming conventions
   - Use `_total` suffix for counters
   - Use base units (seconds, bytes, etc.)
   - Use snake_case

2. **Label Cardinality**: Keep label cardinality low
   - Avoid high-cardinality labels (IDs, timestamps)
   - Use labels for finite sets of values

3. **Metric Types**: Choose appropriate types
   - Counter: Monotonically increasing values
   - Gauge: Values that can go up and down
   - Histogram: Distribution of values

4. **Performance**: Optimize GetValues()
   - Cache values when possible
   - Avoid expensive computations during scrape
   - Use efficient data structures

5. **Error Handling**: Handle errors gracefully
   - Don't panic in GetValues()
   - Return empty slice on errors
   - Log errors separately

## Integration Examples

### Web Application Metrics

```go
package main

import (
    "log"
    "net/http"
    "sync"
    "time"
    
    "github.com/common-library/go/database/prometheus/exporter"
    "github.com/prometheus/client_golang/prometheus"
)

type AppMetrics struct {
    mu              sync.RWMutex
    requestCount    map[string]int64
    activeRequests  int64
    requestDuration map[string]float64
}

type RequestMetric struct {
    desc    *prometheus.Desc
    metrics *AppMetrics
}

func (rm *RequestMetric) GetDesc() *prometheus.Desc {
    return rm.desc
}

func (rm *RequestMetric) GetValueType() prometheus.ValueType {
    return prometheus.CounterValue
}

func (rm *RequestMetric) GetValues() []exporter.Value {
    rm.metrics.mu.RLock()
    defer rm.metrics.mu.RUnlock()
    
    values := make([]exporter.Value, 0, len(rm.metrics.requestCount))
    for endpoint, count := range rm.metrics.requestCount {
        values = append(values, exporter.Value{
            Value:       float64(count),
            LabelValues: []string{endpoint},
        })
    }
    return values
}

func main() {
    appMetrics := &AppMetrics{
        requestCount:    make(map[string]int64),
        requestDuration: make(map[string]float64),
    }
    
    // Create metrics
    requestDesc := prometheus.NewDesc(
        "app_http_requests_total",
        "Total HTTP requests",
        []string{"endpoint"},
        nil,
    )
    
    requestMetric := &RequestMetric{
        desc:    requestDesc,
        metrics: appMetrics,
    }
    
    // Setup exporter
    collector := exporter.NewCollector([]exporter.Metric{requestMetric})
    exporter.RegisterCollector(collector)
    
    go exporter.Start(":9090", "/metrics", func(err error) {
        log.Println(err)
    })
    
    // Application handler with metrics
    http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
        appMetrics.mu.Lock()
        appMetrics.requestCount[r.URL.Path]++
        appMetrics.mu.Unlock()
        
        w.WriteHeader(http.StatusOK)
    })
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Database Connection Pool Monitoring

```go
type DBPoolMetrics struct {
    desc *prometheus.Desc
    db   *sql.DB
}

func (dpm *DBPoolMetrics) GetDesc() *prometheus.Desc {
    return dpm.desc
}

func (dpm *DBPoolMetrics) GetValueType() prometheus.ValueType {
    return prometheus.GaugeValue
}

func (dpm *DBPoolMetrics) GetValues() []exporter.Value {
    stats := dpm.db.Stats()
    
    return []exporter.Value{
        {Value: float64(stats.OpenConnections), LabelValues: []string{"open"}},
        {Value: float64(stats.InUse), LabelValues: []string{"in_use"}},
        {Value: float64(stats.Idle), LabelValues: []string{"idle"}},
    }
}
```

### Cache Hit Rate

```go
type CacheMetrics struct {
    hitDesc  *prometheus.Desc
    missDesc *prometheus.Desc
    cache    *Cache
}

func (cm *CacheMetrics) GetHitDesc() *prometheus.Desc {
    return cm.hitDesc
}

func (cm *CacheMetrics) GetMissDesc() *prometheus.Desc {
    return cm.missDesc
}

func (cm *CacheMetrics) GetHitValues() []exporter.Value {
    hits := cm.cache.GetHits()
    return []exporter.Value{
        {Value: float64(hits), LabelValues: []string{}},
    }
}

func (cm *CacheMetrics) GetMissValues() []exporter.Value {
    misses := cm.cache.GetMisses()
    return []exporter.Value{
        {Value: float64(misses), LabelValues: []string{}},
    }
}
```

## Testing

### Client Testing

```go
func TestClientQuery(t *testing.T) {
    // Setup Prometheus container
    prometheusEndpoint := setupPrometheus(t)
    
    // Create client
    c, err := client.NewClient("http://" + prometheusEndpoint)
    require.NoError(t, err)
    
    // Execute query
    value, warnings, err := c.Query("up", time.Now(), 10*time.Second)
    assert.NoError(t, err)
    assert.NotNil(t, value)
    
    // Verify warnings
    for _, warning := range warnings {
        t.Logf("Warning: %s", warning)
    }
}
```

### Exporter Testing

```go
func TestExporter(t *testing.T) {
    // Create test metric
    desc := prometheus.NewDesc(
        "test_metric",
        "Test metric",
        []string{"label"},
        nil,
    )
    
    metric := &TestMetric{
        desc:      desc,
        valueType: prometheus.CounterValue,
        values: []exporter.Value{
            {Value: 42.0, LabelValues: []string{"test"}},
        },
    }
    
    // Create and register collector
    collector := exporter.NewCollector([]exporter.Metric{metric})
    err := exporter.RegisterCollector(collector)
    assert.NoError(t, err)
    
    // Start server
    go exporter.Start(":19090", "/metrics", func(err error) {})
    time.Sleep(100 * time.Millisecond)
    
    // Query metrics endpoint
    resp, err := http.Get("http://localhost:19090/metrics")
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Cleanup
    exporter.Stop(5 * time.Second)
}
```

## Troubleshooting

### Client Issues

**Connection refused:**
```go
// Check if Prometheus is running
c, err := client.NewClient("http://localhost:9090")
if err != nil {
    log.Fatal("Failed to create client:", err)
}

// Verify with simple query
_, _, err = c.Query("up", time.Now(), 5*time.Second)
if err != nil {
    log.Fatal("Prometheus not reachable:", err)
}
```

**Query timeout:**
```go
// Increase timeout for complex queries
value, warnings, err := c.Query(
    "complex_query",
    time.Now(),
    30*time.Second,  // Longer timeout
)
```

**Authentication fails:**
```go
// Verify credentials
c, err := client.NewClientWithBasicAuth(
    "http://localhost:9090",
    "correct_username",
    "correct_password",
)
```

### Exporter Issues

**Metrics not appearing:**
```go
// Verify collector is registered
err := exporter.RegisterCollector(collector)
if err != nil {
    log.Fatal("Registration failed:", err)
}

// Check server is running
resp, err := http.Get("http://localhost:9090/metrics")
if err != nil {
    log.Fatal("Server not responding:", err)
}
```

**Label mismatch:**
```go
// Ensure label values match descriptor
desc := prometheus.NewDesc(
    "my_metric",
    "Help text",
    []string{"label1", "label2"},  // 2 labels
    nil,
)

// Must provide exactly 2 label values
values := []exporter.Value{
    {Value: 1.0, LabelValues: []string{"value1", "value2"}},  // ✓ Correct
    {Value: 2.0, LabelValues: []string{"value1"}},             // ✗ Wrong
}
```

**Duplicate metrics:**
```go
// Unregister before re-registering
result := exporter.UnRegisterCollector(collector)
if result {
    err := exporter.RegisterCollector(collector)
}
```

## Performance Considerations

1. **Query Optimization**: Use appropriate time ranges and step intervals
2. **Metric Cardinality**: Limit number of unique label combinations
3. **Scrape Interval**: Balance between freshness and overhead
4. **Data Retention**: Configure appropriate retention policies
5. **Resource Usage**: Monitor exporter memory and CPU usage

## References

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [PromQL Guide](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Writing Exporters](https://prometheus.io/docs/instrumenting/writing_exporters/)
- [Metric Types](https://prometheus.io/docs/concepts/metric_types/)
- [Best Practices](https://prometheus.io/docs/practices/naming/)
