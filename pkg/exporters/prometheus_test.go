package exporters

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type TestMetrics struct {}

func (m TestMetrics) GetMetricValue(name string) (float64, error) {
	if name != "memory" {
		return 0.0, errors.New("metric not found")
	}
	return 1.5, nil
}

func (m TestMetrics) GetLabels(name string) (map[string]string, error) {
	if name != "memory" {
		return nil, errors.New("metric not found")
	}
	return map[string]string{"queue": "q33", "state":"running"}, nil
}

type MockedCollector struct{
	mock.Mock
}

func (c *MockedCollector) Collect() ([]IMetrics, error) {
	args := c.Called()
	return []IMetrics{&TestMetrics{}}, args.Error(0)
}

func TestExporterOk(t *testing.T) {
	testCollector := new(MockedCollector)
	exporter := buildTestExporter([]ICollector{testCollector})
	testCollector.On("Collect").Return(nil)
	ch := make(chan prometheus.Metric, 10)
	exporter.Collect(ch)
	assert.Equal(t, 1, len(ch))
	metric := <- ch
	assert.Equal(t, expected, metric.Desc().String())
}

func TestExporterMultipleCollectorsOneFailing(t *testing.T) {
	testCollectorOk := new(MockedCollector)
	testCollectorFail := new(MockedCollector)
	exporter := buildTestExporter([]ICollector{testCollectorFail, testCollectorOk})
	testCollectorOk.On("Collect").Return(nil)
	testCollectorFail.On("Collect").Return(errors.New("some error"))
	ch := make(chan prometheus.Metric, 10)
	exporter.Collect(ch)
	assert.Equal(t, 1, len(ch))

	metric := <- ch
	assert.Equal(t, expected, metric.Desc().String())
}

func TestExporterMultipleCollectorsOk(t *testing.T) {
	testCollectorOk1 := new(MockedCollector)
	testCollectorOk2 := new(MockedCollector)
	exporter := buildTestExporter([]ICollector{testCollectorOk1, testCollectorOk2})
	testCollectorOk1.On("Collect").Return(nil)
	testCollectorOk2.On("Collect").Return(nil)
	ch := make(chan prometheus.Metric, 10)
	exporter.Collect(ch)
	assert.Equal(t, 2, len(ch))
	metric := <- ch
	assert.Equal(t, expected, metric.Desc().String())
}

func TestExporterMultipleCollectorsFailing(t *testing.T) {
	testCollectorFail1 := new(MockedCollector)
	testCollectorFail2 := new(MockedCollector)
	exporter := buildTestExporter([]ICollector{testCollectorFail1, testCollectorFail2})
	testCollectorFail1.On("Collect").Return(errors.New("some error"))
	testCollectorFail2.On("Collect").Return(errors.New("some error"))
	ch := make(chan prometheus.Metric, 10)
	exporter.Collect(ch)
	assert.Equal(t, 0, len(ch))
}

func buildTestExporter(c []ICollector) *PrometheusExporter {
	return NewPrometheusExporter("prefix_", 9999, c)
}
