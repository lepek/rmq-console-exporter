package collectors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueueParserOkWithNoMessageTimestamp(t *testing.T) {
	parser := NewQueueParser()
	line := "q33	running	1	2	3	4	34664	5	0.1	"
	metrics, err := parser.Parse(line)

	expectedLabels := map[string]string{
		"queue": "q33",
		"state": "running",
	}

	assert.Equal(t, nil, err)

	checkValue(t, metrics, "messages_ready", 1)
	checkLabels(t, metrics, "messages_ready", expectedLabels)

	checkValue(t, metrics, "message_bytes_ready", 2)
	checkLabels(t, metrics, "message_bytes_ready", expectedLabels)

	checkValue(t, metrics, "messages_unacknowledged", 3)
	checkLabels(t, metrics, "messages_unacknowledged", expectedLabels)

	checkValue(t, metrics, "message_bytes_unacknowledged", 4)
	checkLabels(t, metrics, "message_bytes_unacknowledged", expectedLabels)

	checkValue(t, metrics, "memory", 34664)
	checkLabels(t, metrics, "memory", expectedLabels)

	checkValue(t, metrics, "consumers", 5)
	checkLabels(t, metrics, "consumers", expectedLabels)

	checkValue(t, metrics, "consumer_utilisation", 0.1)
	checkLabels(t, metrics, "consumer_utilisation", expectedLabels)

	_, err = metrics.GetMetricValue("head_message_timestamp")
	assert.NotEqual(t, nil, err)
}

func TestQueueParserMatchNotFound(t *testing.T) {
	parser := NewQueueParser()
	line := "some random string"
	metrics, err := parser.Parse(line)
	assert.IsType(t, &NonFatalError{}, err)
	assert.Nil(t, metrics)
}

func checkLabels(t *testing.T, metrics *Metrics, metricName string, expectedLabels map[string]string) {
	l, err := metrics.GetLabels(metricName)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedLabels["queue"], l["queue"])
	assert.Equal(t, expectedLabels["state"], l["state"])
}

func checkValue(t *testing.T, metrics *Metrics, metricName string, expected float64) {
	v, err := metrics.GetMetricValue(metricName)
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, v)
}