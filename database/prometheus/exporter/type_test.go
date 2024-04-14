package exporter_test

import (
	"io"
	"strings"

	"github.com/common-library/go/database/prometheus/exporter"
	"github.com/prometheus/client_golang/prometheus"
)

type metric01 struct {
}

func (this *metric01) GetDesc() *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName("sample01", "", "metric01"),
		"metric01 of sample01",
		[]string{"label_01", "label_02"},
		prometheus.Labels{
			"const_label_01": "const-value-01",
			"const_label_02": "const-value-02"})
}

func (this *metric01) GetValueType() prometheus.ValueType {
	return prometheus.GaugeValue
}

func (this *metric01) GetValues() []exporter.Value {
	return []exporter.Value{
		exporter.Value{Value: 1, LabelValues: []string{"value-01", "value-02"}},
		exporter.Value{Value: 1, LabelValues: []string{"value-03", "value-04"}},
		exporter.Value{Value: 2, LabelValues: []string{"value-05", "value-06"}},
	}
}

func (this *metric01) getExpected() io.Reader {
	expected := `# HELP sample01_metric01 metric01 of sample01
       # TYPE sample01_metric01 gauge
       sample01_metric01{const_label_01="const-value-01",const_label_02="const-value-02",label_01="value-01",label_02="value-02"} 1
       sample01_metric01{const_label_01="const-value-01",const_label_02="const-value-02",label_01="value-03",label_02="value-04"} 1
       sample01_metric01{const_label_01="const-value-01",const_label_02="const-value-02",label_01="value-05",label_02="value-06"} 2
       `

	return strings.NewReader(expected)
}
