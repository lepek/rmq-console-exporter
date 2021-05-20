package collectors

import (
	"reflect"
)

type RMQQueueMetrics struct {
	Name 						string	`regroup:"Name"`
	State 						string	`regroup:"State"`
	MessagesReady				float64	`regroup:"MessagesReady"`
	MessageBytesReady			float64	`regroup:"MessageBytesReady"`
	MessagesUnacknowledged		float64	`regroup:"MessagesUnacknowledged"`
	MessageBytesUnacknowledged	float64	`regroup:"MessageBytesUnacknowledged"`
	Memory						float64	`regroup:"Memory"`
	Consumers					float64	`regroup:"Consumers"`
	ConsumerUtilisation			float64	`regroup:"ConsumerUtilisation"`
	HeadMessageTimestamp		float64	`regroup:"HeadMessageTimestamp"`
}

func (m *RMQQueueMetrics) GetMetricByName(name string) float64 {
	r := reflect.ValueOf(m)
	f := reflect.Indirect(r).FieldByName(name)
	return f.Float()
}

func (m *RMQQueueMetrics) GetLabelsForMetric(name string) map[string]string {
	// We can ignore the metric name for now since all the metrics have the same labels right now
	labelPairs := make(map[string]string, 2)
	labelPairs["queue"] = m.Name
	labelPairs["state"] = m.State
	return labelPairs
}