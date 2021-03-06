package collectors

import (
	"encoding/json"
	"errors"
	"strings"
)

type IConfig interface {
	filterQueue(string) bool
}

type QueueJSONParser struct {
	Config		IConfig
	Cmd			string
	Arguments	[]string
}

func NewQueueJSONParser(config IConfig) *QueueJSONParser {
	return &QueueJSONParser{
		Config: config,
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

func (p *QueueJSONParser) GetCmd() string {
	return p.Cmd
}

func (p *QueueJSONParser) GetArguments() []string {
	return p.Arguments
}

func (p *QueueJSONParser) Parse(line string) (*Metrics, error) {
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
	return p.parseQueue(jsonMetrics)
}

func (p *QueueJSONParser) parseQueue(jsonMetrics map[string]interface{}) (*Metrics, error) {
	queueMetrics := NewMetrics()
	queue, state := jsonMetrics["name"].(string), jsonMetrics["state"].(string)
	// If it doesn't go through the filters then we ignore the queue metric
	if !p.Config.filterQueue(queue) {
		return nil, nil
	}
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
