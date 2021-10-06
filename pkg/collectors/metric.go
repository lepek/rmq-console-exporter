package collectors

import (
	"errors"
)

type Metric struct {
	Value		float64
	LabelPairs	map[string]string
}

type Metrics struct {
	MetricPairs	map[string]Metric
}

func NewMetrics() *Metrics {
	return &Metrics{
		MetricPairs: make(map[string]Metric),
	}
}

func (m *Metrics) AddMetric(name string, value float64, labelPairs map[string]string) {
	m.MetricPairs[name] = Metric{
		Value: value,
		LabelPairs: labelPairs,
	}
}

func (m Metrics) GetMetricValue(name string) (float64, error) {
	me, err := m.getMetric(name)
	if err != nil { return 0.0, err }
	return me.Value, nil
}

func (m Metrics) GetLabels(name string) (map[string]string, error) {
	me, err := m.getMetric(name)
	if err != nil { return nil, err }
	return me.LabelPairs, nil
}

func (m Metrics) getMetric(name string) (Metric, error) {
	if me, found := m.MetricPairs[name]; found {
		return me, nil
	}
	return Metric{}, errors.New("metric not found")
}