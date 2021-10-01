package collectors

import (
	"encoding/json"
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
			"--formatter json",
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
