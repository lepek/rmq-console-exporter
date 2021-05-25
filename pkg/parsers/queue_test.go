package parsers

import (
	"github.com/stretchr/testify/assert"
	"rmq-console-exporter/pkg/collectors"
	"testing"
)

func TestQueueParserOkWithNoMessageTimestamp(t *testing.T) {
	parser := NewQueueParser()
	line := "q33	running	1	2	3	4	34664	5	6	"
	metrics, err := parser.Parse(line)
	assert.Equal(t, nil, err)

	checkValue(t, metrics, "MessagesReady", 1)
	checkLabels(t, metrics, "MessagesReady")

	checkValue(t, metrics, "MessageBytesReady", 2)
	checkLabels(t, metrics, "MessageBytesReady")

	checkValue(t, metrics, "MessagesUnacknowledged", 3)
	checkLabels(t, metrics, "MessagesUnacknowledged")

	checkValue(t, metrics, "MessageBytesUnacknowledged", 4)
	checkLabels(t, metrics, "MessageBytesUnacknowledged")

	checkValue(t, metrics, "Memory", 34664)
	checkLabels(t, metrics, "Memory")

	checkValue(t, metrics, "Consumers", 5)
	checkLabels(t, metrics, "Consumers")

	checkValue(t, metrics, "ConsumerUtilisation", 6)
	checkLabels(t, metrics, "ConsumerUtilisation")

	_, err = metrics.GetMetricValue("HeadMessageTimestamp")
	assert.NotEqual(t, nil, err)
}

func TestQueueParserMatchNotFound(t *testing.T) {
	parser := NewQueueParser()
	line := "some random string"
	metrics, err := parser.Parse(line)
	assert.Equal(t, nil, metrics)
	assert.Equal(t, nil, err)
}

func checkLabels(t *testing.T, metrics collectors.IMetrics, metricName string) {
	l, err := metrics.GetLabels(metricName)
	assert.Equal(t, nil, err)
	assert.Equal(t, "q33", l["queue"])
	assert.Equal(t, "running", l["state"])
}

func checkValue(t *testing.T, metrics collectors.IMetrics, metricName string, expected float64) {
	v, err := metrics.GetMetricValue(metricName)
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, v)
}