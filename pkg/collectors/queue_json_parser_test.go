package collectors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TrueFilterConfig struct{}
func (c *TrueFilterConfig) filterQueue(name string) bool {
	return true
}

type FalseFilterConfig struct{}
func (c *FalseFilterConfig) filterQueue(name string) bool {
	return false
}

func TestQueueJsonParserOk(t *testing.T) {
	var parser ICmdParser
	parser = NewQueueJSONParser(&TrueFilterConfig{})
	line := `{"name":"delegate_encryption_test_3579441e-1f41-4455-90e4-04c3228f1305.tenant_3667d578-644d-4930-b965-4f7bd45ee537.dev","state":"running","messages_ready":1,"message_bytes_ready":288,"messages_unacknowledged":0,"message_bytes_unacknowledged":0,"memory":34788,"consumers":6,"consumer_utilisation":"","head_message_timestamp":1630920836}`
	metrics, err := parser.Parse(line)

	expectedLabels := map[string]string{
		"queue":  "delegate_encryption_test_3579441e-1f41-4455-90e4-04c3228f1305.tenant_3667d578-644d-4930-b965-4f7bd45ee537.dev",
		"state": "running",
	}

	assert.Equal(t, nil, err)

	checkValue(t, metrics, "messages_ready", 1)
	checkLabels(t, metrics, "messages_ready", expectedLabels)

	checkValue(t, metrics, "message_bytes_ready", 288)
	checkLabels(t, metrics, "message_bytes_ready", expectedLabels)

	checkValue(t, metrics, "messages_unacknowledged", 0)
	checkLabels(t, metrics, "messages_unacknowledged", expectedLabels)

	checkValue(t, metrics, "message_bytes_unacknowledged", 0)
	checkLabels(t, metrics, "message_bytes_unacknowledged", expectedLabels)

	checkValue(t, metrics, "memory", 34788)
	checkLabels(t, metrics, "memory", expectedLabels)

	checkValue(t, metrics, "consumers", 6)
	checkLabels(t, metrics, "consumers", expectedLabels)

	checkValue(t, metrics, "head_message_timestamp", 1630920836)
	checkLabels(t, metrics, "head_message_timestamp", expectedLabels)
}

func TestStatusJsonParser(t *testing.T) {
	parser := NewQueueJSONParser(&TrueFilterConfig{})
	line := `{"command_executed":"rabbitmqctl list_queues --formatter json name","command_runtime":0.5655179}`
	metrics, err := parser.Parse(line)
	assert.Equal(t, nil, err)
	checkValue(t, metrics, "command_runtime", 0.5655179)
	checkLabels(t, metrics, "command_runtime", map[string]string{"command_executed":"rabbitmqctl list_queues --formatter json name"})
}

func TestQueueJsonParserTrailingOk(t *testing.T) {
	var parser ICmdParser
	parser = NewQueueJSONParser(&TrueFilterConfig{})
	line := `,{"name":"delegate_encryption_test_3579441e-1f41-4455-90e4-04c3228f1305.tenant_3667d578-644d-4930-b965-4f7bd45ee537.dev","state":"running","messages_ready":1,"message_bytes_ready":288,"messages_unacknowledged":0,"message_bytes_unacknowledged":0,"memory":34788,"consumers":6,"consumer_utilisation":"","head_message_timestamp":1630920836}`
	metrics, err := parser.Parse(line)

	expectedLabels := map[string]string{
		"queue":  "delegate_encryption_test_3579441e-1f41-4455-90e4-04c3228f1305.tenant_3667d578-644d-4930-b965-4f7bd45ee537.dev",
		"state": "running",
	}

	assert.Equal(t, nil, err)

	checkValue(t, metrics, "messages_ready", 1)
	checkLabels(t, metrics, "messages_ready", expectedLabels)

	checkValue(t, metrics, "message_bytes_ready", 288)
	checkLabels(t, metrics, "message_bytes_ready", expectedLabels)

	checkValue(t, metrics, "messages_unacknowledged", 0)
	checkLabels(t, metrics, "messages_unacknowledged", expectedLabels)

	checkValue(t, metrics, "message_bytes_unacknowledged", 0)
	checkLabels(t, metrics, "message_bytes_unacknowledged", expectedLabels)

	checkValue(t, metrics, "memory", 34788)
	checkLabels(t, metrics, "memory", expectedLabels)

	checkValue(t, metrics, "consumers", 6)
	checkLabels(t, metrics, "consumers", expectedLabels)

	checkValue(t, metrics, "head_message_timestamp", 1630920836)
	checkLabels(t, metrics, "head_message_timestamp", expectedLabels)
}

func TestQueueJsonParserMatchNotFound(t *testing.T) {
	var parser ICmdParser
	parser = NewQueueJSONParser(&TrueFilterConfig{})
	line := `[`
	metrics, err := parser.Parse(line)
	assert.IsType(t, &NonFatalError{}, err)
	assert.Nil(t, metrics)
}

func TestQueueJsonParserFiltered(t *testing.T) {
	var parser ICmdParser
	parser = NewQueueJSONParser(&FalseFilterConfig{})
	line := `{"name":"delegate_encryption_test_3579441e-1f41-4455-90e4-04c3228f1305.tenant_3667d578-644d-4930-b965-4f7bd45ee537.dev","state":"running","messages_ready":1,"message_bytes_ready":288,"messages_unacknowledged":0,"message_bytes_unacknowledged":0,"memory":34788,"consumers":6,"consumer_utilisation":"","head_message_timestamp":1630920836}`
	metrics, err := parser.Parse(line)

	assert.Nil(t, err)
	assert.Nil(t, metrics)
}
