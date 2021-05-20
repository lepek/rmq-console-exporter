package collectors

type Metric struct {
	name	string
	value	float64
	labels	map[string]string
}

func NewMetric(name string, value float64, labels map[string]string) *Metric {
	return &Metric{
		name: name,
		value: value,
		labels: labels,
	}
}

func (m *Metric) GetField() string {
	return m.name
}

func (m *Metric) GetName() string {
	return m.name
}

func (m *Metric) GetValue() float64 {
	return m.value
}

func (m *Metric) GetLabels() map[string]string {
	return m.labels
}
