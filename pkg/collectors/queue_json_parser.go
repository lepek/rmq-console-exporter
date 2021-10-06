package collectors

import (
	"encoding/json"
	"errors"
	"strings"
)

type QueueJsonParser struct {
	Cmd			string
	Arguments	[]string
}

func NewQueueJsonParser() *QueueJsonParser {
	return &QueueJsonParser{
		Cmd: "rabbitmqctl",
		Arguments: []string{
			"list_queues",
			"--formatter",
			"json",
			"name",
			"state",
			"messages_ready",
			"message_bytes_ready",
			"messages_unacknowledged",
			"message_bytes_unacknowledged",
			"memory",
			"consumers",
			"consumer_utilisation",
			"head_message_timestamp",
		},
	}
}

func (p *QueueJsonParser) GetCmd() string {
	return p.Cmd
}

func (p *QueueJsonParser) GetArguments() []string {
	return p.Arguments
}

func (p *QueueJsonParser) Parse(line string) (*Metrics, error) {
	var jsonMetrics map[string]interface{}
	err := json.Unmarshal([]byte(strings.Trim(line,",")), &jsonMetrics)
	if err != nil {
		return nil, NewNonFatalError(err)
	}

	// Not super elegant to check to identify a valid queue metric
	_, okQueue := jsonMetrics["name"]
	_, okState := jsonMetrics["state"]

	if !okQueue || !okState {
		return parseStatus(jsonMetrics)
	}
	return parseQueue(jsonMetrics)
}

func parseQueue(jsonMetrics map[string]interface{}) (*Metrics, error) {
	queueMetrics := NewMetrics()
	queue, state := jsonMetrics["name"].(string), jsonMetrics["state"].(string)
	for name, value := range jsonMetrics {
		if name == "name" || name == "state" { continue }
		fValue, ok := value.(float64)
		if !ok { continue }
		queueMetrics.AddMetric(name, fValue, map[string]string{"queue": queue, "state": state})
	}

	return queueMetrics, nil
}

func parseStatus(jsonMetrics map[string]interface{}) (*Metrics, error) {
	commandRuntime, okRuntime := jsonMetrics["command_runtime"]
	commandExecuted, okExecuted := jsonMetrics["command_executed"]
	if okRuntime && okExecuted {
		statusMetrics := NewMetrics()
		labels := map[string]string{"command_executed": commandExecuted.(string)}
		statusMetrics.AddMetric("command_runtime", commandRuntime.(float64), labels)
		return statusMetrics, nil
	}
	return nil, NewNonFatalError(errors.New("unknown JSON line"))
}
