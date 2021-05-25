package exporters

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rmq-console-exporter/pkg/collectors"
	"testing"
)

var (
	expected = fmt.Sprintf(
		"Desc{fqName: %q, help: %q, constLabels: {%s}, variableLabels: %v}",
		"prefix_memory",
		"Bytes of memory allocated by the runtime for the queue, including stack, heap and internal structures.",
		"",
		[]string{"queue", "state"},
	)
)

type MockedCollector struct{
	mock.Mock
}

func (c *MockedCollector) Collect() ([]collectors.IMetrics, error) {
	args := c.Called()
	metric := collectors.NewMetrics()
	metric.AddMetric("Memory", 1.5, map[string]string{"queue": "q33"})
	return []collectors.IMetrics{metric}, args.Error(0)
}

func TestExporterOk(t *testing.T) {
	testCollector := new(MockedCollector)
	exporter := buildTestExporter([]collectors.ICollector{testCollector})
	testCollector.On("Collect").Return(nil)
	ch := make(chan prometheus.Metric, 1)
	exporter.Collect(ch)
	metric := <- ch
	assert.Equal(t, expected, metric.Desc().String())
}

func TestExporterMultipleCollectorsOneFailing(t *testing.T) {
	testCollectorOk := new(MockedCollector)
	testCollectorFail := new(MockedCollector)
	exporter := buildTestExporter([]collectors.ICollector{testCollectorFail, testCollectorOk})
	testCollectorOk.On("Collect").Return(nil)
	testCollectorFail.On("Collect").Return(errors.New("some error"))
	ch := make(chan prometheus.Metric, 2)
	exporter.Collect(ch)
	assert.Equal(t, 1, len(ch))

	metric := <- ch
	assert.Equal(t, expected, metric.Desc().String())
}

func TestExporterMultipleCollectorsOk(t *testing.T) {
	testCollectorOk1 := new(MockedCollector)
	testCollectorOk2 := new(MockedCollector)
	exporter := buildTestExporter([]collectors.ICollector{testCollectorOk1, testCollectorOk2})
	testCollectorOk1.On("Collect").Return(nil)
	testCollectorOk2.On("Collect").Return(nil)
	ch := make(chan prometheus.Metric, 2)
	exporter.Collect(ch)
	assert.Equal(t, 2, len(ch))
	metric := <- ch
	assert.Equal(t, expected, metric.Desc().String())
}

func TestExporterMultipleCollectorsFailing(t *testing.T) {
	testCollectorFail1 := new(MockedCollector)
	testCollectorFail2 := new(MockedCollector)
	exporter := buildTestExporter([]collectors.ICollector{testCollectorFail1, testCollectorFail2})
	testCollectorFail1.On("Collect").Return(errors.New("some error"))
	testCollectorFail2.On("Collect").Return(errors.New("some error"))
	ch := make(chan prometheus.Metric, 2)
	exporter.Collect(ch)
	assert.Equal(t, 0, len(ch))
}

func buildTestExporter(c []collectors.ICollector) *PrometheusExporter {
	return NewPrometheusExporter("prefix_", 9999, c)
}


