package exporter_test

import (
	"testing"

	"github.com/common-library/go/database/prometheus/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

type TestMetric struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	values    []exporter.Value
}

func (tm *TestMetric) GetDesc() *prometheus.Desc {
	return tm.desc
}

func (tm *TestMetric) GetValueType() prometheus.ValueType {
	return tm.valueType
}

func (tm *TestMetric) GetValues() []exporter.Value {
	return tm.values
}

func TestMetricInterface(t *testing.T) {
	desc := prometheus.NewDesc(
		"test_metric",
		"A test metric",
		[]string{"label1", "label2"},
		prometheus.Labels{"const_label": "const_value"},
	)

	values := []exporter.Value{
		{Value: 1.0, LabelValues: []string{"value1", "value2"}},
		{Value: 2.0, LabelValues: []string{"value3", "value4"}},
	}

	testMetric := &TestMetric{
		desc:      desc,
		valueType: prometheus.CounterValue,
		values:    values,
	}

	assert.Equal(t, desc, testMetric.GetDesc())
	assert.Equal(t, prometheus.CounterValue, testMetric.GetValueType())
	assert.Equal(t, values, testMetric.GetValues())
}

func TestNewCollector(t *testing.T) {
	desc1 := prometheus.NewDesc(
		"test_counter",
		"A test counter",
		[]string{"label1"},
		nil,
	)

	desc2 := prometheus.NewDesc(
		"test_gauge",
		"A test gauge",
		[]string{"label2"},
		nil,
	)

	metrics := []exporter.Metric{
		&TestMetric{
			desc:      desc1,
			valueType: prometheus.CounterValue,
			values: []exporter.Value{
				{Value: 10.0, LabelValues: []string{"counter_value"}},
			},
		},
		&TestMetric{
			desc:      desc2,
			valueType: prometheus.GaugeValue,
			values: []exporter.Value{
				{Value: 5.5, LabelValues: []string{"gauge_value"}},
			},
		},
	}

	collector := exporter.NewCollector(metrics)
	assert.NotNil(t, collector)

	var _ prometheus.Collector = collector
}

func TestCollectorDescribe(t *testing.T) {
	desc := prometheus.NewDesc(
		"test_metric_describe",
		"A test metric for describe",
		[]string{"label"},
		nil,
	)

	metrics := []exporter.Metric{
		&TestMetric{
			desc:      desc,
			valueType: prometheus.CounterValue,
			values: []exporter.Value{
				{Value: 1.0, LabelValues: []string{"test"}},
			},
		},
	}

	collector := exporter.NewCollector(metrics)

	descChan := make(chan *prometheus.Desc, 1)
	go func() {
		collector.Describe(descChan)
		close(descChan)
	}()

	receivedDesc := <-descChan
	assert.Equal(t, desc, receivedDesc)
}

func TestCollectorCollect(t *testing.T) {
	desc := prometheus.NewDesc(
		"test_metric_collect",
		"A test metric for collect",
		[]string{"label"},
		nil,
	)

	expectedValue := 42.0
	expectedLabel := "test_label"

	metrics := []exporter.Metric{
		&TestMetric{
			desc:      desc,
			valueType: prometheus.CounterValue,
			values: []exporter.Value{
				{Value: expectedValue, LabelValues: []string{expectedLabel}},
			},
		},
	}

	collector := exporter.NewCollector(metrics)

	metricChan := make(chan prometheus.Metric, 1)
	go func() {
		collector.Collect(metricChan)
		close(metricChan)
	}()

	receivedMetric := <-metricChan
	assert.NotNil(t, receivedMetric)

	metricDesc := receivedMetric.Desc()
	assert.Equal(t, desc, metricDesc)
}

func TestValueStruct(t *testing.T) {
	value := exporter.Value{
		Value:       123.45,
		LabelValues: []string{"label1", "label2", "label3"},
	}

	assert.Equal(t, 123.45, value.Value)
	assert.Equal(t, []string{"label1", "label2", "label3"}, value.LabelValues)
	assert.Len(t, value.LabelValues, 3)
}

func TestMultipleMetricsCollector(t *testing.T) {
	desc1 := prometheus.NewDesc("metric1", "First metric", []string{"l1"}, nil)
	desc2 := prometheus.NewDesc("metric2", "Second metric", []string{"l2"}, nil)
	desc3 := prometheus.NewDesc("metric3", "Third metric", []string{"l3"}, nil)

	metrics := []exporter.Metric{
		&TestMetric{
			desc:      desc1,
			valueType: prometheus.CounterValue,
			values: []exporter.Value{
				{Value: 1.0, LabelValues: []string{"v1"}},
				{Value: 2.0, LabelValues: []string{"v2"}},
			},
		},
		&TestMetric{
			desc:      desc2,
			valueType: prometheus.GaugeValue,
			values: []exporter.Value{
				{Value: 3.0, LabelValues: []string{"v3"}},
			},
		},
		&TestMetric{
			desc:      desc3,
			valueType: prometheus.UntypedValue,
			values: []exporter.Value{
				{Value: 4.0, LabelValues: []string{"v4"}},
				{Value: 5.0, LabelValues: []string{"v5"}},
				{Value: 6.0, LabelValues: []string{"v6"}},
			},
		},
	}

	collector := exporter.NewCollector(metrics)

	descChan := make(chan *prometheus.Desc, 10)
	go func() {
		collector.Describe(descChan)
		close(descChan)
	}()

	var descriptions []*prometheus.Desc
	for desc := range descChan {
		descriptions = append(descriptions, desc)
	}

	assert.Len(t, descriptions, 3)
	assert.Contains(t, descriptions, desc1)
	assert.Contains(t, descriptions, desc2)
	assert.Contains(t, descriptions, desc3)

	metricChan := make(chan prometheus.Metric, 10)
	go func() {
		collector.Collect(metricChan)
		close(metricChan)
	}()

	var collectedMetrics []prometheus.Metric
	for metric := range metricChan {
		collectedMetrics = append(collectedMetrics, metric)
	}

	assert.Len(t, collectedMetrics, 6)
}
